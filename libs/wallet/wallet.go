package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	config "github.com/appditto/pippin_nano_wallet/libs/config/models"
	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent/account"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	nanorpc "github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/wallet/models"
	"github.com/google/uuid"
)

type NanoWallet struct {
	DB         *ent.Client
	Ctx        context.Context
	RpcClient  *nanorpc.RPCClient
	WorkClient *pow.PippinPow
	Config     *config.PippinConfig
	Banano     bool
}

var ErrInvalidSeed = errors.New("invalid seed")
var ErrInvalidWallet = errors.New("invalid wallet")
var ErrInvalidAccount = errors.New("invalid account")
var ErrInvalidPrivKey = errors.New("invalid private key")
var ErrInvalidAccountCount = errors.New("invalid count")
var ErrWalletNotFound = errors.New("wallet not found")

// Retrieves wallet
func (w *NanoWallet) GetWallet(walletID string) (*ent.Wallet, error) {
	parsedUuid, err := uuid.Parse(walletID)
	if err != nil {
		return nil, ErrWalletNotFound
	}
	wallet, err := w.DB.Wallet.Get(w.Ctx, parsedUuid)
	if ent.IsNotFound(err) {
		return nil, ErrWalletNotFound
	} else if err != nil {
		return nil, err
	}

	return wallet, nil
}

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

func (w *NanoWallet) WalletDestroy(wallet *ent.Wallet) error {
	if wallet == nil {
		return ErrInvalidWallet
	}

	// Determine if wallet is locked or not
	_, err := GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return err
	}

	// Delete wallet
	err = w.DB.Wallet.DeleteOne(wallet).Exec(w.Ctx)
	if err != nil {
		return err
	}

	return nil
}

func (w *NanoWallet) WalletInfo(wallet *ent.Wallet) (*models.WalletInfo, error) {
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
	_, err = GetDecryptedKeyFromStorage(wallet, "seed")
	if err != nil {
		return nil, err
	}

	curAccount, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.AccountIndexNotNil()).Order(ent.Desc(account.FieldAccountIndex)).First(w.Ctx)

	if err != nil {
		return nil, err
	}

	currentIndex := *curAccount.AccountIndex

	// Get all accounts on wallet
	accounts, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.AccountIndexNotNil()).Count(w.Ctx)
	if err != nil {
		return nil, err
	}

	// Get all adhoc accounts on wallet
	adhocAccounts, err := w.DB.Account.Query().Where(account.WalletID(wallet.ID), account.PrivateKeyNotNil()).Count(w.Ctx)
	if err != nil {
		return nil, err
	}

	return &models.WalletInfo{
		AccountsCount:      accounts + adhocAccounts,
		AdhocCount:         adhocAccounts,
		DeterministicCount: accounts,
		DeterministicIndex: currentIndex,
	}, nil
}

func (w *NanoWallet) WalletRepresentativeSet(wallet *ent.Wallet, representative string, changeExisting bool, bpowKey *string) error {
	if wallet == nil {
		return ErrInvalidWallet
	}

	// Update wallet with representative
	wallet, err := wallet.Update().SetRepresentative(representative).Save(w.Ctx)
	if err != nil {
		return err
	}

	if !changeExisting {
		return nil
	}

	// Create and publish change blocks for every account on the wallet.
	_, addresses, err := w.AccountsList(wallet, 0)
	if err != nil {
		return err
	}
	for _, address := range addresses {
		_, err := w.CreateAndPublishChangeBlock(wallet, address, representative, nil, bpowKey, true)
		if err != nil && !errors.Is(err, ErrSameRepresentative) && !errors.Is(err, nanorpc.ErrAccountNotFound) {
			return err
		}
	}
	return nil
}
