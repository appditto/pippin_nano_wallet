package controller

import (
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog/v2"
)

// Wallet handlers, reserved for the handlers that directly interact with the wallet
// e.g. wallet_create, wallet_add, wallet_destroy

func (hc *HttpController) HandleWalletCreate(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	// mapstructure decode
	var walletCreateRequest requests.WalletCreateRequest
	if err := mapstructure.Decode(request, &walletCreateRequest); err != nil {
		klog.Errorf("Error unmarshalling wallet_create request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	var seed string
	var err error
	if walletCreateRequest.Seed != nil {
		seed = *walletCreateRequest.Seed
	} else {
		seed, err = utils.GenerateSeed(nil)
		if err != nil {
			ErrInternalServerError(w, r, err.Error())
			return
		}
	}

	newWallet, err := hc.Wallet.WalletCreate(seed)
	if errors.Is(err, wallet.ErrInvalidSeed) {
		ErrInvalidSeed(w, r)
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// We have a wallet
	// Return the wallet id
	walletCreateResponse := responses.WalletCreateResponse{
		Wallet: newWallet.ID.String(),
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, &walletCreateResponse)
}
