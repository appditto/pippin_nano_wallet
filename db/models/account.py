from tortoise.models import Model
from tortoise import fields

class Account(Model):
    wallets  = fields.ForeignKeyField('db.User', related_name='accounts', index=True)
    address = fields.CharField(max_length=65)
    index = fields.IntField()

    class Meta:
        table = 'accounts'