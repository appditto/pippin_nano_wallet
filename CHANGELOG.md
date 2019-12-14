# Changelog

## [1.0.7] - 2019-12-14

- (Fix) Compatibility with PyYAML < 5.1

## [1.0.6] - 2019-12-14

- Compatibility with PyYAML < 5.1

## [1.0.5] - 2019-12-14

- Add support for dynamic PoW (requires websocket configured)
- Fix issues related to PyYAML dependency

## [1.0.4] - 2019-12-13

- Minor bug fixes

## [1.0.3] - 2019-12-13

- Add support for Python 3.6

## [1.0.2] - 2019-12-13

- Restructure project for pypi compatibility

## [1.0.1] - 2019-12-12

- Fix an issue with multiple work peers
- Add [Distributed PoW](https://dpow.nanocenter.org) and [BoomPoW support](https://bpow.banano.cc)
- Add account_create to CLI
- Fix issues with closing websockets when Pippin stopped
- Use `subtype` when sending RPC process to the node
- Fix so RPC send/receive/account_rep_set return the exact same response as node wallet

## [1.0.0] - 2019-12-10

- Initial Release