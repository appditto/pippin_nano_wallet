package controller

import (
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
)

// Some common things multiple handlers use

// Get wallet if it exists, set response
func (hc *HttpController) WalletExists(walletId string, w http.ResponseWriter, r *http.Request) *ent.Wallet {
	// See if wallet exists
	dbWallet, err := hc.Wallet.GetWallet(walletId)
	if errors.Is(err, wallet.ErrWalletNotFound) || errors.Is(err, wallet.ErrInvalidWallet) {
		ErrWalletNotFound(w, r)
		return nil
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return nil
	}

	return dbWallet
}
