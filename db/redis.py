import aioredis
import asyncio

class RedisDB(object):
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'RedisDB':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.redis = None
        return cls._instance

    @classmethod
    async def close(cls):
        if hasattr(cls, 'redis') and cls.redis is not None:
            await cls.redis.wait_closed()
        if cls._instance is not None:
            cls._instance = None

    @classmethod
    async def get_redis(cls) -> aioredis.Redis:
        if cls.redis is not None:
            return cls.redis
        # TODO - we should let them override redis host/port in configuration
        cls.redis = await aioredis.create_redis_pool(('localhost', 6379), db=0, encoding='utf-8', minsize=1, maxsize=5)
        return cls.redis

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