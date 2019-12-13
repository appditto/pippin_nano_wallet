from pippin.util.validators import Validators

import os
import string
import unittest

class TestValidators(unittest.TestCase):
    def test_valid_block_hash(self):
        # None type
        self.assertFalse(Validators.is_valid_block_hash(None))
        # Non-hex strings
        self.assertFalse(Validators.is_valid_block_hash('1234'))
        self.assertFalse(Validators.is_valid_block_hash('Appditto LLC'))
        self.assertFalse(Validators.is_valid_block_hash('ban_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xr'))
        # Too short
        self.assertFalse(Validators.is_valid_block_hash('5E5B7C8F97BDA8B90FAA243050D99647F80C25EB4A07E7247114CBB129BDA18'))
        # Too long
        self.assertFalse(Validators.is_valid_block_hash('5E5B7C8F97BDA8B90FAA243050D99647F80C25EB4A07E7247114CBB129BDA1888'))
        # Non-hex character
        self.assertFalse(Validators.is_valid_block_hash('5E5B7C8F97BDA8B90FAA243050D99647F80C25EB4A07E7247114CBB129BDA18Z'))
        # Valid
        self.assertTrue(Validators.is_valid_block_hash('5E5B7C8F97BDA8B90FAA243050D99647F80C25EB4A07E7247114CBB129BDA188'))

    def test_valid_address(self):
        """Test address validation"""
        # Null should always be false
        self.assertFalse(Validators.is_valid_address(None))
        os.environ['BANANO'] = '1'
        # Valid
        self.assertTrue(Validators.is_valid_address('ban_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xr'))
        # Bad checksum
        self.assertFalse(Validators.is_valid_address('ban_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xa'))
        # Bad length
        self.assertFalse(Validators.is_valid_address('ban_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94x'))
        del os.environ['BANANO']
        # Valid
        self.assertTrue(Validators.is_valid_address('nano_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xr'))
        # Bad checksum
        self.assertFalse(Validators.is_valid_address('nano_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xa'))
        # Bad length
        self.assertFalse(Validators.is_valid_address('nano_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94x'))
        # Valid
        self.assertTrue(Validators.is_valid_address('xrb_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xr'))
        # Bad checksum
        self.assertFalse(Validators.is_valid_address('xrb_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94xa'))
        # Bad length
        self.assertFalse(Validators.is_valid_address('xrb_1bananobh5rat99qfgt1ptpieie5swmoth87thi74qgbfrij7dcgjiij94x'))