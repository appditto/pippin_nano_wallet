package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletCreateWithoutSeed(t *testing.T) {
	os.Setenv("MOCK_REDIS", "true")
	defer os.Unsetenv("MOCK_REDIS")
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_create",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	MockController.HandleWalletCreate(&reqBody, w, req)
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
	os.Setenv("MOCK_REDIS", "true")
	defer os.Unsetenv("MOCK_REDIS")
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
	MockController.HandleWalletCreate(&reqBody, w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Invalid seed", respJson["error"])
}

func TestWalletCreateWithValidSeed(t *testing.T) {
	os.Setenv("MOCK_REDIS", "true")
	defer os.Unsetenv("MOCK_REDIS")
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
	MockController.HandleWalletCreate(&reqBody, w, req)
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
