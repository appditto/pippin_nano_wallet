# Pippin API/CLI Wallet for NANO and BANANO

Pippin is a high performance API-based wallet for NANO and BANANO.

Pippin uses the fast, C-based [nanopy](https://github.com/npy0/nanopy/tree/master/nanopy) libary to sign blocks and generate work. It is also completely written using asyncio and utilizes the ultra-fast [uvloop](https://github.com/MagicStack/uvloop)

## Why?

Pippin has several use case, here's an excerpt from the NANO documentation

```
Below are RPC commands that interact with the built-in, QT-based node wallet. This wallet is only recommended for development and testing. For production integrations, setting up custom External Management processes is required.
```

NANO Foundation advises against using the built-in node wallet for production environments, Pippin is a production-ready wallet that provides many of the same APIs developers and apps are already using.

It may be useful for exchanges, nano-based games, tipbots, payment processors, and any other application that needs to store nano keys and create nano blocks.

Also:

- Pippin can be used with *any public node* - you don't need to run your own node to have an API wallet.
- Pippin is extremely fast and lightweight
- Pippin supports multiple database backends (SQLite, PostgreSQL)

## How?

Pippin intercepts every wallet-related RPC and handles them internally. It builds its own blocks, has its own storage, etc. Every non-wallet RPC gets proxied to a remote or local node.

Once pippin is configured, all you need to do is point your existing NANO/BANANO application to it.

## API Documentation

Pippin is designed to be a drop-in replacement for the standard node wallet, with some exceptions (see the next section)

You should reference the [NANO RPC documentation](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs) to view a list of all of the available APIs pippin supports.

## Known Differences: Pippin vs NANO Node Wallet

**Enhanced Behavior**

- `account_list` accepts a `count` parameter that defaults to 1000

**Degraded Behavior**

- `account_create` does not accept an index

**Fuzzy**

- `wallet_change_seed` decrypts the wallet, if it's encrypted

**Missing/Not Implemented**

- `account_move`
- `account_remove`
- `receive_minimum`
- `receive_minimum_set`
- `wallet_add_watch`

## CLI Documentation

Unlike the API, the CLI is entirely different from the nano node CLI.

You can see all available options with the CLI using `./pippin --help`

### Something missing?

Notice an API that's missing a feature or not behaving the same as nano's APIs?

Open a bug report/feature request on the [issues page](https://github.com/bbedward/pippin_nano_wallet/issues)