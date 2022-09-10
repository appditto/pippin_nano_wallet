package wallet

import (
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent/adhocaccount"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/stretchr/testify/assert"
)

func TestEncryptWallet(t *testing.T) {
	// Create a wallet
	seed, err := utils.GenerateSeed(strings.NewReader("88888340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Nil(t, err)

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create some adhoc accounts
	_, priv, _ := ed25519.GenerateKey(strings.NewReader("1111111111111111111111111111111111111111111111111111111111111111"))
	_, _, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)
	_, priv, _ = ed25519.GenerateKey(strings.NewReader("2111111111111111111111111111111111111111111111111111111111111111"))
	_, _, err = MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	// Encrypt the wallet
	password := "mypassword"
	encrypted, err := MockWallet.EncryptWallet(wallet, password)
	assert.Nil(t, err)
	assert.True(t, encrypted)

	// Ensure the wallet is encrypted
	assert.Nil(t, err)
	assert.True(t, wallet.Encrypted)
	assert.NotEqual(t, seed, wallet.Seed)

	// Ensure adhoc keys are encrypted
	adhocs, err := MockWallet.DB.AdhocAccount.Query().Where(adhocaccount.WalletID(wallet.ID)).All(MockWallet.Ctx)
	assert.Nil(t, err)
	assert.Len(t, adhocs, 2)
	for _, adhoc := range adhocs {
		// If private key has this many bits its encrypted
		assert.True(t, len(adhoc.PrivateKey) > 128)
	}

	// Ensure we get an error acting on it since it's locked
	_, err = MockWallet.EncryptWallet(wallet, password)
	assert.ErrorIs(t, ErrWalletLocked, err)

	// Test bad input
	_, err = MockWallet.EncryptWallet(nil, "abcd1234")
	assert.ErrorIs(t, ErrInvalidWallet, err)

	wallet.Encrypted = false
	_, err = MockWallet.EncryptWallet(wallet, "")
	assert.ErrorIs(t, ErrBadPassword, err)
}

func TestSetAndGetDecryptedFromStorage(t *testing.T) {
	// Create a wallet
	seed, err := utils.GenerateSeed(strings.NewReader("22222340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Nil(t, err)

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Wallet isn't encrypted to begin with, so we expect to get the seed from relational DB
	retrievedSeed, err := GetDecryptedKeyFromStorage(wallet, "seed")
	assert.Nil(t, err)
	assert.Equal(t, seed, retrievedSeed)

	// Wallet not locked error
	err = SetDecryptedKeyToStorage(wallet, "seed", seed)
	assert.ErrorIs(t, ErrWalletNotLocked, err)

	// Set decrypted seed
	wallet.Encrypted = true
	err = SetDecryptedKeyToStorage(wallet, "seed", "1234")
	assert.Nil(t, err)

	// Get decrypted seed, we expect redis version to be 1234
	retrievedSeed, err = GetDecryptedKeyFromStorage(wallet, "seed")
	assert.Nil(t, err)
	assert.Equal(t, "1234", retrievedSeed)
}

func TestUnlockWallet(t *testing.T) {
	// Create a wallet
	seed, err := utils.GenerateSeed(strings.NewReader("33333340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040b"))
	assert.Nil(t, err)

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create some adhoc accounts
	_, priv, _ := ed25519.GenerateKey(strings.NewReader("1111111111111111111111111111111111111111111111111111111111111111"))
	acc1, _, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)
	_, priv, _ = ed25519.GenerateKey(strings.NewReader("2111111111111111111111111111111111111111111111111111111111111111"))
	acc2, _, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	// Encrypt the wallet
	password := "mypassword"
	encrypted, err := MockWallet.EncryptWallet(wallet, password)
	assert.Nil(t, err)
	assert.True(t, encrypted)

	// Unlock with bad password
	unlocked, err := MockWallet.UnlockWallet(wallet, "hunter2")
	assert.ErrorIs(t, ErrBadPassword, err)
	assert.False(t, unlocked)

	// Unlock with good password
	unlocked, err = MockWallet.UnlockWallet(wallet, password)
	assert.Nil(t, err)
	assert.True(t, unlocked)

	// Check that key exists
	key, err := GetDecryptedKeyFromStorage(wallet, acc1.Address)
	assert.Nil(t, err)
	assert.Equal(t, acc1.PrivateKey, key)
	key, _ = GetDecryptedKeyFromStorage(wallet, acc2.Address)
	assert.Equal(t, acc2.PrivateKey, key)

	// Check not locked error
	wallet.Encrypted = false
	_, err = MockWallet.UnlockWallet(wallet, password)
	assert.ErrorIs(t, ErrWalletNotLocked, err)
}
