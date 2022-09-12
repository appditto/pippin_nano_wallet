package controller

import (
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	rpc "github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
)

type HttpController struct {
	Wallet    *wallet.NanoWallet
	RpcClient *rpc.RPCClient
	PowClient *pow.PippinPow
}
