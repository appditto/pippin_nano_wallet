package controller

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrUnableToParseJson(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrUnableToParseJson(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Unable to parse json", respJson["error"])
}

func TestErrInvalidSeed(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrInvalidSeed(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Invalid seed", respJson["error"])
}

func TestErrWalletNotFound(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrWalletNotFound(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "wallet not found", respJson["error"])
}

func TestErrWalletLocked(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrWalletLocked(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "wallet locked", respJson["error"])
}

func TestErrWalletNotLocked(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrWalletNotLocked(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "wallet not locked", respJson["error"])
}

func TestErrInvalidKey(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrInvalidKey(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Invalid key", respJson["error"])
}

func TestErrNoWalletPassword(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrNoWalletPassword(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "password not set", respJson["error"])
}

func TestErrInvalidHash(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrInvalidHash(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Invalid hash", respJson["error"])
}

func TestErrInternalServerError(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrInternalServerError(w, req, "server error")
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 500, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "server error", respJson["error"])
}

func TestErrBadRequest(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrBadRequest(w, req, "server error")
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "server error", respJson["error"])
}

func TestErrWorkFailed(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrWorkFailed(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 500, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Failed to generate work", respJson["error"])
}

func TestErrInvalidAccount(t *testing.T) {
	w := httptest.NewRecorder()
	// Build request
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Content-Type", "application/json")
	ErrInvalidAccount(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	assert.Equal(t, 400, resp.StatusCode)

	var respJson map[string]interface{}
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "Invalid account", respJson["error"])
}
