# Pippin

[![GitHub release (latest)](https://img.shields.io/github/v/release/bbedward/pippin_nano_wallet)](https://github.com/bbedward/pippin_nano_wallet/releases) [![License](https://img.shields.io/github/license/bbedward/pippin_nano_wallet)](https://github.com/bbedward/pippin_nano_wallet/blob/master/LICENSE) [![Pipeline](https://gitlab.com/appditto/pippin_nano_wallet/badges/master/pipeline.svg)](https://gitlab.com/appditto/pippin_nano_wallet/pipelines)

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
- `auto_receive_on_send` configuration option will automatically receive pendings when doing `send`, if balance isn't high enough.

**Degraded Behavior**

- `account_create` does not accept an index

**Fuzzy Behavior**

- `wallet_change_seed` via RPC will result in the wallet no longer being encrypted. (But the wallet has to be unlocked before you can do `wallet_change_seed`).

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

Pippin has a CLI interface available, you can see available subcommands with:

`./pippin --help`

The primary goal of the CLI is key management. It's a more secure way to import a seed and backup your seed.

For example a typical flow of creating a new wallet with a specific seed might look like (add --encrypt to wallet_change_seed if you want to lock the wallet with a password):

```
% ./pippin.py wallet_create
Wallet created, ID: d897b5ec-1897-4e7e-8a90-4526f454c8de
First account: nano_31a7wzm4rayik1hthahzkekntsqz86u6dko5adg8jxueehzt5yhmhsqsuzdy
% ./pippin.py wallet_change_seed --wallet d897b5ec-1897-4e7e-8a90-4526f454c8de
Enter new wallet seed: <hidden_input>
Seed changed for wallet d897b5ec-1897-4e7e-8a90-4526f454c8de
First account: nano_3ejy6ha1iuqhi5cshhifu57p5othdcymfbzsmxhjucdks53eh41yd4qpjtxf
```

To backup a seed (**warning:** this prints seed to stdout)

```
% ./pippin wallet_view_seed --wallet <id>
```

## Setting up Pippin

### Requirements

- Python 3.7 or newer
- GCC, for MacOS and Linux
- libb2 (blake2b)
- A Redis server

On MacOS, with homebrew:

```
% brew install gcc@9 python libb2 redis
% launchctl load ~/Library/LaunchAgents/homebrew.mxcl.redis.plist
```

To start redis at boot on MacOS:

```
% ln -sfv /usr/local/opt/redis/*.plist ~/Library/LaunchAgents
```

MacOS users may find it convenient to priorize homebrew binaries.

```
export PATH=/usr/local/bin:$PATH
```

This means you'll be using the homebrew installed python3 by default, to make it permanent:

```
# Catalina
echo "export PATH=/usr/local/bin:$PATH" >> ~/.zprofile
# Others
echo "export PATH=/usr/local/bin:$PATH" >> ~/.profile
```

On Linux, debian-based systems:

```
sudo apt install build-essential python3.7 python3.7-dev libb2-dev redis-server
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

### Configuring Pippin for BANANO

Pippin uses an environmnet variable to switch between banano and nano mode.

To use with banano:

```
% echo "BANANO=1" >> .env
```

### Configuring PostgreSQL or MySQL

Pippin uses SQLite by default, which requires no extra configuration.

To use postgres or mysql, you need to put your database information in some environment variables

**Postgres:**

Required (replace `database_name`, `user_name`, and `mypassword` with the actual values):
```
echo "POSTGRES_DB=database_name" >> .env
echo "POSTGRES_USER=user_name" >> .env
echo "POSTGRES_PASSWORD=mypassword" >> .env
```

Optional:
```
# 127.0.0.1 is default
echo "POSTGRES_HOST=127.0.0.1" >> .env 
# 5432 is default
echo "POSTGRES_PORT=5432" >> .env
```

**MySQL:**

Required (replace `database_name`, `user_name`, and `mypassword` with the actual values):
```
echo "MYSQL_DB=database_name" >> .env
echo "MYSQL_USER=user_name" >> .env
echo "MYSQL_PASSWORD=mypassword" >> .env
```

Optional:
```
# 127.0.0.1 is default
echo "MYSQL_HOST=127.0.0.1" >> .env 
# 3306 is default
echo "MYSQL_PORT=3306" >> .env
```

### Configuring the node

Create a file called `config.yaml` in the root directory of pippin.

Pippin uses a yaml based configuration for everything else. There's a sample provided at `sample.config.yaml` where you can reference all available options.

At the bare minimum, Pippin requires a node for the RPC api.

```
server:
  # The RPC URL of the remote node to connect to
  # Non-wallet RPCs will be routed to this node
  # Default: http://[::1]:7076 for nano, https://[::1]:7072 for banano
  node_rpc_url: https://coolnanonode.com/rpc
```

### Running Pippin

Once configured, just start it with `python3 main.py`

It can be started on boot using systemd (Linux)

Create a file `/etc/systemd/system/pippin.service`

With the contents:

```
[Unit]
Description=Pippin Wallet
After=network.target

[Service]
Type=simple
User=YOUR_LINUX_USER
Group=YOUR_LINUX_USER
WorkingDirectory=/home/YOUR_LINUX_USER/pippin_nano_wallet
EnvironmentFile=/home/YOUR_LINUX_USER/pippin_nano_wallet/.env
ExecStart=/home/YOUR_LINUX_USER/pippin_nano_wallet/venv/bin/python main.py

[Install]
WantedBy=multi-user.target
```

Then enable and start

```
% sudo systemctl enable pippin
% sudo systemctl start pippin
```

### Endpoints

Send HTTP POST requests to Pippin just like you would a normal node.

```
% curl -g -d '{"action":"wallet_create"}' localhost:11338
% curl -g -d '{"action":"account_balance", "account": "nano_3jb1fp4diu79wggp7e171jdpxp95auji4moste6gmc55pptwerfjqu48okse"}' localhost:11338
```

## Feature requests

Notice an API that's missing a feature or not behaving the same as nano's APIs?

Open a bug report/feature request on the [issues page](https://github.com/bbedward/pippin_nano_wallet/issues)