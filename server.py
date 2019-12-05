from aiohttp import web

import config
import datetime
import logging
import rapidjson as json

class PippinServer(object):
    """API for wallet requests"""
    def __init__(self, host: str, port: int):
        self.bot = bot
        self.app = web.Application()
        self.app.add_routes([
            web.post('/wallet_create', self.wallet_create),
            web.get('/ufw/{wallet}', self.ufw),
            web.get('/wfu/{user}', self.wfu),
            web.get('/users', self.users)
        ])
        self.logger = logging.getLogger()
        self.host = host
        self.port = port

    async def wallet_create(self, request: web.Request):
        """Route for creating new wallet"""
        request_json = await request.json(loads=json.loads)
        return web.HTTPOk()

    async def start(self):
        """Start the server"""
        runner = web.AppRunner(self.app, access_log = None if not config.Config.instance().debug else self.logger)
        await runner.setup()
        site = web.TCPSite(runner, self.host, self.port)
        await site.start()