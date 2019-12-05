from aiohttp import web, log
from db.models.wallet import Wallet
from tortoise.transactions import in_transaction
from util.random import RandomUtil

import config
import datetime
import logging
import rapidjson as json

class PippinServer(object):
    """API for wallet requests"""
    def __init__(self, host: str, port: int):
        self.app = web.Application()
        self.app.add_routes([
            web.post('/api', self.gateway)
        ])
        self.host = host
        self.port = port

    async def stop(self):
        await self.app.shutdown()

    async def generic_error(self):
        """The node returns this generic error when the request is bad"""
        return web.json_response(
            data={
                'error':"Unable to parse json"
            },
            dumps=json.dumps
        )

    async def gateway(self, request: web.Request):
        """Gateway route to mimic nano's API of specifying action in a string"""
        request_json = await request.json(loads=json.loads)      
        if 'action' in request_json:
            if request_json['action'] == 'wallet_create':
                return await self.wallet_create(request, request_json)
            elif request_json['action'] == 'account_create':
                return await self.account_create(request, request_json)
            elif request_json['action'] == 'accounts_create':
                return await self.accounts_create(request, request_json)
            elif request_json['action'] == 'account_list':
                return await self.account_list(request, request_json)

        return self.generic_error()

    async def wallet_create(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        new_seed = RandomUtil.generate_seed()
        async with in_transaction() as conn:
            wallet = Wallet(
                seed=new_seed
            )
            await wallet.save(using_db=conn)
            await wallet.account_create(using_db=conn)
        return web.json_response(
            data = {
                'wallet': str(wallet.id)
            },
            dumps = json.dumps
        )

    async def account_create(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'wallet' not in request_json:
            return self.generic_error()

        wallet = await Wallet.filter(id=request_json['wallet']).first()
        if wallet is None:
            return web.json_response(
                data={
                    'error': 'wallet not found'
                },
                dumps=json.dumps
            )

        # Create account
        async with in_transaction() as conn:
            account = await wallet.account_create(using_db=conn)
        return web.json_response(
            data = {
                'account': account
            },
            dumps = json.dumps
        )

    async def accounts_create(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'wallet' not in request_json or 'count' not in request_json or not isinstance(request_json['count'], int):
            return self.generic_error()

        wallet = await Wallet.filter(id=request_json['wallet']).first()
        if wallet is None:
            return web.json_response(
                data={
                    'error': 'wallet not found'
                },
                dumps=json.dumps
            )

        # Create account
        async with in_transaction() as conn:
            accounts = await wallet.accounts_create(count=request_json['count'], using_db=conn)
        return web.json_response(
            data = {
                'accounts': accounts
            },
            dumps = json.dumps
        )

    async def account_list(self, request: web.Request, request_json: dict):
        """Route for creating new wallet"""
        if 'wallet' not in request_json:
            return self.generic_error()
        elif 'count' in request_json and isinstance(request_json['acount'], int):
            count = request_json['count']
        else:
            count = 1000

        wallet = await Wallet.filter(id=request_json['wallet']).first()
        if wallet is None:
            return web.json_response(
                data={
                    'error': 'wallet not found'
                },
                dumps=json.dumps
            )

        return web.json_response(
            data = {'accounts': [a.address for a in await wallet.accounts.all().limit(count)]},
            dumps=json.dumps
        )

    async def start(self):
        """Start the server"""
        runner = web.AppRunner(self.app, access_log = None if not config.Config.instance().debug else log.server_logger)
        await runner.setup()
        site = web.TCPSite(runner, self.host, self.port)
        await site.start()