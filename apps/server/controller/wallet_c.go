package controller

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
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

// For adding adhoc keys to the wallet
func (hc *HttpController) HandleWalletAdd(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	// mapstructure decode
	var walletAddRequest requests.WalletAddRequest
	if err := mapstructure.Decode(request, &walletAddRequest); err != nil {
		klog.Errorf("Error unmarshalling wallet_add request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	// See if key is valid
	valid := utils.Validate64HexHash(walletAddRequest.Key)
	if !valid {
		ErrInvalidKey(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(walletAddRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Add key to wallet
	asHex, err := hex.DecodeString(walletAddRequest.Key)
	if err != nil {
		ErrInvalidKey(w, r)
		return
	}
	priv, err := ed25519.NewKeyFromSeed(asHex)
	if err != nil {
		ErrInvalidKey(w, r)
		return
	}
	acc, adhocAcc, err := hc.Wallet.AdhocAccountCreate(dbWallet, priv)
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil || (acc == nil && adhocAcc == nil) {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// The adhoc account create will return the normal sequential account too if it already exists
	var resp responses.WalletAddResponse
	if adhocAcc != nil {
		resp.Account = adhocAcc.Address
	} else {
		resp.Account = acc.Address
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

// Handle wallet locked, returns whether or not wallet is locked
func (hc *HttpController) HandleWalletLocked(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	// mapstructure decode
	var walletLockedRequest requests.WalletLockedRequest
	if err := mapstructure.Decode(request, &walletLockedRequest); err != nil {
		klog.Errorf("Error unmarshalling wallet_locked request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(walletLockedRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Check if wallet is locked
	_, err := wallet.GetDecryptedKeyFromStorage(dbWallet, "seed")
	var resp responses.WalletLockedResponse
	if err != nil {
		resp.Locked = "1"
	} else {
		resp.Locked = "0"
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

func (hc *HttpController) HandleWalletLock(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	// mapstructure decode
	var walletLockRequest requests.WalletLockRequest
	if err := mapstructure.Decode(request, &walletLockRequest); err != nil {
		klog.Errorf("Error unmarshalling wallet_lock request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(walletLockRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Lock wallet
	err := hc.Wallet.LockWallet(dbWallet)
	var resp = responses.WalletLockResponse{
		Locked: "1",
	}
	if errors.Is(err, wallet.ErrWalletNotLocked) {
		resp.Locked = "0"
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

func (hc *HttpController) HandleWalletDestroy(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var walletDestroyRequest requests.WalletDestroyRequest
	if err := mapstructure.Decode(rawRequest, &walletDestroyRequest); err != nil {
		klog.Errorf("Error unmarshalling wallet_destroy request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(walletDestroyRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Destroy list
	err := hc.Wallet.WalletDestroy(dbWallet)
	var resp = responses.WalletDestroyResponse{
		Destroyed: "1",
	}
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil {
		resp.Destroyed = "0"
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}
