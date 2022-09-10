package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAccountCreate(t *testing.T) {
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
	_, err = utils.AddressToPub(respJson["account"].(string))
	assert.Nil(t, err)
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
		_, err = utils.AddressToPub(account.(string))
		assert.Nil(t, err)
	}
}
