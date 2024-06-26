package controller

import (
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
)

func (hc *HttpController) HandlePasswordChange(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var passwordChangeRequest requests.PasswordChangeRequest
	if err := mapstructure.Decode(rawRequest, &passwordChangeRequest); err != nil {
		log.Errorf("Error unmarshalling password_change request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(passwordChangeRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	changed, err := hc.Wallet.EncryptWallet(dbWallet, passwordChangeRequest.Password)
	var resp = responses.PasswordChangeResponse{Changed: "1"}
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil || !changed {
		resp.Changed = "0"
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

// Password enter is to unlock the wallet
func (hc *HttpController) HandlePasswordEnter(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var passwordEnterRequest requests.PasswordEnterRequest
	if err := mapstructure.Decode(rawRequest, &passwordEnterRequest); err != nil {
		log.Errorf("Error unmarshalling password_enter request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(passwordEnterRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Unlock the wallet
	unlocked, err := hc.Wallet.UnlockWallet(dbWallet, passwordEnterRequest.Password)
	var resp = responses.PasswordEnterResponse{Valid: "1"}
	if errors.Is(err, wallet.ErrWalletNotLocked) {
		ErrWalletNotLocked(w, r)
		return
	} else if err != nil || !unlocked {
		resp.Valid = "0"
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
