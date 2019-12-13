from concurrent.futures.process import ProcessPoolExecutor
from concurrent.futures.thread import ThreadPoolExecutor

import asyncio
import pippin.config as config
import functools
import nanopy

class NanoUtil(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls, max_work_processes: int = 1, max_sign_threads: int = 1) -> 'NanoUtil':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            work_processes = max(config.Config.instance().max_work_processes, 0)
            if work_processes > 0:
                cls.process_pool = ProcessPoolExecutor(max_workers=work_processes)
            else:
                cls.process_pool = None
            cls.thread_pool = ThreadPoolExecutor(max_workers=max(config.Config.instance().max_sign_threads, 1))
            cls.thread_pool
        return cls._instance

    @classmethod
    async def close(cls):
        if hasattr(cls, 'process_pool') and cls.process_pool is not None:
            cls.process_pool.shutdown()
        if hasattr(cls, 'thread_pool') and cls.thread_pool is not None:
            cls.thread_pool.shutdown()
        if cls._instance is not None:
            cls._instance = None

    async def work_generate(self, hash: str, difficulty: str = None) -> str:
        """Run work_generate in ProcessPool"""
        if self.process_pool is None:
            raise WorkDisabled()

        difficulty = difficulty if difficulty is not None else nanopy.work_difficulty
        result = await asyncio.get_event_loop().run_in_executor(
            self.process_pool,
            functools.partial(
                nanopy.work_generate,
                hash,
                difficulty
            )
        )
        return result

    async def sign_block(self, private_key: str, block: dict) -> str:
        """Sign block in ThreadPool"""
        # We want to sign blocks in a ThreadPool to prevent blocking our main event loop
        result = await asyncio.get_event_loop().run_in_executor(
            self.thread_pool,
            functools.partial(
                nanopy.sign,
                private_key,
                block
            )
        )
        return result

class WorkDisabled(Exception):
    pass