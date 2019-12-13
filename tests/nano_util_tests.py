from aiounittest import async_test
from pippin.util.nano_util import NanoUtil

import asyncio
import unittest

class TestNanoUtil(unittest.TestCase):
    @async_test
    async def test_work_generate(self):
        """Test work generation"""
        hash = "5E5B7C8F97BDA8B90FAA243050D99647F80C25EB4A07E7247114CBB129BDA188"
        difficulty = "ff00000000000000"
        result = await NanoUtil.instance().work_generate(hash, difficulty=difficulty)
        self.assertTrue(len(result) > 8)
        await NanoUtil.close()