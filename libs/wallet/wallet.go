package wallet

import (
	"context"
	"errors"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
)

type NanoWallet struct {
	DB  *ent.Client
	Ctx context.Context
}

func (w *NanoWallet) WalletCreate(seed string) (*ent.Wallet, error) {
	if !utils.Validate64HexHash(seed) {
		return nil, errors.New("invalid seed")
	}
	wallet, err := w.DB.Wallet.Create().SetSeed(seed).Save(w.Ctx)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}
