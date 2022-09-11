package controller

import (
	rpc "github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
)

type HttpController struct {
	Wallet    *wallet.NanoWallet
	RpcClient *rpc.RPCClient
}
