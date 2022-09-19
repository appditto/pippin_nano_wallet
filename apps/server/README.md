<p align="center">
  <img src="https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/assets/pippin_header.png?sanitize=true" alt="Pippin Wallet" width="500">
</p>

# Server

This is Pippin's Server. It provides a rest API that provides full wallet functionality for the [Nano](https://nano.org) and [BANANO](https://banano.cc) cryptocurrencies.

It is engineered to replicate the [Nano Developer Wallet APIs](https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs) and to be a drop in replacement.

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

## API Differences - Nano vs Pippin

These are the known differences between Pippin's API and the Nano node wallet API. There may be more that are not listed here, it is up to you to ensure your application properly integrates with Pippin.

If you find an API that behaves differently, [create an issue](https://github.com/appditto/pippin_nano_wallet/issues) with the details.

**Different Behavior**

APIs that are different between Pippin and the Nano node wallet.

- `account_list` accepts a `count` parameter that defaults to 1000
- Pippin has an `auto_receive_on_send` configuration option that will automatically receive pending blocks when you do a `send`, it will only do this if the source balance isn't high enough to make the transaction.

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