import socket
from logging import config
from typing import List

import aiohttp
import asyncio
import config
import nanopy
import rapidjson as json

from util.env import Env


class WorkClient(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'WorkClient':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.work_urls = config.Config.instance().work_peers
            if config.Config.instance().node_work_generate:
                cls.work_urls.append(config.Config.instance().node_url)
            cls.connector = aiohttp.TCPConnector(family=0 ,resolver=aiohttp.AsyncResolver())
            cls.session = aiohttp.ClientSession(connector=cls.connector, json_serialize=json.dumps)
        return cls._instance

    @classmethod
    async def close(cls):
        if hasattr(cls, 'session') and cls.session is not None:
            await cls.session.close()
        if cls._instance is not None:
            cls._instance = None

    async def make_request(self, url: str, req_json: dict):
        async with self.session.post(url ,json=req_json, timeout=300) as resp:
            return await resp.json()

    async def work_generate(self, hash: str, difficulty: str = nanopy.work_difficulty) -> str:
        # If no peers, do it locally
        if len(self.work_urls) == 0:
            return nanopy.work_generate(hash, difficulty=difficulty)

        work_generate = {
            'action': 'work_generate',
            'account': hash,
            'difficulty': difficulty
        }

        # Build work_generate requests
        tasks = []
        for p in self.work_urls:
            tasks.append(self.make_request(p, work_generate))

        # Post work_generate to all peers simultaneously
        while len(tasks):
            done, tasks = await asyncio.wait(tasks, return_when=asyncio.FIRST_COMPLETED, timeout=30)
            for task in done:
                try:
                    result = task.result()
                    if result is None:
                        aiohttp.log.server_logger.info("work_generate task returned None")
                        continue
                    if 'work' in result:
                        cancel_json = {
                            'action': 'work_cancel',
                            'hash': hash
                        }
                        cancel_tasks = []
                        for p in self.work_urls:
                            asyncio.ensure_future(self.make_request(p, cancel_json))
                        return result
                    elif 'error' in result:
                        aiohttp.log.server_logger.info(f'work_generate task returned error {result["error"]}')
                except Exception as exc:
                    aiohttp.log.server_logger.exception("work_generate Task raised an exception")
                    result.cancel()

        # IF we're still here fall back to nanopy
        return nanopy.work_generate(hash, difficulty)