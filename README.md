<p align="left">  <img src="https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/assets/pippin_logo.svg?sanitize=true" alt="Pippin Wallet" width="256" height="50"> </p>

[![PyPI](https://img.shields.io/pypi/v/pippin-wallet)](https://pypi.org/project/pippin-wallet/) ![PyPI - Python Version](https://img.shields.io/pypi/pyversions/pippin-wallet) [![License](https://img.shields.io/github/license/bbedward/pippin_nano_wallet)](https://github.com/bbedward/pippin_nano_wallet/blob/master/LICENSE) [![CI](https://github.com/appditto/pippin_nano_wallet/workflows/CI/badge.svg)](https://github.com/appditto/pippin_nano_wallet/actions?query=workflow%3ACI)

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
- Pippin natively supports [DPoW](https://dpow.nanocenter.org) and [BPoW](https://bpow.banano.cc)
- Pippin supports multiple database backends (SQLite, PostgreSQL, and MySQL)

Pippin can be used by exchanges, games, payment processors, tip bots, faucets, casinos, and a lot more.

## Pippin Performance

<p align="center">  <img src="https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/assets/benchmark.png" alt="Pippin Benchmarks"> </p>

The benchmark script that was used is available [here](https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/benchmark.py)

There were 3 independent runs for each wallet, 62 blocks each, same node, and the same work peer. Pippin was consistently **twice as fast** versus the node wallet.

## How Pippin Works

Pippin provides an API that mimics the [Nano Wallet RPC Protocol](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs)

Every wallet-related RPC gets intercepted by Pippin and handled internally. It builds the blocks and signs them using locally-stored keys, it uses a node to publish the blocks.

Every non-wallet related RPC gets proxied to the publishing node. Which means you can make all of your RPC requests directly to Pippin whether they are wallet-related or not.

## API Documentation

Recommended reference is the [NANO RPC documentation](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs), Pippin's APIs are mostly identical.

You send an HTTP Post request to pippin with the desired action and parameters, example:

```
{
    "action": "accounts_create",
    "wallet": "12345",
    "count": 100
}
```

### Supported

- `wallet_create`
- `account_create`
- `accounts_create`
- `account_list`
- `receive`
- `send` - Use the **id** parameter to prevent duplicate sends!
- `account_representative_set`
- `password_change` - This will also set a password, if one isn't already set
- `password_enter`
- `password_valid`
- `wallet_representative_set`
- `wallet_add` - This is for adding ad-hoc private keys to a wallet
- `wallet_lock`
- `wallet_locked`
- `wallet_balances`
- `wallet_frontiers`
- `wallet_pending`
- `wallet_destroy`
- `wallet_change_seed`
- `wallet_contains`
- `wallet_representative`
- `receive_all` - Not in the nano API, it takes a `wallet` and it will receive every pending block in that wallet (respecting `receive_minimum`).

### Differences: Pippin vs NANO Node Wallet

These are the known differences between Pippin's API and the Nano node wallet API

**Different Behavior**

APIs that are different between Pippin and the Nano node wallet.

- `account_list` accepts a `count` parameter that defaults to 1000
- Pippin has an `auto_receive_on_send` option that will automatically receive pending blocks when you do a `send`, it will only do this if balance isnt high enough to make the transaction.
- `account_create` does not accept an index

**Fuzzy Behavior**

The Nano documentation isn't perfectly clear on these, but these are how Pippin behaves.

- `wallet_change_seed` will result in the wallet no longer being locked, if it is. The wallet has to already be unlocked before you can use this RPC, though.

**Missing/Not Implemented**

APIs that the Nano node wallet supports but are not implemented in Pippin.

- `account_move`
- `account_remove`
- `receive_minimum` - Receive minimum can be set in `config.yaml`
- `receive_minimum_set`
- `wallet_add_watch`
- `wallet_history`
- `search_pending`
- `search_pending_all`
- `wallet_export`
- `wallet_ledger`
- `wallet_republish`
- `wallet_work_get`
- `work_get`
- `work_set`

## CLI Documentation

Pippin has a CLI interface available, you can see available subcommands with:

`pippin-cli --help`

The primary goal of the CLI is key management. It's a more secure way to import a seed and backup your seed.

For example a typical flow of creating a new wallet with a specific seed might look like (add --encrypt to wallet_change_seed if you want to lock the wallet with a password):

```
% pippin-cli wallet_create
Wallet created, ID: d897b5ec-1897-4e7e-8a90-4526f454c8de
First account: nano_31a7wzm4rayik1hthahzkekntsqz86u6dko5adg8jxueehzt5yhmhsqsuzdy
% pippin-cli wallet_change_seed --wallet d897b5ec-1897-4e7e-8a90-4526f454c8de
Enter new wallet seed: <hidden_input>
Seed changed for wallet d897b5ec-1897-4e7e-8a90-4526f454c8de
First account: nano_3ejy6ha1iuqhi5cshhifu57p5othdcymfbzsmxhjucdks53eh41yd4qpjtxf
```

To backup a seed (**warning:** this prints seed to stdout)

```
% pippin-cli wallet_view_seed --wallet <id>
```

## Setting up Pippin

### Requirements

- Python 3.6 or newer
- GCC, for MacOS and Linux
- A Redis server

**MacOS Instructions**

1) Install [homebrew](https://brew.sh/), if it isn't already installed

2) Install Redis (skip if you already have done so)

```
% brew install gcc@9 python redis
% launchctl load ~/Library/LaunchAgents/homebrew.mxcl.redis.plist
```

To start redis at boot (optional):

```
% ln -sfv /usr/local/opt/redis/*.plist ~/Library/LaunchAgents
```

You may find it convenient to priorize homebrew binaries by including the install location first in your PATH.

```
% export PATH=/usr/local/bin:$PATH
```

To make it permanent:

```
# Catalina
% echo "export PATH=/usr/local/bin:$PATH" >> ~/.zprofile
# Others
% echo "export PATH=/usr/local/bin:$PATH" >> ~/.profile
```

**Ubuntu 18.04**

Instructions for other debian-based Linux distributions should be similar.

```
% sudo apt install build-essential python3.6 python3.6-dev python3-pip redis-server
```

**CentOS 8**

Install the required developer tools:
```
# dnf install gcc redis python3-devel
```

The above may require the EPEL and PowerTools repos to be configured first:
```
# dnf install -y install https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm
# dnf config-manager --set-enabled PowerTools
```

**Redis basics for CentOS 8**

After installing Redis, create directories for the Redis config and runtime files:
```
# mkdir /etc/redis /var/redis /var/redis/6379
```
By default, the dnf install places Redis config files in /etc, move them to your dedicated folders:
```
# mv /etc/redis-sentinel.conf /etc/redis/ 
# mv /etc/redis.conf /etc/redis/6379.conf
```
Update your Redis config file to allow it to run in the background as a daemon, supervised by systemd:
```
# vim /etc/redis/6379.conf
	Set daemonize to yes (by default it is set to no).
	Set the pidfile to /var/run/redis_6379.pid (modify the port if needed).
	Change the port if necessary (6379 is the default).	
	Set the logfile to /var/log/redis_6379.log
	Set the working dir to /var/redis/6379 
	supervised systemd
```

Create the Systemd service to have Redis start automatically:
```
# vim /etc/systemd/system/redis.service
[Unit]
Description=Redis
After=syslog.target

[Service]
Type=notify
PIDFile=/var/run/redis_6379.pid
ExecStart=/usr/bin/redis-server /etc/redis/6379.conf --supervised systemd
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Enable and start the service:
```
# systemctl enable --now /etc/systemd/system/redis.service
```

Verify that Redis is running:
```
# systemctl status redis
# redis-cli ping
```

### Installing Pippin

First, update PIP to latest version.

```
% pip3 install -U pip
```

For MacOS you might need to set the following environment variable:

```
export CC=/usr/local/bin/gcc-9
```

To install Pippin on macos or Linux

```
% pip install --user pippin-wallet
```

Windows requires visual c++ and should be prefixed with USE_VC=1

```
% USE_VC=1 pip install pippin-wallet
```

To upgrade Pippin in the future, add --upgrade

```
% pip install --upgrade pippin-wallet
```

### Configuring Pippin

Pippin creates a `PippinData` directory in your home directory.

Run: `pippin-server --generate-config` to generate a sample in `~/PippinData/sample.config.yaml`

### Using Distributed PoW or BoomPoW

Want to use [DPoW](https://dpow.nanocenter.org) or [BPoW](https://bpow.banano.cc)?

Pippin will use them automatically for work generation if the key/user is present in the environment.

For DPoW:
```
% echo "DPOW_USER=mydpowuser" >> ~/PippinData/.env
% echo "DPOW_KEY=mydpowkey" >> ~/PippinData/.env
```

For BPoW:
```
% echo "BPOW_USER=mybpowuser" >> ~/PippinData/.env
% echo "BPOW_KEY=mybpowkey" >> ~/PippinData/.env
```

Replace `mybpowuser` and `mybpowkey` with the actual user and keys you have. If you need keys, visit their respected websites for instructions on how to request them.

### Configuring PostgreSQL or MySQL

Pippin uses SQLite by default, which requires no extra configuration.

To use postgres or mysql, you need to put your database information in some environment variables

**Postgres:**

Required (replace `database_name`, `user_name`, and `mypassword` with the actual values):
```
% echo "POSTGRES_DB=database_name" >> ~/PippinData/.env
% echo "POSTGRES_USER=user_name" >> ~/PippinData/.env
% echo "POSTGRES_PASSWORD=mypassword" >> ~/PippinData/.env
```

Optional:
```
# 127.0.0.1 is default
% echo "POSTGRES_HOST=127.0.0.1" >> ~/PippinData/.env 
# 5432 is default
% echo "POSTGRES_PORT=5432" >> ~/PippinData/.env
```

**MySQL:**

Required (replace `database_name`, `user_name`, and `mypassword` with the actual values):
```
% echo "MYSQL_DB=database_name" >> ~/PippinData/.env
% echo "MYSQL_USER=user_name" >> ~/PippinData/.env
% echo "MYSQL_PASSWORD=mypassword" >> ~/PippinData/.env
```

Optional:
```
# 127.0.0.1 is default
% echo "MYSQL_HOST=127.0.0.1" >> ~/PippinData/.env 
# 3306 is default
% echo "MYSQL_PORT=3306" >> ~/PippinData/.env
```

### Changing Redis Host/Port

Pippin uses Redis for distributed locks, so that every account works on its own chain in a synchronous fashion.

By default, it will look for redis on `127.0.0.1` on port `6379` and use db `0`, you can also change these with environment variables.

```
echo "REDIS_HOST=127.0.0.1" >> ~/PippinData/.env
echo "REDIS_PORT=6379" >> ~/PippinData/.env
echo "REDIS_DB=0" >> ~/PippinData/.env
echo "REDIS_PW=supersecretpassword" >> ~/PippinData/.env
```

## Pippin Configuration

Pippin uses a [yaml](https://yaml.org/) based configuration for everything else.

All available options are in a [sample file](./pippin/sample.config.yaml).

You can override any default by creating a file called `~/PippinData/config.yaml` and choosing your own settings.

It must be in your users home directory in a folder called `PippinData`

### Configuring Pippin for BANANO

In `config.yaml` set banano: true

```
# Settings for the pippin wallet
wallet:
  # Run in banano mode
  # If true, the wallet will operate based on the BANANO protocol
  # Default: false
  banano: true
```

### Configuring the node

At the bare minimum, Pippin requires a node for the RPC api. It will default to `http://[::1]:7076` for Nano, or `http://[::1]:7072` for BANANO. If you want to change it to `https://coolnanonode.com/rpc` then it would look like this:

```
server:
  # The RPC URL of the remote node to connect to
  # Non-wallet RPCs will be routed to this node
  # Default: http://[::1]:7076 for nano, https://[::1]:7072 for banano
  node_rpc_url: https://coolnanonode.com/rpc
```

### Running Pippin

Once configured, just start it with `pippin-server`

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
ExecStart=/home/myuser/.local/bin/pippin-server

[Install]
WantedBy=multi-user.target
```

If you aren't sure what the full path of pippin-server is, run `which pippin-server`

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

### Auto-receive & Dynamic PoW

To automatically pocket pending transactions as they arrive, callback is required.

Hooking up the websocket also adds support for **dynamic PoW** which means that blocks will get confirmed faster if the active difficulty is higher than the minimum.

Pippin only supports the websocket callback, which can be setup like so in `config.yaml`:

```
server:
  # The WebSocket URL of the node to connect to
  # Optional, but required to receive transactions as they arrive to accounts
  # Default: None
  #node_ws_url: ws://[::1]:7078
```

## Feature requests

Notice an API that's missing a feature or not behaving the same as nano's APIs?

Open a bug report/feature request on the [issues page](https://github.com/bbedward/pippin_nano_wallet/issues)
