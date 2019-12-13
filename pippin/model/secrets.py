class SeedStorage(object):
    """Store decrypted seeds in memory"""
    _instance = None

    def __init__(self):
        raise RuntimeError('Call instance() instead')

    @classmethod
    def instance(cls) -> 'SeedStorage':
        if cls._instance is None:
            cls._instance = cls.__new__(cls)
            cls.seeds = {}
        return cls._instance

    def get_decrypted_seed(self, key) -> str:
        key = str(key)
        if key in self.seeds:
            return self.seeds[key]
        return None

    def set_decrypted_seed(self, key, decrypted: str) -> str:
        key = str(key)
        self.seeds[key] = decrypted

    def contains_encrypted(self, key) -> bool:
        key = str(key)
        return key in self.seeds

    def remove(self, key):
        key = str(key)
        if key in self.seeds:
            del self.seeds[key]