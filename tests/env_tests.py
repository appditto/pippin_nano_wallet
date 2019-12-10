import os
import unittest

from util.env import Env

class EnvTest(unittest.TestCase):
    def test_banano_env(self):
        os.environ['BANANO'] = '1'
        self.assertTrue(Env.banano())
        del os.environ['BANANO']
        self.assertFalse(Env.banano())