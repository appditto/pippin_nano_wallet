from tortoise.models import Model
from tortoise import fields

class Account(Model):
    wallet  = fields.ForeignKeyField('db.Wallet', related_name='accounts', index=True)
    address = fields.CharField(max_length=65)
    account_index = fields.IntField()

    class Meta:
        table = 'accounts'
        unique_together = ('address', 'account_index')