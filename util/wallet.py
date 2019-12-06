from aioredis_lock import RedisLock
from aiohttp import log

from db.models.account import Account
from db.models.wallet import Wallet
from db.redis import RedisDB
from network.rpc_client import RPCClient, AccountNotFound
from network.work_client import WorkClient

import config
import nanopy
import rapidjson

class WalletUtil(object):
    """Wallet utilities, like signing, creating blocks, etc."""

    def __init__(self, acct: Account, wallet: Wallet):
        self.account = acct
        self.wallet = wallet

    def get_representative(self):
        if self.wallet.representative is None:
            return config.Config.instance().get_random_rep()
        return self.wallet.representative

    async def receive(self, hash: str) -> str:
        """Receive a block and return hash of published block"""
        async with RedisLock(
            await RedisDB.instance().get_redis(),
            key=f"pippin:{self.account.address}",
            timeout=30,
            wait_timeout=30            
        ) as lock:
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
            try:
                work = await WorkClient.instance().work_generate(workbase)
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
            state_block['link'] = hash
            state_block['work'] = work

            # Sign block
            pk = self.account.private_key(self.wallet.seed)
            state_block['signature'] = nanopy.sign(pk, block=state_block)

            # Publish block
            process_hash = await RPCClient.instance().process(state_block)
            if process_hash is None:
                raise ProcessFailed(rapidjson.dumps(state_block))
            return process_hash

class WorkFailed(Exception):
    pass

class ProcessFailed(Exception):
    pass