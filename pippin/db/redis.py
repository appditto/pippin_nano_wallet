import asyncio
import os

import aioredis
import aioredlock

class RedisDB(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'RedisDB':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.redis = None
            cls.lock_manager = None
            cls.redis_host = os.getenv('REDIS_HOST', '127.0.0.1')
            cls.redis_port = int(os.getenv('REDIS_PORT', '6379'))
            cls.redis_db = int(os.getenv('REDIS_DB', 0))
        return cls._instance

    @classmethod
    async def close(cls):
        if hasattr(cls, 'redis') and cls.redis is not None:
            cls.redis.close()
            await cls.redis.wait_closed()
        if cls._instance is not None:
            cls._instance = None

    @classmethod
    async def get_redis(cls) -> aioredis.Redis:
        if cls.redis is not None:
            return cls.redis
        cls.redis = await aioredis.create_redis_pool((cls.redis_host, cls.redis_port), db=cls.redis_db, encoding='utf-8', minsize=1, maxsize=5)
        return cls.redis

    @classmethod
    async def get_lock_manager(cls) -> aioredlock.Aioredlock:
        if cls.lock_manager is not None:
            return cls.lock_manager
        cls.lock_manager = aioredlock.Aioredlock([await cls.get_redis()], lock_timeout=300, retry_count=3, retry_delay_min=30, retry_delay_max=90)
        return cls.lock_manager

    async def set(self, key: str, value: str, expires: int = 0):
        """Basic redis SET"""
        # Add a prefix
        key = f"pippin:{key}"
        redis = await self.get_redis()
        await redis.set(key, value, expire=expires)

    async def get(self, key: str):
        """Redis GET"""
        # Add a prefix
        key = f"pippin:{key}"
        redis = await self.get_redis()
        return await redis.get(key)

    async def delete(self, key: str):
        """Redis DELETE"""
        key = f"pippin:{key}"
        await self._delete(key)

    async def _delete(self, key: str):
        """Redis DELETE"""
        redis = await self.get_redis()
        await redis.delete(key)

    async def exists(self, key: str):
        """See if a key exists"""
        key = f"pippin:{key}"
        redis = await self.get_redis()
        return (await redis.get(key)) is not None