# Changelog

## [1.1.19] - 2021-03-14

- Support work_generate action, supports 'subtype' to automatically choose correct difficulty
- Subscribe to active_difficulty websocket for dynamic PoW

## [1.1.18] - 2021-02-23

- Fix server hanging when DPoW/BPoW is unavailable
- Fix auto receives working for unopened accounts

## [1.1.17] - 2020-08-20

- Fix transactions made to accounts with xrb_ prefix
- Fix pyyaml dependency

## [1.1.16] - 2020-08-17

- Fix backwards compatibility

## [1.1.15] - 2020-08-17

- Support Nano v21 PoW Difficulty

## [1.1.13] - 2020-08-09

- Require aioredlock 0.3.0

## [1.1.12] - 2020-05-19

- Fix pending RPC when receive_minimum is < 1

## [1.1.10] - 2020-05-05

- Remove non-voting default rep

## [1.1.9] - 2020-05-04

- Windows compatibility

## [1.1.7] - 2020-04-27

- Reduce log spam

## [1.1.6] - 2020-04-27

- Fix difficulty adjustment

## [1.1.5] - 2020-04-26

- Keep sane difficulty levels in DPoW and BPoW requests

## [1.1.4] - 2020-04-22

- Bump tortoise-orm minimum version to 0.15.24

## [1.1.3] - 2020-04-20

- Don't generate sample config if config.yaml exists already

## [1.1.2] - 2020-04-20

- Add log_to_stdout option
- Changes for docker environments

## [1.1.1] - 2020-04-18

- Fix ipv6 rpc/work URLs

## [1.1.0] - 2020-04-15

- Make requirements less strict
- Fix aiohttp TCPConnector in some cases

## [1.0.8] - 2019-12-16

- Fix update_existing with wallet_representative_set when accounts aren't open
- Add CLI commands: wallet_destroy, wallet_representative_set, wallet_representative_get
- Add setuptools dependency

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
