# Work In Progress

This is the active development rewrite of pippin, for the current version see the [legacy branch](https://github.com/appditto/pippin_nano_wallet/tree/legacy)

[![Release](https://img.shields.io/github/v/release/appditto/pippin_nano_wallet)](https://github.com/appditto/pippin_nano_wallet/releases/latest) ![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/appditto/pippin_nano_wallet?filename=apps%2Fserver%2Fgo.mod) [![License](https://img.shields.io/github/license/appditto/pippin_nano_wallet)](https://github.com/appditto/pippin_nano_wallet/blob/master/LICENSE) [![CI](https://github.com/appditto/pippin_nano_wallet/workflows/CI/badge.svg)](https://github.com/appditto/pippin_nano_wallet/actions?query=workflow%3ACI)

<p align="center">
  <img src="https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/assets/pippin_header.png?sanitize=true" alt="Pippin Wallet" width="500">
</p>

## About Pippin

Pippin is a production-ready, high-performance **API/developer** wallet for [Nano](https://nano.org) and [BANANO](https://banano.cc). Pippin's API is a drop-in replacement for the Nano developer wallet that is built in to the Nano node software.

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


## Components

Pippin has a *server* and a *cli* interface. See the appropriate README for more details on each.

- [Server](https://github.com/appditto/pippin_nano_wallet/tree/master/apps/server)
- [CLI](https://github.com/appditto/pippin_nano_wallet/tree/master/apps/cli)

The **CLI** is the main entrypoint for the entire app, including the server.

## Setting up Pippin

Pippin comes as a pre-compiled binary for multiple architectures, which can be downloaded on the [releases](https://github.com/appditto/pippin_nano_wallet/releases) page.

You can start the server with `pippin --start-server`, or use `pippin --help` for a full list of available CLI options.

## Server Configuration

Pippin is configured through a `yaml` style configuration for most things, but some things are configured via `env` variables.

The first step is to determine if you want to change the location where pippin stores its data, by default it will be in the users home directory (`~`)

Which translates to `/home/<username>` on Linux, `/Users/<username>` on MacOS and `C:\Users\<username>`

You can override this with the `PIPPIN_HOME` environment variable.

### Configuring Pippin

Pippin creates a `PippinData` directory in your home directory.

Run: `pippin --generate-config` to generate a sample in `~/PippinData/sample.config.yaml`

After editing your parameters in this file, **move it to ~/PippinData/config.yaml**

### Configuring Database

By default, Pippin will use a SQLite database that is created in `$PIPPIN_HOME/PippinData/pippingo.db`

Pippin also supports `MySQL` and `PostgreSQL` which is configured in the environment.

You can set these variables the same way that you normally set environment variables, but for convenience pippin will read `$PIPPIN_HOME/PippinData/.env`

For MySQL
```bash
% echo "MYSQL_DB=database_name" >> ~/PippinData/.env
% echo "MYSQL_USER=user_name" >> ~/PippinData/.env
% echo "MYSQL_PASSWORD=mypassword" >> ~/PippinData/.env
```

For Postgres
```bash
% echo "POSTGRES_DB=database_name" >> ~/PippinData/.env
% echo "POSTGRES_USER=user_name" >> ~/PippinData/.env
% echo "POSTGRES_PASSWORD=mypassword" >> ~/PippinData/.env
```

### Configuring Redis

[Redis](https://redis.io) is a non-optional requirement for Pippin. It allows Pippin to be scalable across multiple instances and handles distributed locking.

By default it will use `localhost:6379` to connect and `0` as the database. These can be overriden with the environment variables:

- `REDIS_HOST`
- `REDIS_PORT`
- `REDIS_DB`

### Using BoomPoW

Want to use [BoomPoW](https://boompow.banano.cc)?

Pippin will use them automatically for work generation if the key is present in the environment.

For BPoW:

```
% echo "BPOW_KEY=service:mybpowkey" >> ~/PippinData/.env
```

### Using GPU/OpenCL To Generate PoW Locally

The pre-compiled pippin distributions do not support GPU PoW out of the box, however Pippin can be compiled that way to enable it with something like:

`go build -tags cl -o pippin ./apps/cli`

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

The `node_ws_url` corresponds to the URL to use for the [Node Websocket API](https://docs.nano.org/integration-guides/websockets/)

It is **optional** but should take the form of `ws://[::1]:7078`

The websocket is only used to automatically receive transactions for unlocked wallets.

### Running Pippin

After configuration is complete, simply run `pippin --start-server`

### Endpoints

Send HTTP POST requests to Pippin just like you would a normal node.

```
% curl -g -d '{"action":"wallet_create"}' localhost:11338
% curl -g -d '{"action":"account_balance", "account": "nano_3jb1fp4diu79wggp7e171jdpxp95auji4moste6gmc55pptwerfjqu48okse"}' localhost:11338
```

## Feature requests

Notice an API that's missing a feature or not behaving the same as nano's APIs?

Open a bug report/feature request on the [issues page](https://github.com/bbedward/pippin_nano_wallet/issues)
