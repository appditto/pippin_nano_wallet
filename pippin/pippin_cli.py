import pathlib
from dotenv import load_dotenv
load_dotenv()
from pippin.util.utils import Utils
load_dotenv(dotenv_path=Utils.get_project_root().joinpath(pathlib.PurePath('.env')))

import argparse
import asyncio
import getpass

from pippin.db.models.wallet import Wallet, WalletLocked, WalletNotFound
from pippin.db.tortoise_config import DBConfig
from tortoise import Tortoise
from tortoise.transactions import in_transaction
from pippin.util.crypt import AESCrypt, DecryptionError
from pippin.util.random import RandomUtil
from pippin.util.validators import Validators
from pippin.version import __version__

from pippin.config import Config
import os

# Set and patch nanopy
import nanopy
nanopy.account_prefix = 'ban_' if Config.instance().banano else 'nano_'
if Config.instance().banano:
    nanopy.standard_exponent = 29
    nanopy.work_difficulty = 'fffffe0000000000'

parser = argparse.ArgumentParser(description=f'Pippin v{__version__}')
subparsers = parser.add_subparsers(title='available commands', dest='command')

wallet_parser = subparsers.add_parser('wallet_list')

wallet_create_parser = subparsers.add_parser('wallet_create')
wallet_create_parser.add_argument('--seed', type=str, help='Seed for wallet (optional)', required=False)

wallet_change_seed_parser = subparsers.add_parser('wallet_change_seed')
wallet_change_seed_parser.add_argument('--wallet', type=str, help='ID of wallet to change seed for', required=True)
wallet_change_seed_parser.add_argument('--seed', type=str, help='New seed for wallet (optional)', required=False)
wallet_change_seed_parser.add_argument('--encrypt', action='store_true', help='If specified, will get prompted for a password to encrypt the wallet', default=False)

wallet_view_seed_parser = subparsers.add_parser('wallet_view_seed')
wallet_view_seed_parser.add_argument('--wallet', type=str, help='Wallet ID', required=True)
wallet_view_seed_parser.add_argument('--password', type=str, help='Password needed to decrypt wallet (if encrypted)', required=False)
wallet_view_seed_parser.add_argument('--all-keys', action='store_true', help='Also show all of the wallet address and keys', default=False)

account_create_parser = subparsers.add_parser('account_create')
account_create_parser.add_argument('--wallet', type=str, help='Wallet ID', required=True)
account_create_parser.add_argument('--key', type=str, help='AdHoc Account Key', required=False)
account_create_parser.add_argument('--count', type=int, help='Number of accounts to create (min: 1)', required=False)

wallet_destroy_parser = subparsers.add_parser('wallet_destroy')
wallet_destroy_parser.add_argument('--wallet', type=str, help='Wallet ID', required=True)

repget_parser = subparsers.add_parser('wallet_representative_get')
repget_parser.add_argument('--wallet', type=str, help='Wallet ID', required=True)

repset_parser = subparsers.add_parser('wallet_representative_set')
repset_parser.add_argument('--wallet', type=str, help='Wallet ID', required=True)
repset_parser.add_argument('--representative', type=str, help='New Wallet Representative', required=True)
repset_parser.add_argument('--update-existing', action='store_true', help='Update existing accounts', default=False)

options = parser.parse_args()

async def wallet_list():
    wallets = await Wallet.all().prefetch_related('accounts', 'adhoc_accounts')
    if len(wallets) == 0:
        print("There aren't any wallets")
        return

    for w in wallets:
        print(f"ID:{w.id}")
        print("Accounts:")
        for a in w.accounts:
            print(a.address)

async def wallet_create(seed):
    async with in_transaction() as conn:
        wallet = Wallet(
            seed=RandomUtil.generate_seed() if seed is None else seed
        )
        await wallet.save(using_db=conn)
        new_acct = await wallet.account_create(using_db=conn)
    print(f"Wallet created, ID: {wallet.id}\nFirst account: {new_acct}")

async def wallet_change_seed(wallet_id: str, seed: str, password: str) -> str:
    encrypt = False
    old_password = None
    if len(password) > 0:
        encrypt = True

    # Retrieve wallet
    try:
        wallet = await Wallet.get_wallet(wallet_id)
    except WalletNotFound:
        print(f"No wallet found with ID: {wallet_id}")
        exit(1)
    except WalletLocked as wl:
        wallet = wl.wallet
        while True:
            try:
                npass = getpass.getpass(prompt='Enter current password:')
                crypt = AESCrypt(npass)
                try:
                    decrypted = crypt.decrypt(wallet.seed)
                    async with in_transaction() as conn:
                        wallet.seed = decrypted
                        wallet.encrypted = False
                        await wallet.save(using_db=conn, update_fields=['seed', 'encrypted'])
                        for a in await wallet.adhoc_accounts.all():
                            a.private_key = crypt.decrypt(a.private_key)
                            await a.save(using_db=conn, update_fields=['private_key'])
                    old_password = npass
                    break
                except DecryptionError:
                    print("**Invalid password**")
            except KeyboardInterrupt:
                break
                exit(0)

    # Change key
    await wallet.change_seed(seed)

    # Encrypt if necessary
    if encrypt:
        await wallet.encrypt_wallet(password)

    # Get newest account
    newest = await wallet.get_newest_account()

    print(f"Seed changed for wallet {wallet.id}\nFirst account: {newest.address}")

async def wallet_view_seed(wallet_id: str, password: str, all_keys: bool) -> str:
    # Retrieve wallet
    crypt = None
    try:
        wallet = await Wallet.get_wallet(wallet_id)
    except WalletNotFound:
        print(f"No wallet found with ID: {wallet_id}")
        exit(1)
    except WalletLocked as wl:
        wallet = None
        if password is not None:
            crypt = AESCrypt(password)
            try:
                decrypted = crypt.decrypt(wl.wallet.seed)
                wallet = wl.wallet
                wallet.seed = decrypted
            except DecryptionError:
                pass
        if wallet is None:
            while True:
                try:
                    npass = getpass.getpass(prompt='Enter current password:')
                    crypt = AESCrypt(npass)
                    try:
                        decrypted = crypt.decrypt(wl.wallet.seed)
                        wallet = wl.wallet
                        wallet.seed = decrypted
                    except DecryptionError:
                        print("**Invalid password**")
                except KeyboardInterrupt:
                    break
                    exit(0)

    print(f"Seed: {wallet.seed}")
    if all_keys:
        for a in await wallet.accounts.all():
            print(f"Addr: {a.address} PrivKey: {nanopy.deterministic_key(wallet.seed, index=a.account_index)[0].upper()}")
    else:
        print(f"AdHoc accounts:")
        for a in await wallet.adhoc_accounts.all():
            if not wallet.encrypted:
                print(f"Addr: {a.address} PrivKey: {a.private_key.upper()}")
            else:
                print(f"Addr: {a.address} PrivKey: {crypt.decrypt(a.private_key)}")

async def account_create(wallet_id: str, key: str, count: int = 1) -> str:
    # Retrieve wallet
    crypt = None
    password=None
    if count is None:
        count = 1
    try:
        wallet = await Wallet.get_wallet(wallet_id)
    except WalletNotFound:
        print(f"No wallet found with ID: {wallet_id}")
        exit(1)
    except WalletLocked as wl:
        wallet = wl.wallet
        if key is not None:
            while True:
                try:
                    npass = getpass.getpass(prompt='Enter current password to encrypt ad-hoc key:')
                    crypt = AESCrypt(npass)
                    try:
                        decrypted = crypt.decrypt(wl.wallet.seed)
                        wallet = wl.wallet
                        wallet.seed = decrypted
                        password=npass
                    except DecryptionError:
                        print("**Invalid password**")
                except KeyboardInterrupt:
                    break
                    exit(0)

    if key is None:
        if count == 1:
            a = await wallet.account_create()
            print(f"account: {a}")
        else:
            async with in_transaction() as conn:
                ass = await wallet.accounts_create(count=count)
                for a in ass:
                    print(f"account: {a}")
    else:
        a = await wallet.adhoc_account_create(key, password=password)
        print(f"account: {a}")

async def wallet_destroy(wallet_id: str):
    # Retrieve wallet
    try:
        wallet = await Wallet.get_wallet(wallet_id)
    except WalletNotFound:
        print(f"No wallet found with ID: {wallet_id}")
        exit(1)
    except WalletLocked as wl:
        wallet = wl.wallet

    await wallet.delete()
    print("Wallet destroyed")

async def wallet_representative_get(wallet_id: str):
    # Retrieve wallet
    try:
        wallet = await Wallet.get_wallet(wallet_id)
    except WalletNotFound:
        print(f"No wallet found with ID: {wallet_id}")
        exit(1)
    except WalletLocked as wl:
        wallet = wl.wallet

    if wallet.representative is None:
        print("Representative not set")
    else:
        print(f"Wallet representative: {wallet.representative}")

async def wallet_representative_set(wallet_id: str, rep: str, update_existing: bool = False):
    # Retrieve wallet
    # Retrieve wallet
    crypt = None
    password=None
    if not Validators.is_valid_address(rep):
        print("Invalid representative")
        exit(1)
    try:
        wallet = await Wallet.get_wallet(wallet_id)
    except WalletNotFound:
        print(f"No wallet found with ID: {wallet_id}")
        exit(1)
    except WalletLocked as wl:
        wallet = wl.wallet
        if update_existing:
            while True:
                try:
                    npass = getpass.getpass(prompt='Enter current password to decrypt wallet:')
                    crypt = AESCrypt(npass)
                    try:
                        decrypted = crypt.decrypt(wl.wallet.seed)
                        wallet = wl.wallet
                        wallet.seed = decrypted
                        password=npass
                    except DecryptionError:
                        print("**Invalid password**")
                except KeyboardInterrupt:
                    break
                    exit(0)

    wallet.representative = rep
    await wallet.save(update_fields=['representative'])
    await wallet.bulk_representative_update(rep)
    print(f"Representative changed")

def main():
    loop = asyncio.new_event_loop()
    try:
        loop.run_until_complete(DBConfig().init_db())
        if options.command == 'wallet_list':
            loop.run_until_complete(wallet_list())
        elif options.command == 'wallet_create':
            if options.seed is not None:
                if not Validators.is_valid_block_hash(options.seed):
                    print("Invalid seed specified")
                    exit(1)
            loop.run_until_complete(wallet_create(options.seed))
        elif options.command == 'wallet_change_seed':
            if options.seed is not None:
                if not Validators.is_valid_block_hash(options.seed):
                    print("Invalid seed specified")
                    exit(1)
            else:
                while True:
                    try:
                        options.seed = getpass.getpass(prompt='Enter new wallet seed:')
                        if Validators.is_valid_block_hash(options.seed):
                            break
                        print("**Invalid seed**, should be a 64-character hex string")
                    except KeyboardInterrupt:
                        break
                        exit(0)
            password = ''
            if options.encrypt:
                while True:
                    try:
                        password = getpass.getpass(prompt='Enter password to encrypt wallet:')
                        if password.strip() == '':
                            print("**Bad password** - cannot be blanke")
                        break
                    except KeyboardInterrupt:
                        break
                        exit(0)
            loop.run_until_complete(wallet_change_seed(options.wallet, options.seed, password))
        elif options.command == 'wallet_view_seed':
            loop.run_until_complete(wallet_view_seed(options.wallet, options.password, options.all_keys))
        elif options.command == 'account_create':
            if options.key is not None:
                if not Validators.is_valid_block_hash(options.key):
                    print("Invalid Private Key")
                    exit(0)
            elif options.key is not None and options.count is not None:
                print("You can only specify one: --key or --count")
                print("--count can only be used for deterministic accounts")
            elif options.count is not None:
                if options.count < 1:
                    print("Count needs to be at least 1...")
            loop.run_until_complete(account_create(options.wallet, options.key, options.count))
        elif options.command == 'wallet_destroy':
            loop.run_until_complete(wallet_destroy(options.wallet))
        elif options.command == 'wallet_representative_get':
            loop.run_until_complete(wallet_representative_get(options.wallet))
        elif options.command == 'wallet_representative_set':
            loop.run_until_complete(wallet_representative_set(options.wallet, options.representatives, update_existing=options.update_existing))
        else:
            parser.print_help()
    except Exception as e:
        print(str(e))
        raise e
    finally:
        loop.run_until_complete(Tortoise.close_connections())
        loop.close()

if __name__ == "__main__":
    main()