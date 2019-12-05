import secrets
import string

class RandomUtil(object):

    @staticmethod
    def generate_seed() -> str:
        """Generate a random seed and return it"""
        seed = "".join([secrets.choice(string.hexdigits) for i in range(64)]).upper()
        return seed