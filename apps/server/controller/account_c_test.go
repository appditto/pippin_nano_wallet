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

func TestAccountCreate(t *testing.T) {
	// Create wallet first
	// Request JSON
	seed := "3ab3de2721b57af86637e2c4c8994adf4f8eecfd62f59d013de0353500b9b823"
	reqBody := map[string]interface{}{
		"action": "wallet_create",
		"seed":   seed,
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

	// Create account
	reqBody = map[string]interface{}{
		"action": "account_create",
		"wallet": respJson["wallet"],
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

	// Make sure address is valid
	assert.Contains(t, respJson, "account")
	_, err = utils.AddressToPub(respJson["account"].(string), false)
	assert.Nil(t, err)
	assert.Equal(t, "nano_13coy8t4jzd516m5ydw8a7mdfguttcm6nkm4t69fwd1dzm87mgj5p8ijge8w", respJson["account"].(string))

	// Create an account at index 500 to make sure it comes back correctly
	// First derive expected
	pub, _, err := utils.KeypairFromSeed(seed, 500)
	assert.Nil(t, err)
	addr := utils.PubKeyToAddress(pub, false)

	reqBody = map[string]interface{}{
		"action": "account_create",
		"wallet": respJson["wallet"],
		"index":  500,
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

	// Make sure address is valid
	assert.Contains(t, respJson, "account")
	_, err = utils.AddressToPub(respJson["account"].(string), false)
	assert.Nil(t, err)
	assert.Equal(t, addr, respJson["account"].(string))
}

func TestAccountsCreate(t *testing.T) {
	// Create wallet first
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

	// Create account
	reqBody = map[string]interface{}{
		"action": "accounts_create",
		"wallet": respJson["wallet"],
		"count":  2,
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

	// Make sure address is valid
	assert.Contains(t, respJson, "accounts")
	// Convert to slice
	accounts := respJson["accounts"].([]interface{})
	assert.Equal(t, 2, len(accounts))
	for _, account := range accounts {
		_, err = utils.AddressToPub(account.(string), false)
		assert.Nil(t, err)
	}
}

func TestAccountList(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("f39a07504c76978f47e6630bb97e6fc169dd734d25ddcb323609a5699789b104"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Create some accounts
	MockController.Wallet.AccountsCreate(wallet, 3)

	// Test API
	reqBody := map[string]interface{}{
		"action": "account_list",
		"wallet": wallet.ID.String(),
		"count":  2,
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

	// Make sure address is valid
	assert.Contains(t, respJson, "accounts")
	// Convert to slice
	accounts := respJson["accounts"].([]interface{})
	assert.Equal(t, 2, len(accounts))
	for _, account := range accounts {
		_, err := utils.AddressToPub(account.(string), false)
		assert.Nil(t, err)
	}
}
