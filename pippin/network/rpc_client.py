import aiohttp
import ipaddress
import rapidjson as json
import socket
from pippin.config import Config
from typing import List

class RPCClient(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'RPCClient':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.node_url = Config.instance().node_url
            cls.session = aiohttp.ClientSession(json_serialize=json.dumps)
        return cls._instance


    @classmethod
    async def close(cls):
        if hasattr(cls, 'session') and cls.session is not None:
            await cls.session.close()
        if cls._instance is not None:
            cls._instance = None

    async def make_request(self, req_json: dict):
        async with self.session.post(self.node_url ,json=req_json, timeout=300) as resp:
            return await resp.json()

    async def account_balance(self, account: str) -> dict:
        account_balance = {
            'action': 'account_balance',
            'account': account
        }
        respjson = await self.make_request(account_balance)
        if 'balance' in respjson:
            return respjson
        return None

    async def account_info(self, account: str) -> dict:
        info_action = {
            'action': 'account_info',
            'account': account,
            'representative': True,
            'pending': True
        }
        respjson = await self.make_request(info_action)
        if 'error' not in respjson:
            return respjson
        elif respjson['error'].lower() == 'account not found':
            raise AccountNotFound(account)
        return None

    async def block_info(self, hash: str) -> dict:
        info_action = {
            'action': 'block_info',
            'hash': hash
        }
        respjson = await self.make_request(info_action)
        if 'error' not in respjson and 'contents' in respjson:
            respjson['contents'] = json.loads(respjson['contents'])
            return respjson
        elif respjson['error'].lower() == 'block not found':
            raise BlockNotFound(hash)
        return None

    async def process(self, block: dict, subtype: str = None) -> str:
        """RPC Process, return hash if successful"""
        process_action = {
            'action': 'process',
            'json_block': False,
            'block': json.dumps(block)
        }
        if subtype is not None:
            process_action['subtype'] = subtype
        return await self.make_request(process_action)

    async def accounts_balances(self, accounts: List[str]) -> dict:
        """Return accounts_balances for accounts"""
        balances_action = {
            'action': 'accounts_balances',
            'accounts': accounts
        }
        return await self.make_request(balances_action)

    async def accounts_frontiers(self, accounts: List[str]) -> dict:
        """Return accounts_frontiers for accounts"""
        frontiers_action = {
            'action': 'accounts_frontiers',
            'accounts': accounts
        }
        return await self.make_request(frontiers_action)

    async def accounts_pending(self, accounts: List[str]) -> dict:
        """Return accounts_pending for accounts"""
        pending_action = {
            'action': 'accounts_pending',
            'accounts': accounts
        }
        return await self.make_request(pending_action)

    async def pending(self, account: str, threshold: int) -> List[str]:
        """RPC Pending"""
        pending_action = {
            'action': 'pending',
            'account': account,
            'threshold': str(threshold)
        }
        resp = await self.make_request(pending_action)
        ret = []
        if 'blocks' in resp:
            if resp['blocks'] == '':
                return []
            # Response format changes when threshold is <= 0
            # blocks: ['ABC...', 'ABC...']
            if threshold <= 0:
                return resp['blocks']
            # When threshold is > 0
            # blocks: {'ABC..', '10000', 'ABC..': '2000'}             
            for k,v in resp['blocks'].items():
                ret.append(k)
            return ret
        return None
            

    async def is_alive(self) -> bool:
        """Returns whether or not the remote node is alive"""
        # We could use 'block_count' or something simple
        # But some nodes may not whitelist actions like that,
        # so we test with an essential action like account_balance
        test_action = {
            'action': 'account_balance',
            'account': 'ban_1tipbotgges3ss8pso6xf76gsyqnb69uwcxcyhouym67z7ofefy1jz7kepoy' if Config.instance().banano else 'nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj'
        }
        respjson = await self.make_request(test_action)
        if 'error' in respjson or 'balance' in respjson:
            return True
        return False

class AccountNotFound(Exception):
    pass

class BlockNotFound(Exception):
    pass