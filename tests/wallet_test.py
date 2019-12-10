from util.wallet import WalletUtil
from util.utils import Utils
from db.tortoise_config import DBConfig
from db.redis import RedisDB
from db.models.account import Account
from db.models.adhoc_account import AdHocAccount
from db.models.wallet import Wallet

import os
import asyncio
import aiounittest
import unittest
import uuid
import pathlib

class TestWalletUtil(aiounittest.AsyncTestCase):
    async def setUpAsync(self):
        if 'BANANO' in os.environ:
            del os.environ['BANANO']
        await DBConfig(mock=True).init_db()
        self.wallet = Wallet(
            seed="7474F694061FB3E5813986AEC8A65340B5DEDB4DF94E394CB44489BEA6B21FCD",
            representative='nano_1oa4sipdm679m7emx9npmoua14etkrj1e85g3ujkxmn4ti1aifbu7gx1zrgr',
            encrypted=False
        )
        await self.wallet.save()
        self.account = Account(
            wallet=self.wallet,
            address='nano_19zdjp6tfhqzcag9y3w499h36nr16ks6gdsfkawzgomrbxs54xaybmsyamza',
            account_index=0
        )
        await self.account.save()
        self.adhoc_account = AdHocAccount(
            wallet=self.wallet,
            address='nano_1oa4sipdm679m7emx9npmoua14etkrj1e85g3ujkxmn4ti1aifbu7gx1zrgr',
            private_key='86A3D926AB6BEBAA678C13823D7A92A97CAFAFD277EBF4B54C42C8BB9806EAEE'
        )
        await self.adhoc_account.save()
        self.wallet_util = WalletUtil(self.account, self.wallet, await RedisDB.instance().get_redis())
        self.wallet_util_adhoc = WalletUtil(self.adhoc_account, self.wallet, await RedisDB.instance().get_redis())

    def removeMockDB(self):
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db')))
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db-wal')))
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db-shm')))
        except FileNotFoundError:
            pass

    def setUp(self):
        self.removeMockDB()
        loop = asyncio.get_event_loop()
        loop.run_until_complete(self.setUpAsync())

    def tearDown(self):
        self.removeMockDB()

    def test_get_representative(self):
        self.assertEqual(self.wallet_util.get_representative(), 'nano_1oa4sipdm679m7emx9npmoua14etkrj1e85g3ujkxmn4ti1aifbu7gx1zrgr')

    def test_adhoc(self):
        self.assertFalse(self.wallet_util.adhoc())
        self.assertTrue(self.wallet_util_adhoc.adhoc())

    async def test_private_key_derivation(self):
        """Ensure  private keys can be derived correctly"""
        self.assertEqual('DD21B99F1A92D7315BF8592F632B05EC77555DFD6A0B2D99C560537CAFC9A13E', self.wallet_util.private_key().upper())
        self.assertEqual('86A3D926AB6BEBAA678C13823D7A92A97CAFAFD277EBF4B54C42C8BB9806EAEE', self.wallet_util_adhoc.private_key().upper())
