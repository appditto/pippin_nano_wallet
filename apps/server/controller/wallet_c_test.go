package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

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
	newSeed, _ := utils.GenerateSeed(strings.NewReader("E11A48D701EA1F8A66A4EB587CDC8808D726FE75B325DF204F62CA2B43F9ADA1"))
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
