<p align="center">
  <img src="https://raw.githubusercontent.com/appditto/pippin_nano_wallet/master/assets/pippin_header.png?sanitize=true" alt="Pippin Wallet" width="500">
</p>

# CLI

This is Pippin's CLI tool. It (currently) doesn't offer full wallet support like the pippin-server rest APIs, but it provides CLI tools for managing keys and wallets.

It comes bundled in the pippin distributions as `pippin-cli(.exe)`

For usage, see: `pippin-cli --help`

Some examples are:

```bash
# List all wallets and accounts
% pippin-cli wallet --list
# Create a wallet
% pippin-cli wallet --create
# Create a wallet with a specific seed
% pippin-cli wallet --create --seed daaf0390c20e7f646759d1f3b93e55a727147bb5649f7e4945dd0afabd29fe12
```