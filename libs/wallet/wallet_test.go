package wallet

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/config"
	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/account"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	nanorpc "github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var MockWallet *NanoWallet
var bananoWallet *NanoWallet

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	os.Setenv("MOCK_REDIS", "true")
	defer os.Unsetenv("MOCK_REDIS")
	dbconn, _ := database.GetSqlDbConn(true)
	client, _ := database.NewEntClient(dbconn)
	defer client.Close()
	if err := client.Schema.Create(context.TODO()); err != nil {
		panic(err)
	}
	os.Setenv("HOME", ".testdata")
	defer os.Unsetenv("HOME")
	defer os.RemoveAll(".testdata")
	config, _ := config.ParsePippinConfig()
	rpcclient := nanorpc.NewRPCClient("/mockrpcendpoint")
	powClient := pow.NewPippinPow([]string{}, "", "", 30)
	MockWallet = &NanoWallet{
		DB:         client,
		Ctx:        context.TODO(),
		Banano:     false,
		RpcClient:  rpcclient,
		WorkClient: powClient,
		Config:     config,
	}
	bananoWallet = &NanoWallet{
		DB:         client,
		Ctx:        context.TODO(),
		Banano:     true,
		RpcClient:  rpcclient,
		Config:     config,
		WorkClient: powClient,
	}

	return m.Run()
}

func TestGetWallet(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("55555540e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Retrieve walelt
	gotten, err := MockWallet.GetWallet(wallet.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, gotten.ID, wallet.ID)

	// Check not found error
	_, err = MockWallet.GetWallet("invalid")
	assert.ErrorIs(t, ErrWalletNotFound, err)
	_, err = MockWallet.GetWallet(uuid.NewString())
	assert.ErrorIs(t, ErrWalletNotFound, err)
}

func TestGetWallets(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("c1cc82790c35cc98de93958555eea5691754b6d97f841725238717f9ae5e43bb"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Retrieve wallets
	gotten, err := MockWallet.GetWallets()
	assert.Nil(t, err)
	// Creat a list of UUID strings based on result
	var walletIDs []string
	for _, w := range gotten {
		walletIDs = append(walletIDs, w.ID.String())
	}
	assert.Contains(t, walletIDs, wallet.ID.String())
}

func TestWalletCreate(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("8d729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Ensure account is created
	account, err := MockWallet.DB.Account.Query().Where(account.WalletID(wallet.ID)).First(MockWallet.Ctx)
	assert.Nil(t, err)
	assert.Equal(t, "nano_1efa1gxbitary1urzix9h13nkzadtz71n3auyj7uztb8i4qbtipu8cxz61ee", account.Address)
	assert.Equal(t, 0, *account.AccountIndex)
	assert.Equal(t, true, account.Work)
}

func TestWalletCreateBanano(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("bb729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))

	wallet, err := bananoWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Ensure account is created
	account, err := bananoWallet.DB.Account.Query().Where(account.WalletID(wallet.ID)).First(bananoWallet.Ctx)
	assert.Nil(t, err)
	assert.Equal(t, "ban_33mhuqjxr166czm4y37xk7emfnt4zogxmqrhbfxyngrkbdchmpsk6qehhm3n", account.Address)
	assert.Equal(t, 0, *account.AccountIndex)
	assert.Equal(t, true, account.Work)
}

func TestWalletCreateBadInput(t *testing.T) {
	// Empty seed
	_, err := MockWallet.WalletCreate("")
	assert.ErrorIs(t, ErrInvalidSeed, err)

	// Invalid seed
	_, err = MockWallet.WalletCreate("invalid seed")
	assert.ErrorIs(t, ErrInvalidSeed, err)
}

func TestWalletDestroy(t *testing.T) {
	// Create a test wallet
	seed, _ := utils.GenerateSeed(strings.NewReader("783c75f57c76937b2bab1e0ada730d1386bacfa06258ddebfcc976b36c0e5549"))
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create relationships to ensure cascade delete works

	// Create some accounts
	_, err = MockWallet.AccountsCreate(wallet, 10)
	assert.Nil(t, err)

	// Create adhoc accounts
	_, priv, _ := utils.KeypairFromSeed(seed, 100)
	acc, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	// Create a block object
	_, err = MockWallet.DB.Block.Create().SetAccount(acc).SetBlock(map[string]interface{}{
		"block": "hello",
	}).SetBlockHash("sdadasdasd").SetSubtype("change").Save(MockWallet.Ctx)
	assert.Nil(t, err)

	err = MockWallet.WalletDestroy(wallet)
	assert.Nil(t, err)

	err = MockWallet.WalletDestroy(nil)
	assert.ErrorIs(t, ErrInvalidWallet, err)

	err = MockWallet.WalletDestroy(&ent.Wallet{})
	assert.NotNil(t, err)
}

func TestWalletInfo(t *testing.T) {
	// Create a test wallet
	seed, _ := utils.GenerateSeed(strings.NewReader("43ae06048b189e8a15da9765d8ce21edbf2d34eb7b1b7fb928e028e3fb416d53"))
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create some accounts
	_, err = MockWallet.AccountsCreate(wallet, 10)
	assert.Nil(t, err)

	// Create adhoc accounts
	_, priv, _ := utils.KeypairFromSeed(seed, 100)
	_, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	info, err := MockWallet.WalletInfo(wallet)
	assert.Nil(t, err)

	assert.Equal(t, info.AccountsCount, 12)
	assert.Equal(t, info.AdhocCount, 1)
	assert.Equal(t, info.DeterministicCount, 11)
	assert.Equal(t, info.DeterministicIndex, 10)
}

func TestWalletRepresentativeSet(t *testing.T) {
	// Create a test wallet
	seed, _ := utils.GenerateSeed(strings.NewReader("1f447808006c0c50d0193eda500ff482d08effa9187dfe2d57a1c5009b2d5f6c"))
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	wallet, err = MockWallet.GetWallet(wallet.ID.String())
	assert.Nil(t, err)

	err = MockWallet.WalletRepresentativeSet(wallet, "nano_1efa1gxbitary1urzix9h13nkzadtz71n3auyj7uztb8i4qbtipu8cxz61ee", false, nil)
	assert.Nil(t, err)

	// Retrieve wallet
	wallet, err = MockWallet.GetWallet(wallet.ID.String())
	assert.Nil(t, err)
	assert.Equal(t, "nano_1efa1gxbitary1urzix9h13nkzadtz71n3auyj7uztb8i4qbtipu8cxz61ee", *wallet.Representative)
}

func TestWalletChangeSeed(t *testing.T) {
	// Create a test wallet
	seed, _ := utils.GenerateSeed(strings.NewReader("94e6c473cf3d539822e073f64d31a248621ef28be595fc497adf322f29a3d9e4"))
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	wallet, err = MockWallet.GetWallet(wallet.ID.String())
	assert.Nil(t, err)

	acc, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	oldAddress := acc.Address

	// Change seed
	newest, err := MockWallet.WalletChangeSeed(wallet, "c0e319472702d7cbe728ad05647395498a6ad498b9ae7e36a33cc37fef60f27a")
	assert.Nil(t, err)
	assert.Equal(t, "nano_33fj9exam1ppgzaco6hjd7z1nnapnf4gh3ech4fbfkr6eotb7bui3qzukt73", newest.Address)
	assert.NotEqual(t, oldAddress, newest.Address)

	// Test with locked wallet
	_, err = MockWallet.EncryptWallet(wallet, "password")
	assert.Nil(t, err)

	_, err = MockWallet.WalletChangeSeed(wallet, "0c07c237c3e4254aeb2cae8ddd48eb9f13f393a99aefa3f012231e9688641e58")
	assert.ErrorIs(t, ErrWalletLocked, err)

	// Unlock wallet
	_, err = MockWallet.UnlockWallet(wallet, "password")
	assert.Nil(t, err)

	// Change seed
	newest, err = MockWallet.WalletChangeSeed(wallet, "e0e87bf97ac01f4428864aa752a2d7acb9c2ca99ea2e69296c8507d5d71408fb")
	assert.Nil(t, err)
	assert.Equal(t, "nano_16rxu414wbt34tyn7yugup99s4xt1htrfufkwjce19ezfwfbmzrf343ynyoi", newest.Address)
}
