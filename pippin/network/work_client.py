import socket
from logging import config
from typing import List

import aiohttp
import asyncio
import pippin.config as config
import nanopy
import os
import rapidjson as json

from pippin.db.redis import RedisDB
from pippin.network.dpow_websocket import ConnectionClosed, DpowClient
from pippin.util.nano_util import NanoUtil

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
            cls.session = aiohttp.ClientSession(json_serialize=json.dumps)
            cls.dpow_client = None
            cls.dpow_futures = {}
            cls.dpow_id = 1
            # Construct DPoW Client
            cls.dpow_user = os.getenv('DPOW_USER', None)
            cls.dpow_key = os.getenv('DPOW_KEY', None)
            if cls.dpow_user is not None and cls.dpow_key is not None:
                cls.dpow_client = DpowClient(
                    cls.dpow_user,
                    cls.dpow_key,
                    work_futures=cls.dpow_futures,
                    bpow=False
                )
                cls.dpow_fallback_url = 'https://dpow.nanocenter.org/service/'
            else:
                cls.dpow_user = os.getenv('BPOW_USER', None)
                cls.dpow_key = os.getenv('BPOW_KEY', None)
                if cls.dpow_user is not None and cls.dpow_key is not None:
                    cls.dpow_client = DpowClient(
                        cls.dpow_user,
                        cls.dpow_key,
                        work_futures=cls.dpow_futures,
                        bpow=True
                    )
                    cls.dpow_fallback_url = 'https://bpow.banano.cc/service/'

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

    async def work_generate(self, hash: str, difficulty: str) -> str:
        work_generate = {
            'action': 'work_generate',
            'hash': hash,
            'difficulty': difficulty
        }

        # Build work_generate requests
        tasks = []
        for p in self.work_urls:
            tasks.append(self.make_request(p, work_generate))

        # Use DPoW if applicable
        if self.dpow_client is not None:
            dpow_id = str(self.dpow_id)
            self.dpow_id += 1
            self.dpow_futures[dpow_id] = asyncio.get_event_loop().create_future()
            try:
                success = await self.dpow_client.request_work(dpow_id, hash, difficulty=difficulty)
                tasks.append(self.dpow_futures[dpow_id])
            except ConnectionClosed:
                # HTTP fallback for this request
                dp_req = {
                    "user": self.dpow_user,
                    "api_key": self.dpow_key,
                    "hash": hash,
                    "difficulty": difficulty
                }
                tasks.append(self.make_request(self.dpow_fallback_url, dp_req))

        # Do it locally if no peers or if peers have been failing
        if await RedisDB.instance().exists("work_failure") or (len(self.work_urls) == 0 and self.dpow_client is None):
            tasks.append(
                NanoUtil.instance().work_generate(hash, difficulty=difficulty)
            )

        # Post work_generate to all peers simultaneously
        final_result = None
        while len(tasks) > 0:
            # Fire all tasks simultaneously and fire when first one is completed
            done, pending = await asyncio.wait(tasks, return_when=asyncio.FIRST_COMPLETED, timeout=30)
            # Set to True if any tasks completed before the timeout
            has_done = len(done) > 0
            try:
                for task in done:
                    result = task.result()
                    # Some tasks return different types of responses, e.g. DPoW is different than a normal work peer
                    if result is None:
                        aiohttp.log.server_logger.info("work_generate task returned None")
                    if isinstance(result, list):
                        result = json.loads(result[1])
                    elif isinstance(result, str):
                        result = {'work':result}
                    if result is not None and 'work' in result:
                        cancel_json = {
                            'action': 'work_cancel',
                            'hash': hash
                        }
                        for p in self.work_urls:
                            asyncio.ensure_future(self.make_request(p, cancel_json))
                        final_result = result['work']
                    elif result is not None and 'error' in result:
                        aiohttp.log.server_logger.info(f'work_generate task returned error {result["error"]}')
            except Exception:
                aiohttp.log.server_logger.exception("work_generate task raised Exception")
            finally:
                # Either nothing finished before the timeout or something did
                # Cancel any pending tasks
                if not has_done or final_result is not None or len(pending) == 0:
                    for p in pending:
                        try:
                            p.cancel()
                        except Exception:
                            pass
                    break
                elif len(pending) > 0:
                    # wait again for pending tasks
                    tasks = pending

        if final_result is not None:
            return final_result

        # IF we're still here then all requests failed, set failure flag
        await RedisDB.instance().set(f"work_failure", "aa", expires=300)
        return await NanoUtil.instance().work_generate(hash, difficulty=difficulty)
