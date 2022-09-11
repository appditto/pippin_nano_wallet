package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/rpc/mocks"
	rpcresp "github.com/appditto/pippin_nano_wallet/libs/rpc/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletCreateWithoutSeed(t *testing.T) {
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_create",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	// Make sure uuid is valid
	assert.Contains(t, respJson, "wallet")
	_, err := uuid.Parse(respJson["wallet"].(string))
	assert.Nil(t, err)
}

func TestWalletCreateWithInvalidSeed(t *testing.T) {
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_create",
		"seed":   "1234",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Invalid seed", respJson["error"])
}

func TestWalletCreateWithValidSeed(t *testing.T) {
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_create",
		"seed":   "1A2E95A2DCF03143297572EAEC496F6913D5001D2F28A728B35CB274294D5A14",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	// Make sure uuid is valid
	assert.Contains(t, respJson, "wallet")
	_, err := uuid.Parse(respJson["wallet"].(string))
	assert.Nil(t, err)
}

func TestWalletAdd(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("E11A48D701EA1F8A66A4EB587CDC8808D726FE75B325DF204F62CA2B43F9ADA1"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_add",
		"wallet": wallet.ID.String(),
		"key":    "67EDBC8F904091738DF33B4B6917261DB91DD9002D3985A7BA090345264A46C6",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	// Make sure account is valid
	assert.Contains(t, respJson, "account")
	assert.Equal(t, "nano_1p95xji1g5gou8auj8h6qcuezpdpcyoqmawao6mpwj4p44939oouoturkggc", respJson["account"].(string))

	// Test same thing with an invalid key
	// Request JSON
	reqBody = map[string]interface{}{
		"action": "wallet_add",
		"wallet": wallet.ID.String(),
		"key":    "abcd",
	}
	body, _ = json.Marshal(reqBody)
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	// Make sure account is valid
	assert.Contains(t, respJson, "error")
	assert.Equal(t, "Invalid key", respJson["error"].(string))
}

func TestWalletLocked(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("A11A48D701EA1F8A66A4EB587CDC8808D726FE75B325DF204F62CA2B43F9ADA1"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_locked",
		"wallet": wallet.ID.String(),
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "locked")
	assert.Equal(t, "0", respJson["locked"].(string))

	// Actually lock the wallet
	MockController.Wallet.EncryptWallet(wallet, "password")

	reqBody = map[string]interface{}{
		"action": "wallet_locked",
		"wallet": wallet.ID.String(),
	}
	body, _ = json.Marshal(reqBody)
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "locked")
	assert.Equal(t, "1", respJson["locked"].(string))
}

func TestWalletLock(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("ca439f7f9e6a3e2e0291b71391d2a097e6ba9912e5402588ddebf339fe46b270"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_lock",
		"wallet": wallet.ID.String(),
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "locked")
	assert.Equal(t, "0", respJson["locked"].(string))

	// Actually lock the wallet
	MockController.Wallet.EncryptWallet(wallet, "password")

	reqBody = map[string]interface{}{
		"action": "wallet_lock",
		"wallet": wallet.ID.String(),
	}
	body, _ = json.Marshal(reqBody)
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "locked")
	assert.Equal(t, "1", respJson["locked"].(string))
}

func TestWalletDestroy(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("43cededf4d2bacaa096bfe0251519d2adedc31aa1a417073c5a23f30e74b3ed7"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// lock walet
	MockController.Wallet.EncryptWallet(wallet, "password")
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_destroy",
		"wallet": wallet.ID.String(),
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "error")
	assert.Equal(t, "wallet locked", respJson["error"].(string))

	// unlock wallet
	MockController.Wallet.UnlockWallet(wallet, "password")
	reqBody = map[string]interface{}{
		"action": "wallet_destroy",
		"wallet": wallet.ID.String(),
	}
	body, _ = json.Marshal(reqBody)
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "destroyed")
	assert.Equal(t, "1", respJson["destroyed"].(string))

	// check if wallet is destroyed
	_, err := MockController.Wallet.GetWallet(wallet.ID.String())
	assert.NotNil(t, err)
}

func TestWalletBalances(t *testing.T) {
	// mock rpc response
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: mocks.AccountsBalancesResponse,
		}, nil
	}
	newSeed, _ := utils.GenerateSeed(strings.NewReader("b0ec8a9edddc58abba90656853a27cb9968f0ba48b432c04eb0f2e9143d1e34c"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_balances",
		"wallet": wallet.ID.String(),
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.Gateway(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var respJson rpcresp.AccountsBalancesResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	balances := *respJson.Balances
	assert.Len(t, balances, 1)
	assert.Equal(t, "11999999999999999918751838129509869131", balances["nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5"].Balance)
	assert.Equal(t, "0", balances["nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5"].Pending)
	assert.Equal(t, "0", balances["nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5"].Receivable)
}
