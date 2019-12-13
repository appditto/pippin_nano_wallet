from pippin.util.crypt import AESCrypt, DecryptionError

import unittest

class TestAESCrypt(unittest.TestCase):
    def setUp(self):
        self.crypt1 = AESCrypt('mypassword')
        self.crypt2 = AESCrypt('someotherpassword')

    def test_encrypt_decrypt(self):
        """Test encryption and decryption"""
        some_string = "Some Random String to Encrypt"
        encrypted = self.crypt1.encrypt(some_string)
        self.assertNotEqual(encrypted, some_string)
        self.assertEqual(self.crypt1.decrypt(encrypted), some_string)
        with self.assertRaises(DecryptionError) as exc:
            self.crypt2.decrypt(encrypted)