from tortoise.models import Model
from tortoise import fields

class Wallet(Model):
    id = fields.UUIDField(pk=True)
    seed = fields.CharField(max_length=64)

    class Meta:
        table = 'wallets'