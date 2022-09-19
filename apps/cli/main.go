package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/appditto/pippin_nano_wallet/libs/config"
	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	"github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"k8s.io/klog/v2"
)

var Version = "dev"
var walletCmd *flag.FlagSet

func usage() {
	fmt.Println("General commands:")
	fmt.Printf("Usage %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("\nWallet commands:")
	fmt.Printf("Usage: %s wallet [options]\n", os.Args[0])
	fmt.Println("Options:")
	walletCmd.PrintDefaults()
	return
}

func init() {
	walletCmd = flag.NewFlagSet("wallet", flag.ExitOnError)
}

func main() {
	showHelp := flag.Bool("help", false, "Show help")
	version := flag.Bool("version", false, "Display the version")

	flag.Parse()

	walletList := walletCmd.Bool("list", false, "List all wallets along with their accounts")

	if *showHelp {
		usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("Pippin version: %s\n", Version)
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("expected 'wallet' subcommands")
		os.Exit(1)
	}

	// Read yaml configuration
	conf, err := config.ParsePippinConfig()
	if err != nil {
		klog.Fatalf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	// Setup database conn
	ctx := context.Background()
	dbconn, err := database.GetSqlDbConn(false)
	if err != nil {
		klog.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	entClient, err := database.NewEntClient(dbconn)
	if err != nil {
		fmt.Printf("Failed to create ent client: %v", err)
		os.Exit(1)
	}
	defer entClient.Close()

	// Run migrations
	if err := entClient.Schema.Create(ctx); err != nil {
		fmt.Printf("Failed to run migrations: %v", err)
		os.Exit(1)
	}

	// Setup RPC handlers
	rpcClient := rpc.RPCClient{
		Url: conf.Server.NodeRpcUrl,
	}

	// Setup pow client
	pow := pow.NewPippinPow(conf.Wallet.WorkPeers, utils.GetEnv("BPOW_KEY", ""), utils.GetEnv("BPOW_URL", ""))

	// Setup nano wallet instance with DB, options, etc.
	nanoWallet := wallet.NanoWallet{
		DB:         entClient,
		Ctx:        ctx,
		Banano:     conf.Wallet.Banano,
		RpcClient:  &rpcClient,
		WorkClient: pow,
		Config:     conf,
	}

	switch os.Args[1] {

	case "wallet":
		walletCmd.Parse(os.Args[2:])
		if *walletList {
			// Get all wallets
			wallets, err := nanoWallet.GetWallets()
			if err != nil {
				fmt.Printf("Failed to get wallets: %v\n", err)
				os.Exit(1)
			}
			for idx, w := range wallets {
				fmt.Printf("Wallet ID: %s\n", w.ID.String())
				accounts, err := w.QueryAccounts().All(ctx)
				if err != nil {
					fmt.Printf("Failed to get accounts for wallet: %v\n", err)
					os.Exit(1)
				}
				for _, a := range accounts {
					fmt.Printf("Account: %s\n", a.Address)
				}
				if idx < len(wallets)-1 {
					fmt.Println("----------------------------")
				}
			}
		} else {
			usage()
		}
	default:
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}
}
