from tortoise.models import Model
from tortoise import fields
from pippin.model.secrets import SeedStorage

class AdHocAccount(Model):
    wallet  = fields.ForeignKeyField('db.Wallet', on_delete=fields.CASCADE, related_name='adhoc_accounts', index=True)
    address = fields.CharField(max_length=65)
    private_key = fields.CharField(max_length=128)
    work = fields.BooleanField(default=True)
    created_at = fields.DatetimeField(auto_now_add=True)

    class Meta:
        table = 'adhoc_accounts'
        unique_together = ('address', 'private_key')

    def private_key_get(self) -> str:
        """Return account private key in hex"""
        priv_key = SeedStorage.instance().get_decrypted_seed(f"{self.wallet_id}:{self.address}")
        return priv_key if priv_key is not None else self.private_key