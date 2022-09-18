package controller

import (
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/render"
)

// Account handlers, reserved for the handlers that directly interact with the account_ actions

// Create a new account in sequence for given wallet
// ! TODO - we don't support setting a specific index like the node does, not sure the best way around that with how we handle our sequencing
func (hc *HttpController) HandleAccountCreate(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request, idx := hc.DecodeAccountCreateRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Create the account
	newAccount, err := hc.Wallet.AccountCreate(dbWallet, idx)
	if errors.Is(err, wallet.ErrWalletLocked) || errors.Is(err, wallet.ErrInvalidWallet) {
		ErrWalletLocked(w, r)
	} else if errors.Is(err, wallet.ErrAccountExists) {
		ErrBadRequest(w, r, "Account already exists")
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	resp := responses.AccountResponse{
		Account: newAccount.Address,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

// Handle bulk account create based on count param
func (hc *HttpController) HandleAccountsCreate(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request, count := hc.DecodeBaseRequestWithCount(rawRequest, w, r)
	if request == nil {
		return
	} else if count == 0 {
		// Not acceptable count
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Create the accounts
	newAccounts, err := hc.Wallet.AccountsCreate(dbWallet, count)
	if errors.Is(err, wallet.ErrWalletLocked) || errors.Is(err, wallet.ErrInvalidWallet) {
		ErrWalletLocked(w, r)
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Loop all accounts and get address as simple string array
	addresses := make([]string, len(newAccounts))
	for i, account := range newAccounts {
		addresses[i] = account.Address
	}

	resp := responses.AccountsResponse{
		Accounts: addresses,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

// Handle accounts_list
func (hc *HttpController) HandleAccountList(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request, count := hc.DecodeBaseRequestWithCount(rawRequest, w, r)
	if request == nil {
		return
	} else if count == 0 {
		// Default 1000
		count = 1000
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Accounts list
	_, accounts, err := hc.Wallet.AccountsList(dbWallet, count)
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	resp := responses.AccountsResponse{
		Accounts: accounts,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}
