package controller

import (
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog/v2"
)

// Account handlers, reserved for the handlers that directly interact with the account_ actions

// Create a new account in sequence for given wallet
// ! TODO - we don't support setting a specific index like the node does, not sure the best way around that with how we handle our sequencing
func (hc *HttpController) HandleAccountCreate(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var accountCreateRequest requests.AccountCreateRequest
	if err := mapstructure.Decode(rawRequest, &accountCreateRequest); err != nil {
		klog.Errorf("Error unmarshalling account_create request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if accountCreateRequest.Wallet == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	dbWallet, err := hc.Wallet.GetWallet(accountCreateRequest.Wallet)
	if errors.Is(err, wallet.ErrWalletNotFound) {
		ErrWalletNotFound(w, r)
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Create the account
	newAccount, err := hc.Wallet.AccountCreate(dbWallet)
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	resp := responses.AccountCreateResponse{
		Account: newAccount.Address,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

// Handle bulk account create based on count param
func (hc *HttpController) HandleAccountsCreate(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var accountsCreateRequest requests.AccountsCreateRequest
	if err := mapstructure.Decode(rawRequest, &accountsCreateRequest); err != nil {
		klog.Errorf("Error unmarshalling account_create request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if accountsCreateRequest.Wallet == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	dbWallet, err := hc.Wallet.GetWallet(accountsCreateRequest.Wallet)
	if errors.Is(err, wallet.ErrWalletNotFound) {
		ErrWalletNotFound(w, r)
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	if accountsCreateRequest.Count < 1 {
		ErrUnableToParseJson(w, r)
		return
	}

	// Create the accounts
	newAccounts, err := hc.Wallet.AccountsCreate(dbWallet, accountsCreateRequest.Count)
	if errors.Is(err, wallet.ErrWalletLocked) {
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

	resp := responses.AccountsCreateResponse{
		Accounts: addresses,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}
