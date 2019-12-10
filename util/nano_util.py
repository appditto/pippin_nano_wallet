from concurrent.futures.process import ProcessPoolExecutor
from concurrent.futures.thread import ThreadPoolExecutor

import asyncio
import nanopy

class NanoUtil(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls, max_work_processes: int = 1, max_sign_threads: int = 1) -> 'NanoUtil':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.process_pool = ProcessPoolExecutor(max_workers=max_work_processes)
            cls.thread_pool = ThreadPoolExecutor(max_workers=max_sign_threads)
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