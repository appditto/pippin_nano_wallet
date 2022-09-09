package wallet

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/account"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/adhocaccount"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
)

type NanoWallet struct {
	DB     *ent.Client
	Ctx    context.Context
	Banano bool
}

var ErrInvalidSeed = errors.New("invalid seed")
var ErrInvalidWallet = errors.New("invalid wallet")
var ErrInvalidPrivKey = errors.New("invalid private key")
var ErrInvalidAccountCount = errors.New("invalid count")

// Creates a new wallet with provided seed
func (w *NanoWallet) WalletCreate(seed string) (*ent.Wallet, error) {
	if !utils.Validate64HexHash(seed) {
		return nil, ErrInvalidSeed
	}

	// Create wallet and first account
	tx, err := w.DB.Tx(w.Ctx)
	if err != nil {
		return nil, err
	}
	wallet, err := tx.Wallet.Create().SetSeed(seed).Save(w.Ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Derive first account
	pub, _, err := utils.KeypairFromSeed(seed, 0)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	address := utils.PubKeyToAddress(pub, w.Banano)

	// Create first account
	_, err = tx.Account.Create().SetWallet(wallet).SetAccountIndex(0).SetAddress(address).Save(w.Ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

// Create the next account in sequence
func (w *NanoWallet) AccountCreate(wallet *ent.Wallet) (*ent.Account, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	}

	// Obtain a lock, prevent concurrent calls
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("wallet:%s", wallet.ID.String()), time.Second*10, &database.LockRetryStrategy)
	if err != nil {
		return nil, database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	account, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID)).Order(ent.Desc(account.FieldAccountIndex)).First(w.Ctx)

	if err != nil {
		return nil, err
	}

	// Derive next account
	pub, _, err := utils.KeypairFromSeed(wallet.Seed, uint32(account.AccountIndex)+1)
	if err != nil {
		return nil, err
	}
	address := utils.PubKeyToAddress(pub, w.Banano)
	newAcc, err := w.DB.Account.Create().SetWallet(wallet).SetAccountIndex(account.AccountIndex + 1).SetAddress(address).Save(w.Ctx)
	if err != nil {
		return nil, err
	}

	return newAcc, nil
}

func (w *NanoWallet) AccountsCreate(wallet *ent.Wallet, count int) ([]*ent.Account, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	} else if count < 1 {
		return nil, ErrInvalidAccountCount
	}

	// Obtain a lock, prevent concurrent calls
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("wallet:%s", wallet.ID.String()), time.Second*10, &database.LockRetryStrategy)
	if err != nil {
		return nil, database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	account, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID)).Order(ent.Desc(account.FieldAccountIndex)).First(w.Ctx)

	if err != nil {
		return nil, err
	}

	nextIndex := account.AccountIndex + 1

	tx, err := w.DB.Tx(w.Ctx)
	var accounts []*ent.Account
	for i := 0; i < count; i++ {
		// Derive next account
		pub, _, err := utils.KeypairFromSeed(wallet.Seed, uint32(nextIndex))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		address := utils.PubKeyToAddress(pub, w.Banano)
		acct, err := tx.Account.Create().SetWallet(wallet).SetAccountIndex(nextIndex).SetAddress(address).Save(w.Ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		accounts = append(accounts, acct)
		nextIndex++
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

// Create an adhoc account. Returns adhocaccount or account if already exists
func (w *NanoWallet) AdhocAccountCreate(wallet *ent.Wallet, privKey ed25519.PrivateKey) (*ent.AdhocAccount, *ent.Account, error) {
	// Input validations
	if wallet == nil {
		return nil, nil, ErrInvalidWallet
	} else if privKey == nil || len(privKey) != ed25519.PrivateKeySize {
		return nil, nil, ErrInvalidPrivKey
	}

	// Obtain a lock, prevent concurrent calls
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("wallet:%s", wallet.ID.String()), time.Second*10, &database.LockRetryStrategy)
	if err != nil {
		return nil, nil, database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	// Derive address
	pub := privKey.Public().(ed25519.PublicKey)
	address := utils.PubKeyToAddress(pub, w.Banano)

	// See if account already exists in either table
	acct, err := w.DB.Account.Query().Where(account.Address(address)).First(w.Ctx)
	if err == nil {
		return nil, acct, nil
	}
	// Check adhoc accounts table
	adhocAcct, err := w.DB.AdhocAccount.Query().Where(adhocaccount.Address(address)).First(w.Ctx)
	if err == nil {
		return adhocAcct, nil, nil
	}

	// Create adhoc account
	adhocAcct, err = w.DB.AdhocAccount.Create().SetWallet(wallet).SetAddress(address).SetPrivateKey(hex.EncodeToString(privKey)).Save(w.Ctx)
	if err != nil {
		return nil, nil, err
	}

	return adhocAcct, nil, nil
}
