package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
	"k8s.io/klog/v2"
)

type HttpController struct {
	Wallet *wallet.NanoWallet
}

// This API is intended to replace the nano node wallet RPCs
// https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs
// It will:
// 1) Determine if the request is a supported wallet RPC, if so process it
// 2) If not, pass it to the nano node RPC
// The error messages and behavior are also intended to replace what the nano node returns
// The node isn't exactly great at returning errors, and the error messages are not very helpful
// But as we want to be a drop-in replacement we mimic the behavior
func (hc *HttpController) HandleAction(w http.ResponseWriter, r *http.Request) {
	var baseRequest map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&baseRequest); err != nil {
		klog.Errorf("Error unmarshalling http base request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	if _, ok := baseRequest["action"]; !ok {
		ErrUnableToParseJson(w, r)
		return
	}

	action := strings.ToLower(fmt.Sprintf("%v", baseRequest["action"]))

	switch action {
	case "wallet_create":
		hc.handleWalletCreate(&baseRequest, w, r)
		return
	}
}

func (hc *HttpController) handleWalletCreate(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
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
