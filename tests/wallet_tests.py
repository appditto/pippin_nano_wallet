from pippin.util.wallet import WalletUtil
from pippin.util.nano_util import NanoUtil
from pippin.util.utils import Utils
from tortoise import Tortoise
from pippin.db.tortoise_config import DBConfig
from pippin.db.redis import RedisDB
from pippin.db.models.account import Account
from pippin.db.models.adhoc_account import AdHocAccount
from pippin.db.models.wallet import Wallet

import os
import asyncio
import unittest
import uuid
import pathlib

class TestWalletUtil(unittest.TestCase):
    @classmethod
    async def setUpClassAsync(cls):
        if 'BANANO' in os.environ:
            del os.environ['BANANO']
        await DBConfig(mock=True).init_db()
        cls.wallet = Wallet(
            seed="7474F694061FB3E5813986AEC8A65340B5DEDB4DF94E394CB44489BEA6B21FCD",
            representative='nano_1oa4sipdm679m7emx9npmoua14etkrj1e85g3ujkxmn4ti1aifbu7gx1zrgr',
            encrypted=False
        )
        await cls.wallet.save()
        cls.account = Account(
            wallet=cls.wallet,
            address='nano_19zdjp6tfhqzcag9y3w499h36nr16ks6gdsfkawzgomrbxs54xaybmsyamza',
            account_index=0
        )
        await cls.account.save()
        cls.adhoc_account = AdHocAccount(
            wallet=cls.wallet,
            address='nano_1oa4sipdm679m7emx9npmoua14etkrj1e85g3ujkxmn4ti1aifbu7gx1zrgr',
            private_key='86A3D926AB6BEBAA678C13823D7A92A97CAFAFD277EBF4B54C42C8BB9806EAEE'
        )
        await cls.adhoc_account.save()
        cls.wallet_util = WalletUtil(cls.account, cls.wallet)
        cls.wallet_util_adhoc = WalletUtil(cls.adhoc_account, cls.wallet)

    @classmethod
    def removeMockDB(cls):
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db')))
        except FileNotFoundError:
            print(str(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db'))))
            pass
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db-wal')))
        except FileNotFoundError:
            pass
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db-shm')))
        except FileNotFoundError:
            pass

    @classmethod
    def setUpClass(cls):
        cls.removeMockDB()
        cls.loop = asyncio.new_event_loop()
        cls.loop.run_until_complete(cls.setUpClassAsync())

    @classmethod
    def tearDownClass(cls):
        cls.loop.run_until_complete(Tortoise.close_connections())
        cls.loop.run_until_complete(NanoUtil.close())
        cls.loop.close()
        cls.removeMockDB()

    def test_get_representative(self):
        self.assertEqual(self.wallet_util.get_representative(), 'nano_1oa4sipdm679m7emx9npmoua14etkrj1e85g3ujkxmn4ti1aifbu7gx1zrgr')

    def test_adhoc(self):
        self.assertFalse(self.wallet_util.adhoc())
        self.assertTrue(self.wallet_util_adhoc.adhoc())

    def test_private_key_derivation(self):
        """Ensure  private keys can be derived correctly"""
        self.assertEqual('DD21B99F1A92D7315BF8592F632B05EC77555DFD6A0B2D99C560537CAFC9A13E', self.wallet_util.private_key().upper())
        self.assertEqual('86A3D926AB6BEBAA678C13823D7A92A97CAFAFD277EBF4B54C42C8BB9806EAEE', self.wallet_util_adhoc.private_key().upper())

