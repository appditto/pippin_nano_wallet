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

// Handle accounts_list
func (hc *HttpController) HandleReceiveRequest(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var receiveRequest requests.ReceiveRequest
	if err := mapstructure.Decode(rawRequest, &receiveRequest); err != nil {
		klog.Errorf("Error unmarshalling receive request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if receiveRequest.Wallet == "" || receiveRequest.Action == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(receiveRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Accounts list
	resp, err := hc.Wallet.CreateAndPublishReceiveBlock(dbWallet, receiveRequest.Account, receiveRequest.Block, receiveRequest.Work, receiveRequest.BpowKey)
	if err != nil {
		ErrBadRequest(w, r, err.Error())
		return
	}

	blockResponse := responses.BlockResponse{
		Block: resp,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &blockResponse)
}

// Handle accounts_list
func (hc *HttpController) HandleReceiveAllRequest(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Accounts list
	accounts, err := hc.Wallet.AccountsList(dbWallet, 0)
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	receivedCount := 0
	for _, account := range accounts {
		resp, err := hc.Wallet.ReceiveAllBlocks(dbWallet, account, request.BpowKey)
		if err != nil {
			ErrInternalServerError(w, r, err.Error())
			return
		}
		receivedCount += resp
	}

	resp := responses.ReceiveAllResponse{
		Received: receivedCount,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}
