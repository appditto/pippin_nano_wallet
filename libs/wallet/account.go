package wallet

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/account"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
)

var ErrAccountNotFound = errors.New("account not found")
var ErrAccountExists = errors.New("account already exists")
var ErrUnableToCreateAccount = errors.New("unable to create account")

// Retrieve an account or adhoc account for a wallet
func (w *NanoWallet) GetAccount(wallet *ent.Wallet, address string) (*ent.Account, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	}

	// Determine if wallet is locked or not
	_, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	// Check if account exists
	acc, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.Address(address)).First(w.Ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrAccountNotFound
		}
		return nil, err
	}

	return acc, nil
}

// Create the next account in sequence, or at index
func (w *NanoWallet) AccountCreate(wallet *ent.Wallet, index *int) (*ent.Account, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	}

	// Obtain a lock, prevent concurrent calls
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("wallet:%s", wallet.ID.String()), time.Second*10, &database.LockRetryStrategy)
	if err != nil {
		return nil, database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	// Get seed
	seed, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	if index != nil {
		// See if account exists at index
		pub, priv, err := utils.KeypairFromSeed(seed, uint32(*index))
		if err != nil {
			return nil, err
		}
		address := utils.PubKeyToAddress(pub, w.Banano)
		exists, err := w.AccountExists(wallet, address)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrAccountExists
		}
		// Create it as an adhoc
		acc, err := w.DB.Account.Create().SetWallet(wallet).SetAddress(address).SetPrivateKey(hex.EncodeToString(priv)).Save(w.Ctx)
		if err != nil {
			return nil, err
		}
		return acc, nil
	}

	account, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.AccountIndexNotNil()).Order(ent.Desc(account.FieldAccountIndex)).First(w.Ctx)

	if err != nil {
		return nil, err
	}

	// We repeat as many times as necessary to avoid collisions
	runningIndex := *account.AccountIndex + 1
	for true {
		// Derive next account
		pub, _, err := utils.KeypairFromSeed(seed, uint32(runningIndex))
		if err != nil {
			return nil, err
		}
		address := utils.PubKeyToAddress(pub, w.Banano)
		exists, err := w.AccountExists(wallet, address)
		if err != nil {
			return nil, err
		}
		if exists {
			runningIndex++
			continue
		}
		newAcc, err := w.DB.Account.Create().SetWallet(wallet).SetAccountIndex(runningIndex).SetAddress(address).Save(w.Ctx)
		if err != nil {
			return nil, err
		}

		return newAcc, nil
	}
	return nil, ErrUnableToCreateAccount
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

	// Get seed
	seed, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	acc, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.AccountIndexNotNil()).Order(ent.Desc(account.FieldAccountIndex)).First(w.Ctx)

	if err != nil {
		return nil, err
	}

	nextIndex := *acc.AccountIndex + 1

	tx, err := w.DB.Tx(w.Ctx)
	var accounts []*ent.Account
	for i := 0; i < count; i++ {
		// Derive next account
		pub, _, err := utils.KeypairFromSeed(seed, uint32(nextIndex))
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		address := utils.PubKeyToAddress(pub, w.Banano)
		count, err := tx.Account.Query().Where(account.WalletID(wallet.ID), account.Address(address)).Count(w.Ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if count > 0 {
			// Iterate target index and roll back iteration count
			nextIndex++
			i--
			continue
		}
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
// Already existing account may be adhoc or non-adhoc
func (w *NanoWallet) AdhocAccountCreate(wallet *ent.Wallet, privKey ed25519.PrivateKey) (*ent.Account, error) {
	// Input validations
	if wallet == nil {
		return nil, ErrInvalidWallet
	} else if privKey == nil || len(privKey) != ed25519.PrivateKeySize {
		return nil, ErrInvalidPrivKey
	}

	// Obtain a lock, prevent concurrent calls
	lock, err := database.GetRedisDB().Locker.Obtain(w.Ctx, fmt.Sprintf("wallet:%s", wallet.ID.String()), time.Second*10, &database.LockRetryStrategy)
	if err != nil {
		return nil, database.ErrLockNotObtained
	}
	defer lock.Release(w.Ctx)

	// Determine if wallet is locked or not
	_, err = GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	// Derive address
	pub := privKey.Public().(ed25519.PublicKey)
	address := utils.PubKeyToAddress(pub, w.Banano)

	// See if account already exists
	acct, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.Address(address)).First(w.Ctx)
	if err == nil {
		return acct, nil
	} else if !ent.IsNotFound(err) {
		// Some unknown error we didn't expect
		return nil, err
	}

	// Create adhoc account
	adhocAcct, err := w.DB.Account.Create().SetWallet(wallet).SetAddress(address).SetPrivateKey(hex.EncodeToString(privKey)).Save(w.Ctx)
	if err != nil {
		return nil, err
	}

	return adhocAcct, nil
}

// Retrieve list of accounts on a wallet, if not locked
func (w *NanoWallet) AccountsList(wallet *ent.Wallet, limit int) ([]string, error) {
	if wallet == nil {
		return nil, ErrInvalidWallet
	}

	// Determine if wallet is locked or not
	_, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	// Get accounts
	accounts, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID)).Limit(limit).All(w.Ctx)
	if err != nil {
		return nil, err
	}

	// Concatenate them together as an array of addresses
	var addresses []string
	for _, acct := range accounts {
		addresses = append(addresses, acct.Address)
	}

	return addresses, nil
}

func (w *NanoWallet) AccountExists(wallet *ent.Wallet, address string) (bool, error) {
	if wallet == nil {
		return false, ErrInvalidWallet
	}

	// Determine if wallet is locked or not
	_, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return false, err
	}

	// Get accounts
	count, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.Address(address)).Count(w.Ctx)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
