from tortoise.models import Model
from tortoise import fields
from enum import Enum, unique

class Payment(Model):

    def asdict(self):
        return {'address': self.address, 'business_memo_id': self.business_memo_id, 'is_paid': self.is_paid}

    id = fields.IntField(pk=True)
    created_at = fields.DatetimeField(auto_now_add=True)
    address = fields.CharField(max_length=65, index=True, unique=True)
    business_memo_id = fields.CharField(max_length=40)
    is_paid = fields.BooleanField()
    amount = fields.CharField(max_length=64) # in raw

    class Meta:
        table = 'payments'