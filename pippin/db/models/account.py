from tortoise.models import Model
from tortoise import fields

import nanopy

class Account(Model):
    wallet  = fields.ForeignKeyField('db.Wallet', on_delete=fields.CASCADE, related_name='accounts', index=True)
    address = fields.CharField(max_length=65)
    account_index = fields.IntField()
    work = fields.BooleanField(default=True)
    created_at = fields.DatetimeField(auto_now_add=True)

    class Meta:
        table = 'accounts'
        unique_together = ('address', 'account_index')

    def private_key(self, seed: str) -> str:
        """Return account private key in hex"""
        return nanopy.deterministic_key(seed, self.account_index)[0]