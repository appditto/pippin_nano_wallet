package main

import (
	"context"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/appditto/pippin_nano_wallet/apps/server"
	"github.com/appditto/pippin_nano_wallet/libs/config"
	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/account"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	"github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"golang.org/x/term"
)

var Version = "dev"
var walletCmd *flag.FlagSet
var accountCmd *flag.FlagSet

func usage() {
	fmt.Println("General commands:")
	fmt.Printf("Usage %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Println("\n\nWallet commands:")
	fmt.Printf("Usage: %s wallet [options]\n", os.Args[0])
	fmt.Println("Options:")
	walletCmd.PrintDefaults()
	fmt.Println("\n\nAccount commands:")
	fmt.Printf("Usage: %s account [options]\n", os.Args[0])
	fmt.Println("Options:")
	accountCmd.PrintDefaults()
	return
}

func init() {
	walletCmd = flag.NewFlagSet("wallet", flag.ExitOnError)
	accountCmd = flag.NewFlagSet("account", flag.ExitOnError)
}

func getWallet(nanoWallet *wallet.NanoWallet, id string) *ent.Wallet {
	w, err := nanoWallet.GetWallet(id)
	if errors.Is(err, wallet.ErrWalletNotFound) {
		fmt.Println("Wallet not found")
		os.Exit(1)
	} else if err != nil {
		fmt.Printf("Failed to get wallet: %v\n", err)
		os.Exit(1)
	}
	return w
}

func RequireSeed(seed *string) {
	if *seed == "" {
		fmt.Println("--seed is required")
		os.Exit(1)
	} else if !utils.Validate64HexHash(*seed) {
		fmt.Println("Invalid seed")
		os.Exit(1)
	}
}

func RequireID(id *string, msg string) {
	if *id == "" {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func PasswordPrompt() string {
	for {
		var password string

		fmt.Print("➡️ Enter Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))

		if err != nil {
			fmt.Printf("\n⚠️ Error reading password\n")
			continue
		}

		password = strings.TrimSpace(string(bytePassword))

		var confirmedPassword string

		fmt.Print("\n➡️ Confirm Password: ")
		bytePassword, err = term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			fmt.Printf("\n⚠️ Error reading password\n")
			continue
		}

		confirmedPassword = strings.TrimSpace(string(bytePassword))

		if password != confirmedPassword {
			fmt.Printf("\n⚠️ Passwords do not match\n")
			continue
		}

		fmt.Println()
		return password
	}

}

func RequireUnlockedWallet(nanoWallet *wallet.NanoWallet, w *ent.Wallet, password *string) (alreadyUnlocked bool) {
	_, err := wallet.GetDecryptedKeyFromStorage(w, "seed")
	if errors.Is(err, wallet.ErrWalletLocked) {
		if *password == "" {
			fmt.Println("Wallet is locked, please provide password with --password")
			os.Exit(1)
		}
		// Unlock wallet
		ok, err := nanoWallet.UnlockWallet(w, *password)
		if err != nil || !ok {
			fmt.Printf("Failed to unlock wallet: %v\n", err)
			os.Exit(1)
		}
		return false
	} else if err != nil {
		fmt.Printf("Failed to get wallet: %v\n", err)
		os.Exit(1)
	}
	return true
}

func main() {
	showHelp := flag.Bool("help", false, "Show help")
	version := flag.Bool("version", false, "Display the version")
	startServer := flag.Bool("start-server", false, "Start the pippin server")
	flag.Parse()

	walletList := walletCmd.Bool("list", false, "List all wallets along with their accounts")
	walletCreate := walletCmd.Bool("create", false, "Create a new wallet")
	walletChangeSeed := walletCmd.Bool("change-seed", false, "Change the seed of a wallet")
	walletViewSeed := walletCmd.Bool("view-seed", false, "View the seed of a wallet (unsafe)")
	walletEncrypt := walletCmd.Bool("encrypt", false, "Encrypt a wallet with a password")
	walletDecryt := walletCmd.Bool("decrypt", false, "Decrypt a wallet, remove password requirement")
	// Options that may apply to multiple commands
	walletId := walletCmd.String("id", "", "Target wallet ID")
	walletSeed := walletCmd.String("seed", "", "Specify a seed to use when creating/changing wallet (optional for create)")
	walletPassword := walletCmd.String("password", "", "Specify a password to use if the wallet is locked")
	walletAllKeys := walletCmd.Bool("all-keys", false, "Show all priv/pub keys for accounts on this wallet")

	// For accounts
	accountCreate := accountCmd.Bool("create", false, "Create a new account")
	// Options that may apply to multiple commands
	accountWalletId := accountCmd.String("id", "", "Target wallet ID")
	accountWalletPassword := accountCmd.String("password", "", "Specify a password to use if the wallet is locked")
	accountIndex := accountCmd.Int("index", 0, "Specify an index to use when creating account (optional, cannot be used with --count)")
	accountCount := accountCmd.Int("count", 0, "Specify how many accounts to create (optional, cannot be used with --index)")
	accountKey := accountCmd.String("key", "", "Specify a private key to use when creating account (optional, cannot be used with --index or --count)")
	repairAdhoc := accountCmd.Bool("repair-adhoc", false, "Repair adhoc accounts")

	if *showHelp {
		usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("Pippin version: %s\n", Version)
		os.Exit(0)
	}

	if *startServer {
		server.StartPippinServer()
		os.Exit(0)
	}

	if len(os.Args) < 2 {
		fmt.Println("expected 'wallet' subcommands")
		os.Exit(1)
	}

	// Read yaml configuration
	conf, err := config.ParsePippinConfig()
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
		os.Exit(1)
	}

	// Setup database conn
	ctx := context.Background()
	dbconn, err := database.GetSqlDbConn(false)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
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

	switch os.Args[1] {

	case "wallet":
		walletCmd.Parse(os.Args[2:])
		// ** wallet --list
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
			// ** wallet --create (--seed)
		} else if *walletCreate {
			var seed string
			if *walletSeed == "" {
				fmt.Println("Generating secure seed...")
				seed, err = utils.GenerateSeed(nil)
				if err != nil {
					fmt.Printf("Secue random source may not be available on your OS\n")
					fmt.Printf("Failed to generate seed: %v\n", err)
					os.Exit(1)
				}
			} else {
				seed = *walletSeed
			}
			if !utils.Validate64HexHash(seed) {
				fmt.Println("Invalid seed")
				os.Exit(1)
			}
			// Create wallet
			w, err := nanoWallet.WalletCreate(seed)
			if err != nil {
				fmt.Printf("Failed to create wallet: %v\n", err)
				os.Exit(1)
			}
			// Retrieve wallet
			w, err = nanoWallet.GetWallet(w.ID.String())
			acct, err := w.QueryAccounts().First(ctx)
			if err != nil {
				fmt.Printf("Failed to get account for wallet: %v\n", err)
				os.Exit(1)
			}
			//print(f"Wallet created, ID: {wallet.id}\nFirst account: {new_acct}")
			fmt.Printf("Wallet created, ID: %s\n", w.ID.String())
			fmt.Printf("First account: %s\n", acct.Address)
			// ** wallet --change-seed (--id --seed)
		} else if *walletChangeSeed {
			RequireSeed(walletSeed)
			RequireID(walletId, "--id is required for --change-seed")
			w := getWallet(&nanoWallet, *walletId)
			RequireUnlockedWallet(&nanoWallet, w, walletPassword)

			// Change seed
			newest, err := nanoWallet.WalletChangeSeed(w, *walletSeed)
			if err != nil {
				fmt.Printf("Failed to change seed: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Wallet seed changed, newest account: %s with index %d\n", newest.Address, *newest.AccountIndex)
			// ** wallet --view-seed --id
		} else if *walletViewSeed {
			RequireID(walletId, "--id is required for --view-seed")
			w := getWallet(&nanoWallet, *walletId)
			alreadyUnlocked := RequireUnlockedWallet(&nanoWallet, w, walletPassword)
			seed, err := wallet.GetDecryptedKeyFromStorage(w, "seed")
			if err != nil {
				fmt.Printf("Failed to get seed: %v\n", err)
				os.Exit(1)
			}
			if !alreadyUnlocked {
				// Re-lock the wallet
				err = nanoWallet.LockWallet(w)
				if err != nil {
					fmt.Printf("Failed to re-lock wallet: %v\n", err)
					os.Exit(1)
				}
			}
			fmt.Printf("Seed: %s\n", strings.ToUpper(seed))
			// Get accounts
			if *walletAllKeys {
				accounts, err := w.QueryAccounts().Where(account.AccountIndexNotNil()).All(ctx)
				if err != nil {
					fmt.Printf("Failed to get accounts for wallet: %v\n", err)
					os.Exit(1)
				}
				for _, a := range accounts {
					_, priv, _ := utils.KeypairFromSeed(seed, uint32(*a.AccountIndex))
					asStr := strings.ToUpper(hex.EncodeToString(priv))[:64]
					fmt.Printf("Account: %s PrivKey: %s\n", a.Address, asStr)
				}
			}
			// Show adhoc accounts no matter what
			adhocAccounts, err := w.QueryAccounts().Where(account.PrivateKeyNotNil()).All(ctx)
			if err != nil {
				fmt.Printf("Failed to get adhoc accounts for wallet: %v\n", err)
				os.Exit(1)
			}
			if len(adhocAccounts) > 0 {
				fmt.Println("AdHoc Accounts:")
			}
			for _, a := range adhocAccounts {
				asStr := strings.ToUpper(*a.PrivateKey)[:64]
				fmt.Printf("Account: %s PrivKey: %s\n", a.Address, asStr)
			}
			// ** --encrypt --id
		} else if *walletEncrypt {
			RequireID(walletId, "--id is required for --encrypt")
			w := getWallet(&nanoWallet, *walletId)
			if w.Encrypted {
				fmt.Println("Wallet is already encrypted")
				os.Exit(1)
			}
			// Prompt for username and password
			password := PasswordPrompt()
			if password == "" {
				fmt.Println("Password cannot be empty")
				os.Exit(1)
			}
			// Encrypt wallet
			ok, err := nanoWallet.EncryptWallet(w, password)
			if err != nil || !ok {
				fmt.Printf("Failed to encrypt wallet: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Wallet encrypted")
		} else if *walletDecryt {
			RequireID(walletId, "--id is required for --decrypt")
			w := getWallet(&nanoWallet, *walletId)
			if !w.Encrypted {
				fmt.Println("Wallet is not encrypted")
				os.Exit(1)
			}
			// Prompt for username and password
			password := PasswordPrompt()
			if password == "" {
				fmt.Println("Password cannot be empty")
				os.Exit(1)
			}
			// Decrypt wallet
			RequireUnlockedWallet(&nanoWallet, w, &password)
			ok, err := nanoWallet.EncryptWallet(w, "")
			if err != nil || !ok {
				fmt.Printf("Failed to decrypt wallet: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Wallet decrypted")
		} else {
			usage()
		}
	case "account":
		accountCmd.Parse(os.Args[2:])
		// ** account --create (--id --index --count --key)
		if *accountCreate {
			RequireID(accountWalletId, "--id is required for --create")
			w := getWallet(&nanoWallet, *accountWalletId)
			if *accountIndex > 0 && *accountCount > 0 {
				fmt.Println("Cannot specify both --index and --count")
				os.Exit(1)
			}
			if *accountIndex > 0 && *accountKey != "" {
				fmt.Println("Cannot specify both --index and --key")
				os.Exit(1)
			}
			if *accountCount > 0 && *accountKey != "" {
				fmt.Println("Cannot specify both --count and --key")
				os.Exit(1)
			}
			alreadyUnlocked := RequireUnlockedWallet(&nanoWallet, w, accountWalletPassword)
			if !alreadyUnlocked {
				defer nanoWallet.LockWallet(w)
			}
			// Determine type
			if *accountIndex > 0 {
				// Create account at index
				acc, err := nanoWallet.AccountCreate(w, accountIndex)
				if err != nil {
					fmt.Printf("Failed to create account: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Account created: %s\n", acc.Address)
			} else if *accountCount > 0 {
				accs, err := nanoWallet.AccountsCreate(w, *accountCount)
				if err != nil {
					fmt.Printf("Failed to create accounts: %v\n", err)
					os.Exit(1)
				}
				for _, acc := range accs {
					fmt.Printf("Account %d created: %s\n", *acc.AccountIndex, acc.Address)
				}
			} else if *accountKey != "" {
				if !utils.Validate64HexHash(*accountKey) {
					fmt.Println("Invalid private key")
					os.Exit(1)
				}
				asByte, err := hex.DecodeString(*accountKey)
				if err != nil {
					fmt.Printf("Failed to decode private key: %v\n", err)
					os.Exit(1)
				}
				priv, err := ed25519.NewKeyFromSeed(asByte)
				if err != nil {
					fmt.Printf("Failed to create private key: %v\n", err)
					os.Exit(1)
				}
				// Create adhoc account
				acc, err := nanoWallet.AdhocAccountCreate(w, priv)
				if err != nil {
					fmt.Printf("Failed to create adhoc account: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Adhoc account created: %s\n", acc.Address)
			} else {
				// Create account
				acc, err := nanoWallet.AccountCreate(w, nil)
				if err != nil {
					fmt.Printf("Failed to create account: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Account %d created: %s\n", *acc.AccountIndex, acc.Address)
			}
		} else if *repairAdhoc {
			// Get all accounts with 64-length private keys
			accounts, err := nanoWallet.DB.Account.Query().Where(account.PrivateKeyNotNil()).All(ctx)
			if err != nil {
				fmt.Printf("Failed to get accounts: %v\n", err)
				os.Exit(1)
			}
			for _, a := range accounts {
				if len(*a.PrivateKey) != 64 {
					continue
				}
				// Convert to ed25519 key
				asByte, err := hex.DecodeString(*a.PrivateKey)
				if err != nil {
					fmt.Printf("Failed to decode private key: %v\n", err)
					os.Exit(1)
				}
				priv, err := ed25519.NewKeyFromSeed(asByte)
				if err != nil {
					fmt.Printf("Failed to create private key: %v\n", err)
					os.Exit(1)
				}
				// Update account
				_, err = nanoWallet.DB.Account.UpdateOneID(a.ID).SetPrivateKey(hex.EncodeToString(priv)).Save(ctx)
				if err != nil {
					fmt.Printf("Failed to update account: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("Account %s repaired\n", a.Address)
			}
		}
	default:
		fmt.Println("expected 'foo' or 'bar' subcommands")
		os.Exit(1)
	}
}
