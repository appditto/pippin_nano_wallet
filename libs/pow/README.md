# PoW (Proof of Work)

This is the core module for requesting [PoW](https://docs.nano.org/integration-guides/work-generation/), which is required to publish blocks to the Nano or BANANO networks.

There are 3 mechanisms for requesting work within this module.

1) [BoomPoW](https://boompow.banano.cc)
2) API-driven, using something like [Nano work server](https://github.com/nanocurrency/nano-work-server) or a node itself.
3) Local PoW, using [nanopow](https://github.com/inkeliz/nanopow) - if compiled with `-tags cl` it will utilize OpenCL CGO bindings to calculate PoW on GPU.

These mechanisms can all be invoked separately, or by Pippin's mechanism (`WorkGenerateMeta` function).

`WorkGenerateMeta` will

1) Execute concurrent goroutines requesting from BoomPoW and all configured work servers.
2) When first result comes back, cancel all pending goroutines and send work_cancel to all work servers.
3) If API fails, we generate PoW locally and set a flag `WorkFailing`, then subsequent requests will use local PoW along with the peers until the peers are working again

APIs are preferred, if no APIs are configured then local work generation  will be the primary mechanism.