package wallet

import (
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/stretchr/testify/assert"
)

func TestGetBlockFromDatabase(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("42f567f406715877a8d60c6890b4655e8e0fe64f3f82089fe76899f180cd4021"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Create an account and an adhoc account
	account, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)

	_, priv, _ := ed25519.GenerateKey(strings.NewReader("1f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040d"))
	adhocAcct, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	// Create some blocks for each
	// Create a block object
	_, err = MockWallet.DB.Block.Create().SetAccount(account).SetBlock(map[string]interface{}{
		"block": "hello",
	}).SetBlockHash("abc").SetSubtype("send").SetSendID("abc").Save(MockWallet.Ctx)
	assert.Nil(t, err)
	_, err = MockWallet.DB.Block.Create().SetAccount(adhocAcct).SetBlock(map[string]interface{}{
		"block": "world",
	}).SetBlockHash("def").SetSubtype("change").SetSendID("def").Save(MockWallet.Ctx)
	assert.Nil(t, err)

	// Retrieve block
	block, err := MockWallet.GetBlockFromDatabase(wallet, account.Address, "abc")
	assert.Nil(t, err)
	assert.Equal(t, "hello", block.Block["block"])
	block, err = MockWallet.GetBlockFromDatabase(wallet, adhocAcct.Address, "def")
	assert.Nil(t, err)
	assert.Equal(t, "world", block.Block["block"])
	block, err = MockWallet.GetBlockFromDatabase(wallet, account.Address, "nonexistent")
	assert.ErrorIs(t, err, ErrBlockNotFound)
}
