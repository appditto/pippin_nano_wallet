import datetime
import logging
from json import JSONDecodeError

import rapidjson as json
from aiohttp import log, web
from tortoise.transactions import in_transaction

import asyncio
import pippin.config as config
from pippin.db.redis import RedisDB
from pippin.db.models.account import Account
from pippin.db.models.adhoc_account import AdHocAccount
from pippin.db.models.wallet import (AccountAlreadyExists, Wallet, WalletLocked,
                              WalletNotFound)
from pippin.network.rpc_client import AccountNotFound, BlockNotFound, RPCClient
from pippin.network.nano_websocket import WebsocketClient
from pippin.util.crypt import DecryptionError
from pippin.util.random import RandomUtil
from pippin.util.validators import Validators
from pippin.util.wallet import (InsufficientBalance, ProcessFailed, WalletUtil,
                         WorkFailed)

class PippinServer(object):
    """API for wallet requests"""
    def __init__(self, host: str, port: int):
        self.app = web.Application(middlewares=[web.normalize_path_middleware()])
        self.app.add_routes([
            web.post('/', self.gateway)
        ])
        self.host = host
        self.port = port
        self.websocket = None
        if config.Config.instance().node_ws_url is not None:
            self.websocket = WebsocketClient(config.Config.instance().node_ws_url, self.block_arrival_handler)

    async def stop(self):
        await self.app.shutdown()
        if self.websocket:
            await self.websocket.close()

    def json_response(self, data: dict):
        """Wrapper for json responses using custom json parser"""
        return web.json_response(
            data=data,
            dumps=json.dumps
        )

    def generic_error(self):
        """The node returns this generic error when the request is bad"""
        return self.json_response(
            data={
                'error':"Unable to parse json"
            }
        )

    async def gateway(self, request: web.Request):
        """Gateway route to mimic nano's API of specifying action in a string"""
        try:
            request_json = await request.json(loads=json.loads)    
        except json.JSONDecodeError:
            return self.generic_error()
        if 'action' in request_json:
            # Sanitize action
            request_json['action'] = request_json['action'].lower().strip()

            # Handle wallet RPCs
            if request_json['action'] == 'wallet_create':
                return await self.wallet_create(request, request_json)
            elif request_json['action'] == 'account_create':
                return await self.account_create(request, request_json)
            elif request_json['action'] == 'accounts_create':
                return await self.accounts_create(request, request_json)
            elif request_json['action'] == 'account_list':
                return await self.account_list(request, request_json)
            elif request_json['action'] == 'receive':
                return await self.receive(request, request_json)
            elif request_json['action'] == 'send':
                return await self.send(request, request_json)
            elif request_json['action'] == 'account_representative_set':
                return await self.account_representative_set(request, request_json)
            elif request_json['action'] == 'password_change':
                return await self.password_change(request, request_json)
            elif request_json['action'] == 'password_enter':
                return await self.password_enter(request, request_json)
            elif request_json['action'] == 'password_valid':
                return await self.password_valid(request, request_json)
            elif request_json['action'] == 'wallet_representative_set':
                return await self.wallet_representative_set(request, request_json)
            elif request_json['action'] == 'wallet_add':
                return await self.wallet_add(request, request_json)
            elif request_json['action'] == 'wallet_lock':
                return await self.wallet_lock(request, request_json)
            elif request_json['action'] == 'wallet_locked':
                return await self.wallet_locked(request, request_json)
            elif request_json['action'] == 'wallet_balances':
                return await self.wallet_balances(request, request_json)
            elif request_json['action'] == 'wallet_frontiers':
                return await self.wallet_frontiers(request, request_json)
            elif request_json['action'] == 'wallet_pending':
                return await self.wallet_pending(request, request_json)
            elif request_json['action'] == 'wallet_destroy':
                return await self.wallet_destroy(request, request_json)
            elif request_json['action'] == 'wallet_change_seed':
                return await self.wallet_change_seed(request, request_json)
            elif request_json['action'] == 'wallet_contains':
                return await self.wallet_contains(request, request_json)
            elif request_json['action'] == 'wallet_representative':
                return await self.wallet_representative(request, request_json)
            elif request_json['action'] == 'wallet_info':
                return await self.wallet_info(request, request_json)
            elif request_json['action'] == 'receive_all':
                return await self.receive_all(request, request_json)
            elif request_json['action'] in ['account_move', 'account_remove', 'receive_minimum', 'receive_minimum_set', 'search_pending', 'search_pending_all', 'wallet_add_watch', 'wallet_export', 'wallet_history', 'wallet_ledger', 'wallet_republish', 'wallet_work_get', 'work_get', 'work_set']:
                # Prevent unimplemented wallet RPCs from going to the node directly
                return self.json_response(
                    data = {
                        'error': 'not_implemented'
                    }
                )

            # Proxy other requests to the node
            resp_json = await RPCClient.instance().make_request(request_json)
            return self.json_response(
                data = resp_json
            )

        return self.generic_error()

    async def wallet_create(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'seed' in request_json:
            if not Validators.is_valid_block_hash(request_json['seed']):
                return self.json_response(
                    data = {'error': 'Invalid seed'}
                )
            new_seed = request_json['seed']
        else:
            new_seed = RandomUtil.generate_seed()
        async with in_transaction() as conn:
            wallet = Wallet(
                seed=new_seed
            )
            await wallet.save(using_db=conn)
            await wallet.account_create(using_db=conn)
        return self.json_response(
            data = {
                'wallet': str(wallet.id)
            }
        )

    async def account_create(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'wallet' not in request_json:
            return self.generic_error()

        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Create account
        async with in_transaction() as conn:
            account = await wallet.account_create(using_db=conn)
        return self.json_response(
            data = {
                'account': account
            }
        )

    async def accounts_create(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'wallet' not in request_json or 'count' not in request_json or not isinstance(request_json['count'], int):
            return self.generic_error()

        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Create account
        async with in_transaction() as conn:
            accounts = await wallet.accounts_create(count=request_json['count'], using_db=conn)
        return self.json_response(
            data = {
                'accounts': accounts
            }
        )

    async def account_list(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'wallet' not in request_json:
            return self.generic_error()
        elif 'count' in request_json and isinstance(request_json['count'], int):
            count = request_json['count']
        else:
            count = 1000

        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        return self.json_response(
            data = {'accounts': [a.address for a in await wallet.accounts.all().limit(count)]}
        )

    async def receive(self, request: web.Request, request_json: dict):
        """RPC receive"""
        if 'wallet' not in request_json or 'account' not in request_json or 'block' not in request_json:
            return self.generic_error()
        elif not Validators.is_valid_address(request_json['account']):
            return self.json_response(
                data={'error': 'Invalid address'}
            )
        elif not Validators.is_valid_block_hash(request_json['block']):
            return self.json_response(
                data={'error': 'Invalid block'}
            )

        work = request_json['work'] if 'work' in request_json else None

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Retrieve account on wallet
        account = await wallet.get_account(request_json['account'])
        if account is None:
            return self.json_response(
                data={'error': 'Account not found'}
            )

        # Try to receive block
        wallet = WalletUtil(account, wallet)
        try:
            response = await wallet.receive(request_json['block'], work=work)
        except BlockNotFound:
            return self.json_response(
                data={'error': 'Block not found'}
            )
        except WorkFailed:
            return self.json_response(
                data={'error': 'Failed to generate work'}
            )
        except ProcessFailed:
            return self.json_response(
                data={'error': 'RPC Process failed'}
            )

        if response is None:
            return self.json_response(
                data={'error': 'Unable to receive block'}
            )

        return self.json_response(
            data=response
        )

    async def send(self, request: web.Request, request_json: dict):
        """RPC send"""
        if 'wallet' not in request_json or 'source' not in request_json or 'destination' not in request_json or 'amount' not in request_json:
            return self.generic_error()
        elif not Validators.is_valid_address(request_json['source']):
            return self.json_response(
                data={'error': 'Invalid source'}
            )
        elif not Validators.is_valid_address(request_json['destination']):
            return self.json_response(
                data={'error': 'Invalid destination'}
            )

        id = request_json['id'] if 'id' in request_json else None
        work = request_json['work'] if 'work' in request_json else None

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Retrieve account on wallet
        account = await wallet.get_account(request_json['source'])
        if account is None:
            return self.json_response(
                data={'error': 'Account not found'}
            )

        # Try to create and publish send block
        wallet = WalletUtil(account, wallet)
        try:
            resp = await wallet.send(int(request_json['amount']), request_json['destination'], id=id, work=work)
        except AccountNotFound:
            return self.json_response(
                data={'error': 'Account not found'}
            )
        except BlockNotFound:
            return self.json_response(
                data={'error': 'Block not found'}
            )
        except WorkFailed:
            return self.json_response(
                data={'error': 'Failed to generate work'}
            )
        except ProcessFailed:
            return self.json_response(
                data={'error': 'RPC Process failed'}
            )
        except InsufficientBalance:
            return self.json_response(
                data={'error': 'insufficient balance'}
            )

        if resp is None:
            return self.json_response(
                data={'error': 'Unable to create send block'}
            )

        return self.json_response(
            data=resp
        )

    async def account_representative_set(self, request: web.Request, request_json: dict):
        """RPC account_representative_set"""
        if 'wallet' not in request_json or 'account' not in request_json or 'representative' not in request_json:
            return self.generic_error()
        elif not Validators.is_valid_address(request_json['account']):
            return self.json_response(
                data={'error': 'Invalid account'}
            )
        elif not Validators.is_valid_address(request_json['representative']):
            return self.json_response(
                data={'error': 'Invalid representative'}
            )

        work = request_json['work'] if 'work' in request_json else None

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Retrieve account on wallet
        account = await wallet.get_account(request_json['account'])
        if account is None:
            return self.json_response(
                data={'error': 'Account not found'}
            )

        # Try to create and publish CHANGE block
        wallet = WalletUtil(account, wallet)
        try:
            resp = await wallet.representative_set(request_json['representative'], work=work)
        except AccountNotFound:
            return self.json_response(
                data={'error': 'Account not found'}
            )
        except WorkFailed:
            return self.json_response(
                data={'error': 'Failed to generate work'}
            )
        except ProcessFailed:
            return self.json_response(
                data={'error': 'RPC Process failed'}
            )

        if resp is None:
            return self.json_response(
                data={'error': 'Unable to create change block'}
            )

        return self.json_response(
            data=resp
        )

    async def password_change(self, request: web.Request, request_json: dict):
        """RPC password_change"""
        if 'wallet' not in request_json or 'password' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Encrypt
        await wallet.encrypt_wallet(request_json['password'])

        return self.json_response(
            data={'changed': '1'}
        )

    async def password_enter(self, request: web.Request, request_json: dict):
        """RPC password_enter"""
        if 'wallet' not in request_json or 'password' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
            return self.json_response(
                data={
                    'error': 'wallet not locked'
                }
            )
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked as w:
            wallet = w.wallet

        try:
            await wallet.unlock_wallet(request_json['password'])
        except DecryptionError:
            return self.json_response(
                data={'valid': '0'}
            )

        return self.json_response(
            data={'valid': '1'}
        )

    async def password_valid(self, request: web.Request, request_json: dict):
        """RPC password_valid"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
            if not wallet.encrypted:
                return self.json_response(
                    data={
                        'error': 'wallet not locked'
                    }
                )
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={'valid': '0'}
            )

        return self.json_response(
            data={'valid': '1'}
        )

    async def wallet_representative_set(self, request: web.Request, request_json: dict):
        """RPC wallet_representative_set"""
        if 'wallet' not in request_json or 'representative' not in request_json or ('update_existing_accounts' in request_json and not isinstance(request_json['update_existing_accounts'], bool)):
            return self.generic_error()
        elif not Validators.is_valid_address(request_json['representative']):
            return self.json_response(
                data={'error': 'Invalid address'}
            )

        update_existing = False
        if 'update_existing_accounts' in request_json:
            update_existing = request_json['update_existing_accounts']

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        wallet.representative = request_json['representative']
        await wallet.save(update_fields=['representative'])

        if update_existing:
            await wallet.bulk_representative_update(request_json['representative'])

        return self.json_response(
            data={'set': '1'}
        )

    async def wallet_add(self, request: web.Request, request_json: dict):
        """RPC wallet_add"""
        if 'wallet' not in request_json or 'key' not in request_json:
            return self.generic_error()
        elif not Validators.is_valid_block_hash(request_json['key']):
            return self.json_response(
                data={'error': 'Invalid key'}
            )

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Add account
        try:
            address = await wallet.adhoc_account_create(request_json['key'])
        except AccountAlreadyExists:
            return self.json_response(
                data={
                    'error': 'account already exists'
                }
            )

        return self.json_response(
            data={'account':address}
        )

    async def wallet_lock(self, request: web.Request, request_json: dict):
        """RPC wallet_lock"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked as we:
            wallet = we.wallet

        await wallet.lock_wallet()

        return self.json_response(
            data={'locked':'1'}
        )

    async def wallet_locked(self, request: web.Request, request_json: dict):
        """RPC wallet_locked"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={'locked': '1'}
            )

        return self.json_response(
            data={'locked':'0'}
        )

    async def wallet_balances(self, request: web.Request, request_json: dict):
        """RPC wallet_balances"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        return self.json_response(
            data=await RPCClient.instance().accounts_balances([a.address for a in await wallet.get_all_accounts()])
        )

    async def wallet_pending(self, request: web.Request, request_json: dict):
        """RPC wallet_pending"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        return self.json_response(
            data=await RPCClient.instance().accounts_pending([a.address for a in await wallet.get_all_accounts()])
        )

    async def wallet_frontiers(self, request: web.Request, request_json: dict):
        """RPC wallet_frontiers"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        return self.json_response(
            data=await RPCClient.instance().accounts_frontiers([a.address for a in await wallet.get_all_accounts()])
        )

    async def wallet_destroy(self, request: web.Request, request_json: dict):
        """RPC wallet_destroy"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )
        await wallet.delete()
        return self.json_response(
            data={'destroyed': '1'}
        )

    async def wallet_change_seed(self, request: web.Request, request_json: dict):
        """RPC wallet_change_seed"""
        if 'wallet' not in request_json or 'seed' not in request_json:
            return self.generic_error()
        elif not Validators.is_valid_block_hash(request_json['seed']):
            return self.json_response(
                data={'error': 'Invalid seed'}
            )

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Reset password
        if wallet.encrypted:
            await wallet.encrypt_wallet('')

        # Change key
        await wallet.change_seed(request_json['seed'])

        # Get newest account
        newest = await wallet.get_newest_account()

        return self.json_response(
            data={
                "success": "",
                "last_restored_account": newest.address,
                "restored_count": newest.account_index + 1
            }
        )

    async def wallet_contains(self, request: web.Request, request_json: dict):
        """RPC wallet_contains"""
        if 'wallet' not in request_json or 'account' not in request_json:
            return self.generic_error()
        elif not Validators.is_valid_address(request_json['account']):
            return self.json_response(
                data={'error': 'Invalid account'}
            )

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        exists = (await wallet.get_account(request_json['account'])) is not None

        return self.json_response(
            data={'exists': '1' if exists else '0'}
        )

    async def wallet_representative(self, request: web.Request, request_json: dict):
        """RPC wallet_representative"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        if wallet.representative is None:
            wallet.representative = config.Config.instance().get_random_rep()

        return self.json_response(
            data={'representative': wallet.representative}
        )

    async def wallet_info(self, request: web.Request, request_json: dict):
        """RPC wallet_info"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        # Get balances
        balance = 0
        pending_bal = 0
        balance_json = await RPCClient.instance().accounts_balances([a.address for a in await wallet.get_all_accounts()])

        if 'balances' not in balance_json:
            return self.json_response(
                data={'error': 'failed to retrieve balances'}
            )
    
        for k, v in balance_json['balances'].items():
            balance += int(v['balance'])
            pending_bal += int(v['pending'])

        det_count = await wallet.accounts.all().count()
        adhoc_count = await wallet.accounts.all().count()
        det_index = (await wallet.get_newest_account()).account_index

        return self.json_response(
            data={
                'balance': balance,
                'pending': pending_bal,
                'accounts_count': det_count + adhoc_count,
                'adhoc_count': adhoc_count,
                'deterministic_count': det_count,
                'deterministic_index': det_index
            }
        )

    async def receive_all(self, request: web.Request, request_json: dict):
        """receive every pending block in wallet"""
        if 'wallet' not in request_json:
            return self.generic_error()

        # Retrieve wallet
        try:
            wallet = await Wallet.get_wallet(request_json['wallet'])
        except WalletNotFound:
            return self.json_response(
                data={
                    'error': 'wallet not found'
                }
            )
        except WalletLocked:
            return self.json_response(
                data={
                    'error': 'wallet locked'
                }
            )

        balance_json = await RPCClient.instance().accounts_balances([a.address for a in await wallet.get_all_accounts()])

        if 'balances' not in balance_json:
            return self.json_response(
                data={'error': 'failed to retrieve balances'}
            )
    
        received_count = 0
        for k, v in balance_json['balances'].items():
            if int(v['pending']) > 0:
                wallet_util = WalletUtil(await wallet.get_account(k), wallet)
                received_count += await wallet_util.receive_all()

        return self.json_response(
            data={
                'received': received_count
            }
        )

    async def block_arrival_handler(self, data: dict):
        """invoked when we receive a new block"""
        log.server_logger.debug("Received Callback")
        is_send = False
        if 'block' in data and 'subtype' in data['block'] and data['block']['subtype'] == 'send':
            is_send = True
        elif 'is_send' in data and (data['is_send'] == 'true' or data['is_send']):
            is_send = True
        if is_send:
            # Ignore receive_minimum
            if config.Config.instance().receive_minimum > int(data['amount']):
                return
            # Determine if the recipient is one of ours
            destination = data['block']['link_as_account']
            acct = await Account.filter(address=destination).prefetch_related('wallet').first()
            if acct is None:
                acct = await AdHocAccount.filter(address=destination).prefetch_related('wallet').first()
                if acct is None:
                    return
            log.server_logger.debug(f"Auto receiving {data['hash']} for {destination}")
            wu = WalletUtil(acct, acct.wallet)
            try:
                await wu.receive(data['hash'])
            except Exception:
                log.server_logger.debug(f"Failed to receive {data['hash']}")

    async def start(self):
        """Start the server"""
        runner = web.AppRunner(self.app, access_log = None if not config.Config.instance().debug else log.server_logger)
        await runner.setup()
        site = web.TCPSite(runner, self.host, self.port)
        tasks = []
        tasks.append(site.start())
        # Websocket
        if self.websocket:
            await self.websocket.setup()
            tasks.append(self.websocket.loop())
        await asyncio.wait(tasks)
