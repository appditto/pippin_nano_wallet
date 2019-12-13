from pippin.util.random import RandomUtil

import unittest

class TestRandomUtil(unittest.TestCase):
    def test_seed_generate(self):
        """Test seed generation"""
        self.assertTrue(len(RandomUtil.generate_seed()) == 64)
        for i in range(5):
            self.assertNotEqual(RandomUtil.generate_seed(), RandomUtil.generate_seed())