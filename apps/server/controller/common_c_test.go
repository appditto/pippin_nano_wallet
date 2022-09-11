package controller

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/stretchr/testify/assert"
)

func TestWalletExists(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("9d9e1ede8170a7ef7fee2e28990dbc78c150b705ede136c4ab39dec349c38f42"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	we := MockController.WalletExists(wallet.ID.String(), w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, we)

	// Check one that doesn't exist
	w = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/", nil)
	we = MockController.WalletExists("123", w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)
	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "error")
	assert.Equal(t, "wallet not found", respJson["error"].(string))
}

func TestDecodeBaseRequest(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	reqBody := map[string]interface{}{
		"action": "password_change",
		"wallet": "1234",
	}
	be := MockController.DecodeBaseRequest(&reqBody, w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, be)
	assert.Equal(t, "password_change", be.Action)
	assert.Equal(t, "1234", be.Wallet)

	// With error
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("GET", "/", nil)
	reqBody = map[string]interface{}{
		"wallet": "1234",
	}
	be = MockController.DecodeBaseRequest(&reqBody, w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)
	assert.Nil(t, be)
	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "error")
	assert.Equal(t, "Unable to parse json", respJson["error"].(string))
}

func TestDecodeBaseRequestWithCount(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	reqBody := map[string]interface{}{
		"action": "password_change",
		"wallet": "1234",
		"count":  "1",
	}
	be, count := MockController.DecodeBaseRequestWithCount(&reqBody, w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, be)
	assert.Equal(t, "password_change", be.Action)
	assert.Equal(t, "1234", be.Wallet)
	assert.Equal(t, 1, count)

	// Empty count should return 0
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("GET", "/", nil)
	reqBody = map[string]interface{}{
		"action": "password_change",
		"wallet": "1234",
	}
	be, count = MockController.DecodeBaseRequestWithCount(&reqBody, w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
	assert.NotNil(t, be)
	assert.Equal(t, "password_change", be.Action)
	assert.Equal(t, "1234", be.Wallet)
	assert.Equal(t, 0, count)

	// With error
	w = httptest.NewRecorder()
	// Build request
	req = httptest.NewRequest("GET", "/", nil)
	reqBody = map[string]interface{}{
		"wallet": "1234",
	}
	be, count = MockController.DecodeBaseRequestWithCount(&reqBody, w, req)
	resp = w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)
	assert.Nil(t, be)
	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Contains(t, respJson, "error")
	assert.Equal(t, "Unable to parse json", respJson["error"].(string))
}
