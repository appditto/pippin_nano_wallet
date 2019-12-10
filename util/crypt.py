import base64
import hashlib

from Crypto import Random
from Crypto.Cipher import AES

class AESCrypt():

    def __init__(self, secret_key: str):
        self.key = hashlib.sha256(secret_key.encode('utf-8')).digest()
        self.BS = 16
        self.salt = '61606982'

    pad = lambda self, s: s + (self.BS - len(s) % self.BS) * chr(self.BS - len(s) % self.BS)
    unpad = lambda self, s: s[0:-s[-1]]

    def encrypt(self, value: str) -> str:
        value = self.salt + value
        value = self.pad(value)
        iv = Random.new().read(AES.block_size)
        cipher = AES.new(self.key, AES.MODE_CBC, iv)
        return base64.b64encode(iv + cipher.encrypt(value.encode('latin-1'))).decode('latin-1')

    def decrypt(self, encrypted: str) -> str:
        encrypted = base64.b64decode(encrypted.encode('latin-1'))
        iv = encrypted[:16]
        cipher = AES.new(self.key, AES.MODE_CBC, iv)
        decrypted = self.unpad(cipher.decrypt(encrypted[16:])).decode('latin-1')
        # Test if decrypt was successful
        if len(decrypted) < len(self.salt):
            raise DecryptionError()
        elif decrypted[:len(self.salt)] != self.salt:
            raise DecryptionError()
        return decrypted[len(self.salt):]

class DecryptionError(Exception):
    pass
