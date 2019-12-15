from typing import List

import nanopy
import uuid
from aiohttp import log

from aioredlock import LockError
from tortoise import fields
from tortoise.functions import Max
from tortoise.models import Model
from tortoise.transactions import in_transaction

import pippin.db.models.account as acct
import pippin.db.models.adhoc_account as adhoc_acct
from pippin.db.redis import RedisDB
from pippin.model.secrets import SeedStorage
from pippin.util.crypt import AESCrypt
from pippin.util.wallet import WalletUtil
from pippin.network.rpc_client import AccountNotFound


class WalletNotFound(Exception):
    pass

class WalletLocked(Exception):
    def __init__(self, wallet):
        self.wallet = wallet
    pass

class AccountAlreadyExists(Exception):
    pass

class Wallet(Model):
    id = fields.UUIDField(pk=True)
    seed = fields.CharField(max_length=128, unique=True)
    representative = fields.CharField(max_length=65, null=True)
    encrypted = fields.BooleanField(default=False)
    work = fields.BooleanField(default=True)
    created_at = fields.DatetimeField(auto_now_add=True)

    class Meta:
        table = 'wallets'

    @staticmethod
    async def get_wallet(id: str) -> 'Wallet':
        """Get wallet with ID, raise WalletNotFound if not found, WalletLocked if encrypted"""
        wallet = await Wallet.filter(id=uuid.UUID(id)).first()
        if wallet is None:
            raise WalletNotFound()
        elif wallet.encrypted:
            decrypted = SeedStorage.instance().get_decrypted_seed(wallet.id)
            if decrypted is None:
                raise WalletLocked(wallet)
            wallet.seed = decrypted
        return wallet

    async def change_seed(self, seed: str):
        async with in_transaction() as conn:
            self.seed = seed
            for a in await self.accounts.all():
                a.address = nanopy.deterministic_key(self.seed, index=a.account_index)[2]
                await a.save(using_db=conn, update_fields=['address'])
            await self.save(using_db=conn, update_fields=['seed'])

    async def encrypt_wallet(self, password: str):
        """Encrypt wallet seed with password"""
        async with in_transaction() as conn:
            # If password is empty string then decrypted wallet
            if len(password.strip()) == 0:
                self.encrypted = False
                for a in await self.adhoc_accounts.all():
                    decrypted = SeedStorage.instnace().get_decrypted_seed(f"{self.id}:{a.address}")
                    if decrypted is not None:
                        a.private_key = decrypted
                    await a.save(using_db=conn, update_fields=['private_key'])
            else:
                crypt = AESCrypt(password)
                encrypted = crypt.encrypt(self.seed)
                self.seed = encrypted
                self.encrypted = True
                for a in await self.adhoc_accounts.all():
                    a.private_key = crypt.encrypt(a.private_key_get())
                    await a.save(using_db=conn, update_fields=['private_key'])            
            await self.save(using_db=conn, update_fields=['seed', 'encrypted'])

    async def unlock_wallet(self, password: str):
        """Unlock wallet with given password, raise DecryptionError if invalid password"""
        crypt = AESCrypt(password)
        decrypted = crypt.decrypt(self.seed)
        # Store decrypted wallet in memory
        SeedStorage.instance().set_decrypted_seed(self.id, decrypted)
        # Decrypt any ad-hoc accounts
        for a in await self.adhoc_accounts.all():
            SeedStorage.instance().set_decrypted_seed(f"{self.id}:{a.address}", crypt.decrypt(a.private_key))

    async def lock_wallet(self):
        """Lock wallet and remove all decrypted seeds from memory"""
        SeedStorage.instance().remove(self.id)
        for a in await self.adhoc_accounts.all():
            await SeedStorage.instance().remove(f"{self.id}:{a.address}")

    async def is_locked(self) -> bool:
        """Determine whether this wallet is locked or not"""
        return SeedStorage.instance().contains_encrypted(self.id) 

    async def bulk_representative_update(self, rep: str):
        """Set all account representatives to rep"""
        for a in await self.accounts.all():
            w = WalletUtil(a, self)
            try:
                await w.representative_set(rep, only_if_different=True)
            except AccountNotFound:
                pass

    async def adhoc_account_create(self, key: str, password: str = None) -> str:
        """Add an adhoc private key to the wallet, raise AccountAlreadyExists if it already exists"""
        pubkey = nanopy.ed25519_blake2b.publickey(bytes.fromhex(key)).hex()
        address = nanopy.account_get(pubkey)
        # See if address already exists
        a = await self.accounts.filter(address=address).first()
        if a is None:
            a = await self.adhoc_accounts.filter(address=address).first()
        if a is not None:
            raise AccountAlreadyExists(a)
        # Create it
        crypt = None
        if password is not None:
            crypt = AESCrypt(password)
        a = adhoc_acct.AdHocAccount(
            wallet=self,
            private_key=crypt.encrypt(key) if crypt is not None else key,
            address=address
        )
        await a.save()
        return address

    async def get_newest_account(self) -> acct.Account:
        """Get account with highest index beloning to this wallet"""
        return await acct.Account.filter(wallet=self).annotate(max_index=Max("account_index")).order_by('-account_index').first()

    async def account_create(self, using_db=None) -> str:
        """Create an account on this seed and return the created account"""
        async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{str(self.id)}:account_create") as lock:
            account = await self.get_newest_account()
            log.server_logger.debug(f"Creating account for {self.id}")
            index = account.max_index + 1 if account is not None and account.max_index is not None and account.max_index >= 0 else 0
            private_key, public_key, address = nanopy.deterministic_key(self.seed, index=index)
            log.server_logger.debug(f"Created {address} at index {index}")
            account = acct.Account(
                wallet=self,
                account_index=index,
                address=address
            )
            await account.save(using_db=using_db)
            return address

    async def accounts_create(self, count=0, using_db=None) -> List[str]:
        """Create {count} accounts on this seed and return the created accounts"""
        count = max(1, count)
        async with await (await RedisDB.instance().get_lock_manager()).lock(f"pippin:{str(self.id)}:account_create") as lock:
            account = await acct.Account.filter(wallet=self).annotate(max_index=Max("account_index")).order_by('-account_index').first()
            log.server_logger.debug(f"Creating {count} accounts for {self.id}")
            current_index = account.max_index + 1 if account is not None and account.max_index is not None and account.max_index >= 0 else 0
            accounts = []
            for i in range(count):
                private_key, public_key, address = nanopy.deterministic_key(self.seed, index=current_index)
                accounts.append(
                    acct.Account(
                        wallet=self,
                        account_index=current_index,
                        address=address
                    )
                )
                current_index+=1
            await acct.Account.bulk_create(
                accounts, using_db=using_db
            )

            return [a.address for a in accounts]

    async def get_account(self, address: str) -> acct.Account:
        """Get an an account that begins to this wallet"""
        a = await self.accounts.filter(address=address).first()
        if a is None:
            # Check adhoc
            a = await self.adhoc_accounts.filter(address=address).first()
        return a

    async def get_all_accounts(self) -> List[acct.Account]:
        """Get all accounts belong to this wallet"""
        accounts = await self.accounts.all()
        if accounts is None:
            accounts = []
        adhoc_accounts = await self.adhoc_accounts.all()
        if adhoc_accounts is not None:
            accounts += adhoc_accounts
        return accounts