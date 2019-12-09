from tortoise.models import Model
from tortoise import fields

class AdHocAccount(Model):
    wallet  = fields.ForeignKeyField('db.Wallet', related_name='adhoc_accounts', index=True)
    address = fields.CharField(max_length=65)
    private_key = fields.CharField(max_length=128)
    work = fields.BooleanField(default=True)
    created_at = fields.DatetimeField(auto_now_add=True)

    class Meta:
        table = 'adhoc_accounts'
        unique_together = ('address', 'private_key')