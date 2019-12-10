# Pippin

Pippin is a production-ready, high-performance developer wallet for [Nano](https://nano.org) and [BANANO](https://banano.cc). Pippin's API is a drop-in replacement for the Nano developer wallet that is built in to the Nano node software.

## About Pippin

Pippin is written in Python. It achieves high performance across the board using libraries such as [asyncio](https://docs.python.org/3/library/asyncio.html), [uvloop](https://github.com/MagicStack/uvloop), [aiohttp](https://aiohttp.readthedocs.io/en/stable/), [asyncpg](https://github.com/MagicStack/asyncpg)/[aiosqlite](https://github.com/jreese/aiosqlite)/[aiomysql](https://github.com/aio-libs/aiomysql), and [rapidjson](https://rapidjson.org/).

For block signing and work generation, Pippin uses [nanopy](https://github.com/npy0/nanopy), which is a high-performance library that utilizes C-bindings for blake2b and ed25519.

## Benefits of Pippin

The Nano developer wallet (aka "node wallet") is not recommended for production use. One of the goals of Pippin is to provide a production-ready external key management that can be used by developers who are using Nano.

Pippin is the first drop-in replacement for the Nano developer wallet. It's incredibly easy to transition to Pippin if you are already using the Nano developer wallet.

- Pippin is independent of the node. You can use Pippin with any public RPC, so you don't have to run your own node
- Pippin is extremely fast and lightweight
- Pippin supports encrypted secret keys
- Pippin supports multiple database backends (SQLite, PostgreSQL, and MySQL)

Pippin can be used by exchanges, games, payment processors, tip bots, faucets, casinos, and a lot more.

## How Pippin Works

Pippin provides an API that mimics the [Nano Wallet RPC Protocol](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs)

Every wallet-related RPC gets intercepted by Pippin and handled internally. It builds the blocks and signs them using locally-stored keys, it uses a node to publish the blocks.

Every non-wallet related RPC gets proxied to the publishing node. Which means you can make all of your RPC requests directly to Pippin whether they are wallet-related or not.

## API Documentation

Recommended reference is the [NANO RPC documentation](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs), Pippin's APIs are mostly identical.

## Known Differences: Pippin vs NANO Node Wallet

**Pippin supports send indempotency with the `id` send parameter, just like the Nano node it is not required but highly recommended**

**Enhanced Behavior**

- `account_list` accepts a `count` parameter that defaults to 1000
- `receive_all` **new** RPC action: accepts a `wallet` parameter - receives all pendings in the entire wallet. This RPC respects the `receive_minimum` setting

**Degraded Behavior**

- `account_create` does not accept an index

**Fuzzy**

- `wallet_change_seed` decrypts the wallet, if it's encrypted

**Missing/Not Implemented**

- `account_move`
- `account_remove`
- `receive_minimum` - Receive minimum can be set in `config.yaml`
- `receive_minimum_set`
- `wallet_add_watch` - Not certain what this even does
- `wallet_history` - Would be more efficient if NANO supported `accounts_history`
- `search_pending` - Pippin has `receive_all` which should be used to receive all pendings
- `search_pending_all`
- `wallet_export`
- `wallet_ledger` - Pippin doesn't store the ledger
- `wallet_republish` - Same as above, pippin only rebroadcasts `send` RPCs when a duplicate ID is used
- `wallet_work_get` - I'm not really sure what these work RPCs do, pippin doesn't store any information about work
- `work_get`
- `work_set`

## CLI Documentation

Unlike the API, the CLI is entirely different from the Nano node CLI.

You can see all available options with the CLI using `./pippin --help`

## Setting up Pippin

### Requirements

- Python 3.7 or newer
- GCC, for MacOS and Linux
- libb2 (blake2b)

On MacOS, with homebrew:

```
brew install gcc@9 python libb2
```

On Linux, debian-based systems:

```
sudo apt install build-essential python3.7 python3.7-dev libb2-dev
```

### Installing python dependencies

MacOS:

```
CC=/usr/local/bin/gcc-9 python3 -m pip install -U -r requirements.txt
```

Linux:

```
python3.7 -m pip install -r requirements.txt
```

## Feature requests

Notice an API that's missing a feature or not behaving the same as nano's APIs?

Open a bug report/feature request on the [issues page](https://github.com/bbedward/pippin_nano_wallet/issues)