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

    def get_decrypted_seed(self, encrypted: str) -> str:
        if encrypted in self.seeds:
            return self.seeds[encrypted]
        return None

    def set_decrypted_seed(self, encrypted: str, decrypted: str) -> str:
        self.seeds[encrypted] = decrypted

    def contains_encrypted(self, encrypted: str) -> bool:
        return encrypted in self.seeds