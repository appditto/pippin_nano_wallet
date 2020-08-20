from typing import TYPE_CHECKING, List, Union

import aioredis
import nanopy
import rapidjson
from aiohttp import log
from aioredlock import LockError

import pippin.config as config
import pippin.util.nano_util as nano_util
from pippin.db.redis import RedisDB
from pippin.db.models.account import Account
from pippin.db.models.adhoc_account import AdHocAccount
from pippin.db.models.block import Block
from pippin.network.rpc_client import AccountNotFound, RPCClient
from pippin.network.work_client import WorkClient

if TYPE_CHECKING:
    from pippin.db.models.wallet import Wallet

RECEIVE_DIFFICULTY = 'fffffe0000000000' if config.Config.instance().banano else 'ffffffc000000000'#'fffffe0000000000'
SEND_DIFFICULTY = 'fffffe0000000000' if config.Config.instance().banano else 'fffffff800000000'

class WalletUtil(object):
    """Wallet utilities, like signing, creating blocks, etc."""

    def __init__(self, acct: Union[Account, AdHocAccount], wallet: 'Wallet'):
        self.account = acct
        self.wallet = wallet

    def get_representative(self):
        if self.wallet.representative is None:
            return config.Config.instance().get_random_rep()
        return self.wallet.representative

    def adhoc(self) -> bool:
        return not isinstance(self.account, Account)

    def private_key(self) -> str:
        if self.adhoc():
            return self.account.private_key_get()
        return self.account.private_key(self.wallet.seed)


    async def publish(self, state_block: dict, subtype: str = None) -> dict:
        """Publish a state block"""
        try:
            resp = await RPCClient.instance().process(state_block, subtype=subtype)
            # The RPC send/receive uses `block` as a key instead of `hash`
            if 'hash' in resp:
                return {'block': resp['hash']}
            return resp
        except Exception:
            from pippin.db.models.wallet import ProcessFailed
            raise ProcessFailed()

    async def _receive_block_create(self, hash: str, work: str = None) -> dict:
        """Build a state block (receive)"""
        # Get block info
        block_info = await RPCClient.instance().block_info(hash)
        if block_info is None or block_info['contents']['link_as_account'] != self.account.address:
            return None
        # Get account info
        is_open = True
        try:
            account_info = await RPCClient.instance().account_info(self.account.address)
            if account_info is None:
                return None
        except AccountNotFound:
            is_open = False

        # Different workbase for open/receive
        if is_open:
            workbase = account_info['frontier']
        else:
            workbase = nanopy.account_key(self.account.address)

        # Build other fields
        previous = '0000000000000000000000000000000000000000000000000000000000000000' if not is_open else account_info['frontier']
        representative = self.get_representative() if not is_open else account_info['representative']
        balance = block_info['amount'] if not is_open else str(int(account_info['balance']) + int(block_info['amount']))

        # Generate work
        if work is None:
            try:
                work = await WorkClient.instance().work_generate(workbase, RECEIVE_DIFFICULTY)
                if work is None:
                    log.server_logger.error("WORK FAILED")
                    raise WorkFailed(workbase)
            except Exception:
                log.server_logger.exception("work failed")
                raise WorkFailed(workbase)

        # Build final state block
        state_block = nanopy.state_block()
        state_block['account'] = self.account.address
        state_block['previous'] = previous
        state_block['representative'] = representative
        state_block['balance'] = balance
        state_block['link'] = hash
        state_block['work'] = work

        # Sign block
        pk = self.private_key()
        state_block['signature'] = await nano_util.NanoUtil.instance().sign_block(pk, block=state_block)

        return state_block

    async def receive(self, hash: str, work: str = None) -> dict:
        """Receive a block and return hash of published block"""
        async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{self.account.address}") as lock:
            block = await self._receive_block_create(hash, work)
            return await self.publish(block, subtype='receive')

    async def _receive_all(self) -> int:
        """Receive and publish multiple blocks"""
        received_count = 0
        p = await RPCClient.instance().pending(self.account.address, threshold=config.Config.instance().receive_minimum)
        if p is None:
            return received_count
        for block in p:
            state_block = await self._receive_block_create(block)
            await self.publish(state_block, subtype='receive')
            received_count += 1
        return received_count

    async def receive_all(self) -> int:
        """Receive all pending blocks for this account and return # received"""
        received_count = 0
        async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{self.account.address}") as lock:
            received_count = await self._receive_all()
        return received_count

    async def _send_block_create(self, amount: int, destination: str, id: str = None, work: str = None) -> dict:
        """Create a state block (send)"""
        # Get account info
        is_open = True
        account_info = await RPCClient.instance().account_info(self.account.address)
        if account_info is None:
            return None

        # Check balance
        if amount > int(account_info['balance']):
            # Auto-receive blocks if they have it pending
            if config.Config.instance().auto_receive_on_send and int(account_info['balance']) + int(account_info['pending']) >= amount:
                await self._receive_all()
                account_info = await RPCClient.instance().account_info(self.account.address)
                if account_info is None:
                    return None
                if amount > int(account_info['balance']):
                    raise InsufficientBalance(account_info['balance'])
            else:
                raise InsufficientBalance(account_info['balance'])

        workbase = account_info['frontier']

        # Build other fields
        previous = account_info['frontier']
        representative = account_info['representative']
        balance = str(int(account_info['balance']) - amount)

        # Generate work
        if work is None:
            try:
                work = await WorkClient.instance().work_generate(workbase, SEND_DIFFICULTY)
                if work is None:
                    raise WorkFailed(workbase)
            except Exception:
                raise WorkFailed(workbase)

        # Build final state block
        state_block = nanopy.state_block()
        state_block['account'] = self.account.address
        state_block['previous'] = previous
        state_block['representative'] = representative
        state_block['balance'] = balance
        state_block['link'] = nanopy.account_key(destination.replace("xrb_", "nano_"))
        state_block['work'] = work

        # Sign block
        pk = self.private_key()
        state_block['signature'] = await nano_util.NanoUtil.instance().sign_block(pk, block=state_block)

        return state_block

    async def send(self, amount: int, destination: str, id: str = None, work: str = None) -> dict:
        """Create a send block and return hash of published block
            amount is in RAW"""
        
        async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{self.account.address}") as lock:
            # See if block exists, if ID specified
            # If so just rebroadcast it and return the hash
            if id is not None:
                if not self.adhoc():
                    block = await Block.filter(send_id=id, account=self.account).first()
                else:
                    block = await Block.filter(send_id=id, adhoc_account=self.account).first()
                if block is not None:
                    await RPCClient.instance().process(block.block)
                    return {"block": block.block_hash.upper()}
            # Create block
            state_block = await self._send_block_create(amount, destination, id=id, work=work)
            # Publish block
            resp = await self.publish(state_block, subtype='send')
            # Cache if ID specified
            if resp is not None and 'block' in resp:
                # Cache block in database if it has id specified
                if id is not None:
                    block = Block(
                        account=self.account if not self.adhoc() else None,
                        adhoc_account=self.account if self.adhoc() else None,
                        block_hash=nanopy.block_hash(state_block),
                        block=state_block,
                        send_id=id,
                        subtype='send'
                    )
                    await block.save()
            return resp

    async def _change_block_create(self, representative: str, work: str = None, only_if_different: bool = False) -> dict:
        """Create a state block (change)"""
        # Get account info
        account_info = await RPCClient.instance().account_info(self.account.address)
        if account_info is None:
            return None
        elif only_if_different and account_info['representative'] == representative:
            return None

        workbase = account_info['frontier']

        # Build other fields
        previous = account_info['frontier']
        representative = representative
        balance = account_info['balance']

        # Generate work
        if work is None:
            try:
                work = await WorkClient.instance().work_generate(workbase, SEND_DIFFICULTY)
                if work is None:
                    raise WorkFailed(workbase)
            except Exception:
                raise WorkFailed(workbase)

        # Build final state block
        state_block = nanopy.state_block()
        state_block['account'] = self.account.address
        state_block['previous'] = previous
        state_block['representative'] = representative
        state_block['balance'] = balance
        state_block['link'] = '0000000000000000000000000000000000000000000000000000000000000000'
        state_block['work'] = work

        # Sign block
        pk = self.private_key()
        state_block['signature'] = await nano_util.NanoUtil.instance().sign_block(pk, block=state_block)

        return state_block

    async def representative_set(self, representative: str, work: str = None, only_if_different: bool = False) -> dict:
        """Create a change block and return hash of published block"""
        async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{self.account.address}") as lock:
            state_block = await self._change_block_create(representative, work=work, only_if_different=only_if_different)

            if state_block is None and only_if_different:
                return None

            # Publish
            return await self.publish(state_block, subtype='change')

class WorkFailed(Exception):
    pass

class ProcessFailed(Exception):
    pass

class InsufficientBalance(Exception):
    pass
