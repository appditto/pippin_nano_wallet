# Work In Progress

This is the active development rewrite of pippin, for the current version see the [legacy branch](https://github.com/appditto/pippin_nano_wallet/tree/legacy)

[![Release](https://img.shields.io/github/v/release/appditto/pippin_nano_wallet)](https://github.com/appditto/pippin_nano_wallet/releases/latest) ![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/appditto/pippin_nano_wallet?filename=apps%2Fserver%2Fgo.mod) [![License](https://img.shields.io/github/license/appditto/pippin_nano_wallet)](https://github.com/appditto/pippin_nano_wallet/blob/master/LICENSE) [![CI](https://github.com/appditto/pippin_nano_wallet/workflows/CI/badge.svg)](https://github.com/appditto/pippin_nano_wallet/actions?query=workflow%3ACI)

<p align="center">
  <img src="https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/assets/pippin_header.png?sanitize=true" alt="Pippin Wallet" width="500">
</p>

## About Pippin

Pippin is a production-ready, high-performance API and CLI wallet for [Nano](https://nano.org) and [BANANO](https://banano.cc). Pippin's API is a drop-in replacement for the Nano developer wallet that is built in to the Nano node software.

It is recommended for exchanges, wallets, faucets, and any other application that integrates with the Nano or BANANO networks.

## Benefits of Pippin

The Nano developer wallet (aka "node wallet") is not recommended for production use. One of the goals of Pippin is to provide a production-ready external key management that can be used to integrate with Nano or BANANO.

Pippin is the first drop-in replacement for the Nano developer wallet. It's incredibly easy to transition to Pippin if you are already using the Nano developer wallet.

- Pippin is independent of the node. You can use Pippin with any public RPC, so you don't have to run your own node
- Pippin is extremely fast and lightweight
- Pippin supports encrypted secret keys
- Pippin natively supports [BoomPoW](https://boompow.banano.cc)
- Pippin supports multiple database backends (SQLite, PostgreSQL, and MySQL)
- Pippin is scalable and synchronizes state between all instances using [Redis](https://redis.io)

## Pippin Performance

coming soon

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
    "wallet": "186e3283-f27d-4ef5-87e3-84322dd740a2",
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
- `password_change` - This is how you set a password, if one isn't already set
- `password_enter`
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

### Removed in v3.0

- `password_valid` - This one has little to no purpose, `password_enter` indicates if your password is valid or not.

### Wallet Lock

You can optionally encrypt the seed+private keys associated with a wallet, by default seeds are not encrypted in the database backend. (this is ok, if your database is secure).

Details of these APIs are in the [NANO RPC documentation](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs)

The flow would be to:

1) Use `password_change` to set a wallet password, this will also lock the wallet APIs
2) Use `password_enter` to unlock the wallet, using the password
3) When your session is over, use `wallet_lock` to re-lock the wallet.

**If you want to remove the password from the wallet, use `password_change` with an empty password, while the wallet is unlocked**

When locked, any RPCs that interact with the wallet will return an error code, these include:

- `account_create`
- `accounts_create`
- `account_list`
- `receive`
- `send`
- `account_representative_set`
- `password_change`
- `wallet_representative_set`
- `wallet_add`
- `wallet_balances`
- `wallet_frontiers`
- `wallet_pending`
- `wallet_destroy` - You can use the CLI to destroy a wallet if you forget the password
- `wallet_change_seed`
- `wallet_contains`
- `wallet_representative`
- `receive_all`

### Differences: Pippin vs NANO Node Wallet

These are the known differences between Pippin's API and the Nano node wallet API. There may be more that are not listed here, it is up to you to ensure your application properly integrates with Pippin.

**Different Behavior**

APIs that are different between Pippin and the Nano node wallet.

- `account_list` accepts a `count` parameter that defaults to 1000
- Pippin has an `auto_receive_on_send` option that will automatically receive pending blocks when you do a `send`, it will only do this if the source balance isn't high enough to make the transaction.
- `account_create` does not accept an index. If you want to add a specific account to a wallet, you can use `wallet_add`

**Fuzzy Behavior**

The Nano documentation isn't perfectly clear on these, but these are how Pippin behaves.

- `wallet_change_seed` will result in the wallet no longer being locked/encrypted.

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

The primary goal of the CLI is key management. It does not (currently) provide all wallet actions.

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

**Coming Soon**

### Configuring Pippin

Pippin creates a `PippinData` directory in your home directory.

Run: `pippin-server --generate-config` to generate a sample in `~/PippinData/sample.config.yaml`

### Using BoomPoW

Want to use [BoomPoW](https://boompow.banano.cc)?

Pippin will use them automatically for work generation if the key is present in the environment.

For BPoW:

```
% echo "BPOW_KEY=service:mybpowkey" >> ~/PippinData/.env
```

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

Pippin uses Redis for distributed locks and generally synchronizing state, so that every account works on its own chain in a synchronous fashion.

By default, it will look for redis on `127.0.0.1` on port `6379` and use db `0`, you can also change these with environment variables.

```
echo "REDIS_HOST=127.0.0.1" >> ~/PippinData/.env
echo "REDIS_PORT=6379" >> ~/PippinData/.env
echo "REDIS_DB=0" >> ~/PippinData/.env
```

## Pippin Configuration

Pippin uses a [yaml](https://yaml.org/) based configuration for everything else.

All available options are in a [sample file](./sample.config.yaml).

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

**coming soon**

### Endpoints

Send HTTP POST requests to Pippin just like you would a normal node.

```
% curl -g -d '{"action":"wallet_create"}' localhost:11338
% curl -g -d '{"action":"account_balance", "account": "nano_3jb1fp4diu79wggp7e171jdpxp95auji4moste6gmc55pptwerfjqu48okse"}' localhost:11338
```

### Auto-receive & Dynamic PoW

To automatically pocket pending transactions as they arrive, callback is required.

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
