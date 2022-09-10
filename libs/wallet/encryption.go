package wallet

import (
	"errors"

	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/adhocaccount"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/go-redis/redis/v9"
)

var ErrWalletLocked = errors.New("wallet is locked")
var ErrBadPassword = errors.New("bad password")
var ErrWalletNotLocked = errors.New("wallet not locked")

// This is encrypted wallets and adhoc accounts
// We store decrypted seeds in redis, while encrypted ones are stored in the database
// Encryption is an optional behavior

// Encrypt entrypoint
// If password is a blank string, we disable encryption for the wallet
// If password is not a blank string, we enable encryption for the wallet
func (w *NanoWallet) EncryptWallet(wallet *ent.Wallet, password string) (bool, error) {
	if wallet == nil {
		return false, ErrInvalidWallet
	} else if !wallet.Encrypted && password == "" {
		// Wallet is not encrypted and no password is set
		return false, ErrBadPassword
	}

	var seed string
	var err error
	if wallet.Encrypted {
		// Retrieve decrypted seed from storage, wallet has to be unlocked
		seed, err = GetDecryptedKeyFromStorage(wallet, "seed")
		if err != nil {
			return false, err
		}
	} else {
		// Wallet is not encrypted, so we can just use the seed
		seed = wallet.Seed
	}

	if password == "" {
		// Unlock wallet
		tx, err := w.DB.Tx(w.Ctx)
		if err != nil {
			return false, err
		}
		_, err = tx.Wallet.UpdateOne(wallet).SetEncrypted(false).SetSeed(seed).Save(w.Ctx)
		if err != nil {
			tx.Rollback()
			return false, err
		}
		adhocAccts, err := tx.AdhocAccount.Query().Where(adhocaccount.WalletID(wallet.ID)).All(w.Ctx)
		if err != nil {
			tx.Rollback()
			return false, err
		}
		for _, acct := range adhocAccts {
			key, err := GetDecryptedKeyFromStorage(wallet, acct.Address)
			if err != nil {
				tx.Rollback()
				return false, err
			}
			_, err = tx.AdhocAccount.UpdateOne(acct).SetPrivateKey(key).Save(w.Ctx)
			if err != nil {
				tx.Rollback()
				return false, err
			}
		}
		err = tx.Commit()
		if err != nil {
			return false, err
		}
		wallet.Encrypted = false
		wallet.Seed = seed
		database.GetRedisDB().Del(wallet.ID.String())
		return true, nil
	}

	// We are updating the password
	crypter := utils.NewAesCrypt(password)
	encryptedSeed, err := crypter.Encrypt(seed)
	if err != nil {
		return false, err
	}

	tx, err := w.DB.Tx(w.Ctx)
	if err != nil {
		return false, err
	}
	_, err = tx.Wallet.UpdateOne(wallet).SetEncrypted(true).SetSeed(encryptedSeed).Save(w.Ctx)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	// Encrypt all adhoc private keys
	adhocAccts, err := tx.AdhocAccount.Query().Where(adhocaccount.WalletID(wallet.ID)).All(w.Ctx)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	for _, acct := range adhocAccts {
		encryptedKey, err := crypter.Encrypt(acct.PrivateKey)
		if err != nil {
			return false, err
		}
		_, err = tx.AdhocAccount.UpdateOne(acct).SetPrivateKey(encryptedKey).Save(w.Ctx)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	wallet.Encrypted = true
	wallet.Seed = encryptedSeed
	return true, nil
}

func (w *NanoWallet) LockWallet(wallet *ent.Wallet) error {
	if wallet == nil {
		return ErrInvalidWallet
	} else if !wallet.Encrypted {
		return ErrWalletNotLocked
	}

	database.GetRedisDB().Del(wallet.ID.String())
	return nil
}

func (w *NanoWallet) UnlockWallet(wallet *ent.Wallet, password string) (bool, error) {
	if wallet == nil {
		return false, ErrInvalidWallet
	} else if !wallet.Encrypted {
		return false, ErrWalletNotLocked
	}

	crypter := utils.NewAesCrypt(password)
	seed, err := crypter.Decrypt(wallet.Seed)
	if err != nil {
		return false, ErrBadPassword
	}

	err = SetDecryptedKeyToStorage(wallet, "seed", seed)
	if err != nil {
		return false, err
	}

	// Every adhoc account gets decrypted too
	adhocAccts, err := w.DB.AdhocAccount.Query().Where(adhocaccount.WalletID(wallet.ID)).All(w.Ctx)
	if err != nil {
		return false, err
	}
	for _, acct := range adhocAccts {
		key, err := crypter.Decrypt(acct.PrivateKey)
		if err != nil {
			return false, err
		}
		err = SetDecryptedKeyToStorage(wallet, acct.Address, key)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

// Retrieve decrypted key from storage if it exists
func GetDecryptedKeyFromStorage(wallet *ent.Wallet, key string) (string, error) {
	if wallet == nil {
		return "", ErrInvalidWallet
	} else if !wallet.Encrypted {
		return wallet.Seed, nil
	}

	key, err := database.GetRedisDB().Hget(wallet.ID.String(), key)
	if err == redis.Nil {
		return "", ErrWalletLocked
	} else if err != nil {
		// Unknown error
		return "", err
	}

	return key, nil
}

// Set decrypted key to storage
func SetDecryptedKeyToStorage(wallet *ent.Wallet, key string, seed string) error {
	if wallet == nil {
		return ErrInvalidWallet
	} else if !wallet.Encrypted {
		return ErrWalletNotLocked
	}

	return database.GetRedisDB().Hset(wallet.ID.String(), key, seed)
}
