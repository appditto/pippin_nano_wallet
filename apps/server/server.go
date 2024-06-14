package server

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/appditto/pippin_nano_wallet/apps/server/controller"
	"github.com/appditto/pippin_nano_wallet/apps/server/middleware"
	"github.com/appditto/pippin_nano_wallet/apps/server/net"
	"github.com/appditto/pippin_nano_wallet/libs/config"
	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	rpc "github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/chi/v5"
)

func StartPippinServer() {
	// Read yaml configuration
	conf, err := config.ParsePippinConfig()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	// Setup database conn
	ctx := context.Background()
	fmt.Println("🏡 Connecting to database...")
	dbconn, err := database.GetSqlDbConn(false)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	entClient, err := database.NewEntClient(dbconn)
	if err != nil {
		log.Fatalf("Failed to create ent client: %v", err)
		os.Exit(1)
	}
	defer entClient.Close()

	// Run migrations
	log.Info("🦋 Running migrations...")
	if err := entClient.Schema.Create(ctx); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
		os.Exit(1)
	}

	// Setup RPC handlers
	rpcClient := rpc.NewRPCClient(conf.Server.NodeRpcUrl)

	// Setup pow client
	pow := pow.NewPippinPow(conf.Wallet.WorkPeers, utils.GetEnv("BPOW_KEY", ""), utils.GetEnv("BPOW_URL", ""))

	// Setup nano wallet instance with DB, options, etc.
	nanoWallet := wallet.NanoWallet{
		DB:         entClient,
		Ctx:        ctx,
		Banano:     conf.Wallet.Banano,
		RpcClient:  rpcClient,
		WorkClient: pow,
		Config:     conf,
	}

	// Setup nano WS client if configured
	callbackChan := make(chan *net.WSCallbackMsg, 100)
	if conf.Server.NodeWsUrl != "" {
		go net.StartNanoWSClient(conf.Server.NodeWsUrl, &callbackChan, &nanoWallet)
	}

	// Read channel to automatically receive blocks
	go func() {
		for msg := range callbackChan {
			func() {
				// Lock each callback so we don't handle them on multiple instances
				lock, err := database.GetRedisDB().Locker.Obtain(ctx, fmt.Sprintf("blocklock:%s", msg.Hash), time.Second*30, nil)
				if err != nil {
					return
				}
				defer lock.Release(ctx)
				// Ignore non-sends
				if msg.Block.Subtype != "send" || msg.IsSend != "true" {
					return
				}
				amount, ok := big.NewInt(0).SetString(msg.Amount, 10)
				if !ok {
					return
				}
				// Compare to receive minimum
				receiveMinimum, ok := big.NewInt(0).SetString(conf.Wallet.ReceiveMinimum, 10)
				if !ok {
					return
				}
				if amount.Cmp(receiveMinimum) < 0 {
					return
				}

				// See if destination is in our wallet
				dbAccount, err := nanoWallet.GetAccountByAddress(msg.Block.LinkAsAccount)
				if err != nil {
					return
				}
				wallet, err := dbAccount.QueryWallet().Only(ctx)
				if err != nil {
					return
				}

				// Actually receive the block
				nanoWallet.CreateAndPublishReceiveBlock(wallet, dbAccount.Address, msg.Hash, nil, nil)
			}()

		}
	}()

	// Create app
	app := chi.NewRouter()

	// Setup controller
	hc := controller.HttpController{Wallet: &nanoWallet, RpcClient: rpcClient, PowClient: pow}

	// HTTP Routes
	app.Use(middleware.Logger)
	app.Post("/", hc.Gateway)

	http.ListenAndServe(fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port), app)
}
