import json
import os
import pathlib

from aiohttp.test_utils import AioHTTPTestCase, unittest_run_loop
from tortoise import Tortoise

from pippin.config import Config
from pippin.db.models.wallet import Wallet
from pippin.db.tortoise_config import DBConfig
from pippin.db.redis import RedisDB
from pippin.server.pippin_server import PippinServer
from pippin.util.utils import Utils


class PippinServerTest(AioHTTPTestCase):
    async def get_application(self):
        self.server = PippinServer('127.0.0.1', 9999)
        return self.server.app

    @classmethod
    def setUpClass(cls):
        return super().setUpClass()

    async def setUpAsync(self):
        if 'BANANO' in os.environ:
            del os.environ['BANANO']
        await DBConfig(mock=True).init_db()
        await RedisDB.instance().get_redis()

    async def tearDownAsync(self):
        if 'BANANO' in os.environ:
            del os.environ['BANANO']
        await Tortoise.close_connections()
        await RedisDB.close()
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db')))
        except FileNotFoundError:
            pass
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db-wal')))
        except FileNotFoundError:
            pass
        try:
            os.remove(Utils.get_project_root().joinpath(pathlib.PurePath('mock.db-shm')))
        except FileNotFoundError:
            pass

    async def json_request(self, data: str):
        return await self.client.request(
            'POST',
            '/',
            data=data,
            headers={
                'Content-Type': 'application/json'
            }
        )

    async def create_test_wallet(self):
        body = json.dumps(
            {
                'action': 'wallet_create',
                'seed':'C273AB6E1D8121C5DA0B99DD44CF9AA29D51C40B009ACB9410CA1649E28170E8'
            }
        )
        response = await self.json_request(body)
        json_resp = await response.json()
        return json_resp['wallet']


    @unittest_run_loop
    async def test_bad_request(self):
        body = json.dumps(
            {
                'sdasd': 'asdadas'
            }
        )
        response = await self.json_request(body)
        json_resp = await response.json()
        self.assertTrue('error' in json_resp)

    @unittest_run_loop
    async def test_wallet_create(self):
        body = json.dumps(
            {
                'action': 'wallet_create',
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)

        json_resp = await response.json()
        self.assertTrue('wallet' in json_resp)

    @unittest_run_loop
    async def test_wallet_create_with_seed(self):
        wallet_id = await self.create_test_wallet()

        # Make sure we can get wallet
        wallet: Wallet = await Wallet.get_wallet(wallet_id)
        self.assertEqual(wallet.seed, 'C273AB6E1D8121C5DA0B99DD44CF9AA29D51C40B009ACB9410CA1649E28170E8')
        account = await wallet.get_account('nano_3n5jcd764a1t5duoc455zhrz6ernn534c9onrao4tghhfpw8y6m65og66qtc')
        self.assertEqual(account.address, 'nano_3n5jcd764a1t5duoc455zhrz6ernn534c9onrao4tghhfpw8y6m65og66qtc')
        self.assertEqual(account.account_index, 0)

    @unittest_run_loop
    async def test_account_create(self):
        wallet_id = await self.create_test_wallet()

        body = json.dumps(
            {
                'action': 'account_create',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('account' in json_resp)
        # Equal to index 1 account
        self.assertEqual(json_resp['account'], 'nano_3yfyxm9aeyawqwfcd6c6zd36ph877xwwpxeqppn5wp6rzeifrnpi1p49fhr9')

    @unittest_run_loop
    async def test_accounts_create(self):
        wallet_id = await self.create_test_wallet()

        body = json.dumps(
            {
                'action': 'accounts_create',
                'wallet':wallet_id,
                'count':5
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('accounts' in json_resp)
        # Equal to index 1-5 accounts
        self.assertEqual(json_resp['accounts'][0], 'nano_3yfyxm9aeyawqwfcd6c6zd36ph877xwwpxeqppn5wp6rzeifrnpi1p49fhr9')
        self.assertEqual(json_resp['accounts'][1], 'nano_1pdqbxoe81uzdd3u7ft91u49ztdkgdgfoeak7jmysbia7msigbctunmq9php')
        self.assertEqual(json_resp['accounts'][2], 'nano_3a34k3bgc147r6emhac6jxgzquqif695hi8xqypowdru3eum81je3no4iwty')
        self.assertEqual(json_resp['accounts'][3], 'nano_3im7mxo6717mhc54t9ypyber981y8perupjap5p9ysf4wymf5e4a64ugnjax')
        self.assertEqual(json_resp['accounts'][4], 'nano_3cptr3ispmjte8ky1e15imt6py7yppwp1h4r7kfyxy1z3xjbjeujmca4g3m1')

    @unittest_run_loop
    async def test_account_list(self):
        wallet_id = await self.create_test_wallet()
        # Create accounts
        body = json.dumps(
            {
                'action': 'accounts_create',
                'wallet':wallet_id,
                'count':5
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)

        # Test list
        body = json.dumps(
            {
                'action': 'account_list',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('accounts' in json_resp)
        # Equal to index 0-5 accounts
        self.assertEqual(json_resp['accounts'][0], 'nano_3n5jcd764a1t5duoc455zhrz6ernn534c9onrao4tghhfpw8y6m65og66qtc')
        self.assertEqual(json_resp['accounts'][1], 'nano_3yfyxm9aeyawqwfcd6c6zd36ph877xwwpxeqppn5wp6rzeifrnpi1p49fhr9')
        self.assertEqual(json_resp['accounts'][2], 'nano_1pdqbxoe81uzdd3u7ft91u49ztdkgdgfoeak7jmysbia7msigbctunmq9php')
        self.assertEqual(json_resp['accounts'][3], 'nano_3a34k3bgc147r6emhac6jxgzquqif695hi8xqypowdru3eum81je3no4iwty')
        self.assertEqual(json_resp['accounts'][4], 'nano_3im7mxo6717mhc54t9ypyber981y8perupjap5p9ysf4wymf5e4a64ugnjax')
        self.assertEqual(json_resp['accounts'][5], 'nano_3cptr3ispmjte8ky1e15imt6py7yppwp1h4r7kfyxy1z3xjbjeujmca4g3m1')

    @unittest_run_loop
    async def test_account_list(self):
        wallet_id = await self.create_test_wallet()
        # Create accounts
        body = json.dumps(
            {
                'action': 'accounts_create',
                'wallet':wallet_id,
                'count':5
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)

        # Test list
        body = json.dumps(
            {
                'action': 'account_list',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('accounts' in json_resp)
        # Equal to index 0-5 accounts
        self.assertEqual(json_resp['accounts'][0], 'nano_3n5jcd764a1t5duoc455zhrz6ernn534c9onrao4tghhfpw8y6m65og66qtc')
        self.assertEqual(json_resp['accounts'][1], 'nano_3yfyxm9aeyawqwfcd6c6zd36ph877xwwpxeqppn5wp6rzeifrnpi1p49fhr9')
        self.assertEqual(json_resp['accounts'][2], 'nano_1pdqbxoe81uzdd3u7ft91u49ztdkgdgfoeak7jmysbia7msigbctunmq9php')
        self.assertEqual(json_resp['accounts'][3], 'nano_3a34k3bgc147r6emhac6jxgzquqif695hi8xqypowdru3eum81je3no4iwty')
        self.assertEqual(json_resp['accounts'][4], 'nano_3im7mxo6717mhc54t9ypyber981y8perupjap5p9ysf4wymf5e4a64ugnjax')
        self.assertEqual(json_resp['accounts'][5], 'nano_3cptr3ispmjte8ky1e15imt6py7yppwp1h4r7kfyxy1z3xjbjeujmca4g3m1')

    @unittest_run_loop
    async def test_password_change(self):
        wallet_id = await self.create_test_wallet()
        # Set password
        body = json.dumps(
            {
                'action': 'password_change',
                'wallet':wallet_id,
                'password': 'abcd123'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('changed' in json_resp)
        self.assertEqual(json_resp['changed'], '1')

        # Ensure we get wallet locked error when changing again
        # Set password
        body = json.dumps(
            {
                'action': 'password_change',
                'wallet':wallet_id,
                'password': 'def567'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('error' in json_resp)
        self.assertEqual(json_resp['error'].lower(), 'wallet locked')

    @unittest_run_loop
    async def test_password_enter(self):
        wallet_id = await self.create_test_wallet()
        # Set password
        body = json.dumps(
            {
                'action': 'password_change',
                'wallet':wallet_id,
                'password': 'abcd123'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('changed' in json_resp)
        self.assertEqual(json_resp['changed'], '1')

        # Ensure bad password doesn't unlock
        body = json.dumps(
            {
                'action': 'password_enter',
                'wallet':wallet_id,
                'password': 'def567'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('valid' in json_resp)
        self.assertEqual(json_resp['valid'], '0')
        body = json.dumps(
            {
                'action': 'password_change',
                'wallet':wallet_id,
                'password': 'def567'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('error' in json_resp)
        self.assertEqual(json_resp['error'].lower(), 'wallet locked')

        # Test good password
        body = json.dumps(
            {
                'action': 'password_enter',
                'wallet':wallet_id,
                'password': 'abcd123'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('valid' in json_resp)
        self.assertEqual(json_resp['valid'], '1')
        body = json.dumps(
            {
                'action': 'password_change',
                'wallet':wallet_id,
                'password': 'def567'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('changed' in json_resp)
        self.assertEqual(json_resp['changed'], '1')

    @unittest_run_loop
    async def test_wallet_lock_locked(self):
        wallet_id = await self.create_test_wallet()
        # Set password
        body = json.dumps(
            {
                'action': 'password_change',
                'wallet':wallet_id,
                'password': 'abcd123'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('changed' in json_resp)
        self.assertEqual(json_resp['changed'], '1')

        # Check locked stats
        body = json.dumps(
            {
                'action': 'wallet_locked',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('locked' in json_resp)
        self.assertEqual(json_resp['locked'], '1') # Locked

        # Unlock
        body = json.dumps(
            {
                'action': 'password_enter',
                'wallet':wallet_id,
                'password': 'abcd123'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('valid' in json_resp)
        self.assertEqual(json_resp['valid'], '1')

        # Check locked stats
        body = json.dumps(
            {
                'action': 'wallet_locked',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('locked' in json_resp)
        self.assertEqual(json_resp['locked'], '0') # Unlocked

        # Lock
        body = json.dumps(
            {
                'action': 'wallet_lock',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('locked' in json_resp)
        self.assertEqual(json_resp['locked'], '1')

        # Check locked stats
        body = json.dumps(
            {
                'action': 'wallet_locked',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('locked' in json_resp)
        self.assertEqual(json_resp['locked'], '1') # Locked

    @unittest_run_loop
    async def test_wallet_representative_set(self):
        wallet_id = await self.create_test_wallet()
        # Set rep
        body = json.dumps(
            {
                'action': 'wallet_representative_set',
                'wallet':wallet_id,
                'representative': 'nano_1he11o1darknessmyo1dfriend11111111111111111111111111ki7nzfx6'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('set' in json_resp)
        self.assertEqual(json_resp['set'], '1')

        # Ensure we get right representative when querying
        body = json.dumps(
            {
                'action': 'wallet_representative',
                'wallet':wallet_id
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('representative' in json_resp)
        self.assertEqual(json_resp['representative'], 'nano_1he11o1darknessmyo1dfriend11111111111111111111111111ki7nzfx6')

        # Ensure we can't set a bad rep
        body = json.dumps(
            {
                'action': 'wallet_representative_set',
                'wallet':wallet_id,
                'representative': 'nano_1he11o1darknessmyo1dfriend11111111111111111111111111ki7nzfx5'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('error' in json_resp)
        self.assertEqual(json_resp['error'].lower(), 'invalid address')

    @unittest_run_loop
    async def test_adhoc_account(self):
        wallet_id = await self.create_test_wallet()
        # Set rep
        body = json.dumps(
            {
                'action': 'wallet_add',
                'wallet':wallet_id,
                'key': '11A780CD68BF64AE703C9AA7C138E5D6F917EBB60C55D50BE8BD6AFBA65066F8'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('account' in json_resp)
        self.assertEqual(json_resp['account'], 'nano_1997zr1zatrca8giqmyyz45a7q8qqgqkxdntg8uashgrqe5odmz93q5u954d')

    @unittest_run_loop
    async def test_wallet_change_seed(self):
        wallet_id = await self.create_test_wallet()
        # Change seed
        body = json.dumps(
            {
                'action': 'wallet_change_seed',
                'wallet':wallet_id,
                'seed': 'D8BF40997C092631640856E2CD47FCA1BB39FE678ED34665580E8FF61FA1C049'
            }
        )
        response = await self.json_request(body)
        self.assertEqual(response.status, 200)
        json_resp = await response.json()
        self.assertTrue('success' in json_resp)
        self.assertTrue('last_restored_account' in json_resp)
        self.assertTrue('restored_count' in json_resp)
        self.assertEqual(json_resp['last_restored_account'], 'nano_1997zr1zatrca8giqmyyz45a7q8qqgqkxdntg8uashgrqe5odmz93q5u954d')
        self.assertEqual(json_resp['restored_count'], 1)
