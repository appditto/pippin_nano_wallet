package controller

import (
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
)

// Wallet handlers, reserved for the handlers that directly interact with the wallet
// e.g. wallet_create, wallet_add, wallet_destroy

func (hc *HttpController) HandleWalletCreate(request *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	// mapstructure decode
	var walletCreateRequest requests.WalletCreateRequest
	if err := mapstructure.Decode(request, &walletCreateRequest); err != nil {
		log.Errorf("Error unmarshalling wallet_create request %s", err)
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
		log.Errorf("Error unmarshalling wallet_add request %s", err)
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
	acc, err := hc.Wallet.AdhocAccountCreate(dbWallet, priv)
	if errors.Is(err, wallet.ErrWalletLocked) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil || acc == nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// The adhoc account create will return the normal sequential account too if it already exists
	var resp responses.AccountResponse
	resp.Account = acc.Address

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &resp)
}

// Handle wallet locked, returns whether or not wallet is locked
func (hc *HttpController) HandleWalletLocked(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
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

func (hc *HttpController) HandleWalletLock(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	// mapstructure decode
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Lock wallet
	err := hc.Wallet.LockWallet(dbWallet)
	var resp = responses.WalletLockedResponse{
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
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
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

func (hc *HttpController) HandleWalletBalances(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Get accounts on wallet
	_, accounts, err := hc.Wallet.AccountsList(dbWallet, math.MaxInt)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Get RPC balances
	resp, err := hc.RpcClient.MakeAccountsBalancesRequest(accounts)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Return balances
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (hc *HttpController) HandleWalletFrontiers(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Get accounts on wallet
	_, accounts, err := hc.Wallet.AccountsList(dbWallet, math.MaxInt)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Get RPC balances
	resp, err := hc.RpcClient.MakeAccountsFrontiersRequest(accounts)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Return balances
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (hc *HttpController) HandleWalletPending(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Get accounts on wallet
	_, accounts, err := hc.Wallet.AccountsList(dbWallet, math.MaxInt)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Get RPC balances
	resp, err := hc.RpcClient.MakeAccountsPendingRequest(accounts)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Return balances
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (hc *HttpController) HandleWalletInfo(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Get accounts on wallet
	_, accounts, err := hc.Wallet.AccountsList(dbWallet, math.MaxInt)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Get RPC balances
	resp, err := hc.RpcClient.MakeAccountsBalancesRequest(accounts)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	balance := big.NewInt(0)
	pendingBalance := big.NewInt(0)

	// for k, v in balance_json['balances'].items():
	// balance += int(v['balance'])
	// pending_bal += int(v['pending'])
	for _, v := range *resp.Balances {
		// Balance to bigint
		balanceBigInt, ok := new(big.Int).SetString(v.Balance, 10)
		if !ok {
			ErrInternalServerError(w, r, "Could not parse balance")
			return
		}
		balance = balance.Add(balance, balanceBigInt)
		pendingBigInt, ok := new(big.Int).SetString(v.Pending, 10)
		if !ok {
			ErrInternalServerError(w, r, "Could not parse balance")
			return
		}
		pendingBalance = pendingBalance.Add(pendingBalance, pendingBigInt)
	}

	// Retrieve wallet info from database
	walletInfo, err := hc.Wallet.WalletInfo(dbWallet)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	// Return balances
	render.Status(r, http.StatusOK)
	render.JSON(w, r, responses.WalletInfoResponse{
		Balance:            balance.String(),
		Pending:            pendingBalance.String(),
		Receivable:         pendingBalance.String(),
		AccountsCount:      walletInfo.AccountsCount,
		AdhocCount:         walletInfo.AdhocCount,
		DeterministicCount: walletInfo.DeterministicCount,
		DeterministicIndex: walletInfo.DeterministicIndex,
	})
}

func (hc *HttpController) HandleWalletContains(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var request requests.WalletContainsRequest
	if err := mapstructure.Decode(rawRequest, &request); err != nil {
		log.Errorf("Error unmarshalling request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if request.Wallet == "" || request.Action == "" || request.Account == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Validate account
	_, err := utils.AddressToPub(request.Account, hc.Wallet.Banano)
	if err != nil {
		ErrInvalidAccount(w, r)
		return
	}

	// See if account exists
	exists, err := hc.Wallet.AccountExists(dbWallet, request.Account)
	if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	var resp responses.WalletContainsResponse
	if exists {
		resp = responses.WalletContainsResponse{
			Exists: "1",
		}
	} else {
		resp = responses.WalletContainsResponse{
			Exists: "0",
		}
	}

	// Return balances
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (hc *HttpController) HandleWalletRepresentativeSetRequest(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var changeRequest requests.WalletRepresentativeSetRequest
	if err := mapstructure.Decode(rawRequest, &changeRequest); err != nil {
		log.Errorf("Error unmarshalling receive request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if changeRequest.Wallet == "" || changeRequest.Action == "" || changeRequest.Representative == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(changeRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Validate account
	_, err := utils.AddressToPub(changeRequest.Representative, hc.Wallet.Config.Wallet.Banano)
	if err != nil {
		ErrBadRequest(w, r, "Invalid representative account")
		return
	}

	// Parse as bool
	updateExisting := false
	if changeRequest.UpdateExistingAccounts != nil {
		updateExisting, err = utils.ToBool(*changeRequest.UpdateExistingAccounts)
	}

	err = hc.Wallet.WalletRepresentativeSet(dbWallet, changeRequest.Representative, updateExisting, changeRequest.BpowKey)
	setResponse := responses.SetResponse{
		Set: "1",
	}
	if err != nil {
		setResponse.Set = "0"
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, setResponse)
}

func (hc *HttpController) HandleWalletRepresentativeRequest(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	request := hc.DecodeBaseRequest(rawRequest, w, r)
	if request == nil {
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(request.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// See if wallet locked
	_, err := wallet.GetDecryptedKeyFromStorage(dbWallet, "seed")
	if errors.Is(err, wallet.ErrWalletLocked) || errors.Is(err, wallet.ErrInvalidWallet) {
		ErrWalletLocked(w, r)
		return
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	var representative string
	if dbWallet.Representative == nil {
		// Just get a random one since we don't have one set
		representative, err = hc.Wallet.Config.GetRandomRep()
		if err != nil {
			ErrInternalServerError(w, r, err.Error())
			return
		}
	} else {
		representative = *dbWallet.Representative
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &responses.RepresentativeResponse{
		Representative: representative,
	})
}

func (hc *HttpController) HandleWalletChangeSeedRequest(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var changeRequest requests.WalletChangeSeedRequest
	if err := mapstructure.Decode(rawRequest, &changeRequest); err != nil {
		log.Errorf("Error unmarshalling change seed request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if changeRequest.Wallet == "" || changeRequest.Action == "" || changeRequest.Seed == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	// See if wallet exists
	dbWallet := hc.WalletExists(changeRequest.Wallet, w, r)
	if dbWallet == nil {
		return
	}

	// Change the seed
	newest, err := hc.Wallet.WalletChangeSeed(dbWallet, changeRequest.Seed)
	if errors.Is(err, wallet.ErrWalletLocked) || errors.Is(err, wallet.ErrInvalidWallet) {
		ErrWalletLocked(w, r)
		return
	} else if errors.Is(err, wallet.ErrInvalidSeed) {
		ErrBadRequest(w, r, err.Error())
	} else if err != nil {
		ErrInternalServerError(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &responses.WalletChangeSeedResponse{
		Success:             "",
		LastRestoredAccount: newest.Address,
		RestoredCount:       *newest.AccountIndex + 1,
	})
}
