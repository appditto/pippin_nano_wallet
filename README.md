# Pippin NANO and BANANO Wallet

Pippin is a highly performing API-based wallet for NANO and BANANO.

It's API is generally designed to be a drop-in replacement for the standard node wallet with some exceptions that are documented below:

- `account_list` accepts a `count` parameter to limit results that defaults to 1000
- `account_create` does not accept an index
- `account_move` is missing
- `account_remove` is missing

The missing items are missing due to not having been implemented yet. Generally, the most commonly uses wallet RPCs are implemented in Pippin including `send`, `receive`, `wallet_representative_set`, `account_representative_set`, etc.