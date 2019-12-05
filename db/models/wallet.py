from aiohttp import log
from aioredis_lock import RedisLock, LockTimeoutError
from db.redis import RedisDB
from tortoise.functions import Max
from tortoise.models import Model
from tortoise import fields
from typing import List

import nanopy
import db.models.account as acct

class Wallet(Model):
    id = fields.UUIDField(pk=True)
    seed = fields.CharField(max_length=64, unique=True)

    class Meta:
        table = 'wallets'

    async def account_create(self, using_db=None) -> str:
        """Create an account on this seed and return the created account"""
        async with RedisLock(
            await RedisDB.instance().get_redis(),
            key=f"pippin:{str(self.id)}:account_create",
            timeout=30,
            wait_timeout=30
        ):
            account = await acct.Account.filter(wallet=self).annotate(max_index=Max("account_index")).order_by('-account_index').first()
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
        async with RedisLock(
            await RedisDB.instance().get_redis(),
            key=f"pippin:{str(self.id)}:account_create",
            timeout=30,
            wait_timeout=30
        ):
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