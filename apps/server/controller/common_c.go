package controller

import (
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/mitchellh/mapstructure"
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

// Common map decoding for most requests
func (hc *HttpController) DecodeBaseRequest(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) *requests.BaseRequest {
	var baseRequest requests.BaseRequest
	if err := mapstructure.Decode(request, &baseRequest); err != nil {
		log.Errorf("Error unmarshalling request %s", err)
		ErrUnableToParseJson(w, r)
		return nil
	} else if baseRequest.Wallet == "" || baseRequest.Action == "" {
		ErrUnableToParseJson(w, r)
		return nil
	}
	return &baseRequest
}

// Common map decoding for requests with count added
func (hc *HttpController) DecodeBaseRequestWithCount(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) (*requests.BaseRequestWithCount, int) {
	var baseRequest requests.BaseRequestWithCount
	if err := mapstructure.Decode(request, &baseRequest); err != nil {
		log.Errorf("Error unmarshalling request with count %s", err)
		ErrUnableToParseJson(w, r)
		return nil, 0
	} else if baseRequest.Wallet == "" || baseRequest.Action == "" {
		ErrUnableToParseJson(w, r)
		return nil, 0
	}

	var count int
	var err error
	if baseRequest.Count != nil {
		count, err = utils.ToInt(*baseRequest.Count)
		if err != nil || count < 1 {
			ErrUnableToParseJson(w, r)
			return nil, 0
		}
		if count < 1 {
			count = 1
		}
	}

	return &baseRequest, count
}

// ! TODO - can we reduce duplication with generics or something ?
func (hc *HttpController) DecodeAccountCreateRequest(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) (*requests.AccountCreateRequest, *int) {
	var accountCreateRequest requests.AccountCreateRequest
	if err := mapstructure.Decode(request, &accountCreateRequest); err != nil {
		log.Errorf("Error unmarshalling request with count %s", err)
		ErrUnableToParseJson(w, r)
		return nil, nil
	} else if accountCreateRequest.Wallet == "" || accountCreateRequest.Action == "" {
		ErrUnableToParseJson(w, r)
		return nil, nil
	}

	var idx *int
	if accountCreateRequest.Index != nil {
		index, err := utils.ToInt(*accountCreateRequest.Index)
		if err != nil || index < 0 {
			ErrUnableToParseJson(w, r)
			return nil, nil
		}
		idx = &index
	}

	return &accountCreateRequest, idx
}
