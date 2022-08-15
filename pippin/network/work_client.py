import socket
from logging import config
from typing import List

import aiohttp
import asyncio
import pippin.config as config
import nanolib
import nanopy
import os
import rapidjson as json

from pippin.db.redis import RedisDB
from pippin.util.nano_util import NanoUtil

from aiographql.client import (GraphQLClient, GraphQLRequest, GraphQLResponse)

BPOW_URL = os.getenv("BPOW_URL", "https://boompow.banano.cc/graphql")


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
            cls.bpow_key = os.getenv('BPOW_KEY', None)
            if cls.bpow_key is not None:
                cls.bpow_client = GraphQLClient(
                    endpoint=BPOW_URL,
                    headers={"Authorization": f"{cls.bpow_key}"},
                )

        return cls._instance

    @classmethod
    async def close(cls):
        if hasattr(cls, 'session') and cls.session is not None:
            await cls.session.close()
        if cls._instance is not None:
            cls._instance = None

    async def make_request(self, url: str, req_json: dict):
        async with self.session.post(url, json=req_json, timeout=300) as resp:
            return await resp.json()

    async def work_generate(self, hash: str, difficulty: str, blockAward: bool = True) -> str:
        work_generate = {
            'action': 'work_generate',
            'hash': hash,
            'difficulty': difficulty
        }

        # Build work_generate requests
        tasks = []
        for p in self.work_urls:
            tasks.append(self.make_request(p, work_generate))

        # Use BPoW if applicable
        if self.bpow_client is not None:
            multiplier = int(nanolib.work.derive_work_multiplier(
                difficulty, base_difficulty="fffffe0000000000"))
            if multiplier < 1:
                multiplier = 1
            request = GraphQLRequest(
                validate=False,
                query="""
                    mutation($hash:String!, $difficultyMultiplier: Int!, $blockAward: Boolean) {
                        workGenerate(input:{hash:$hash, difficultyMultiplier:$difficultyMultiplier, blockAward:$blockAward})
                    }
                """,
                variables={
                    "hash": hash, "difficultyMultiplier":  multiplier, "blockAward": blockAward}
            )
            tasks.append(self.bpow_client.query(request=request))

        # Do it locally if no peers or if peers have been failing
        if await RedisDB.instance().exists("work_failure") or (len(self.work_urls) == 0 and self.bpow_client is None):
            tasks.append(
                NanoUtil.instance().work_generate(hash, difficulty=difficulty)
            )

        # Post work_generate to all peers simultaneously
        final_result = None
        while len(tasks) > 0:
            # Fire all tasks simultaneously and fire when first one is completed
            done, pending = await asyncio.wait(tasks, return_when=asyncio.FIRST_COMPLETED, timeout=100)
            # Set to True if any tasks completed before the timeout
            has_done = len(done) > 0
            try:
                for task in done:
                    result = task.result()
                    # Some tasks return different types of responses, e.g. DPoW is different than a normal work peer
                    if result is None:
                        aiohttp.log.server_logger.info(
                            "work_generate task returned None")
                    if isinstance(result, GraphQLResponse):
                        if len(result.errors) > 0:
                            result = None
                        else:
                            result = {"work": result.data["workGenerate"]}
                    elif isinstance(result, list):
                        result = json.loads(result[1])
                    elif isinstance(result, str):
                        result = {'work': result}
                    if result is not None and 'work' in result:
                        cancel_json = {
                            'action': 'work_cancel',
                            'hash': hash
                        }
                        for p in self.work_urls:
                            asyncio.ensure_future(
                                self.make_request(p, cancel_json))
                        final_result = result['work']
                    elif result is not None and 'error' in result:
                        aiohttp.log.server_logger.info(
                            f'work_generate task returned error {result["error"]}')
            except Exception:
                aiohttp.log.server_logger.exception(
                    "work_generate task raised Exception")
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
