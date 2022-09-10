package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	pw "github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/stretchr/testify/assert"
)

func TestPasswordChange(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("ffffff04c76978f47e6630bb97e6fc169dd734d25ddcb323609a5699789b104"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	// Wallet non-encrypted and empty password should return error
	reqBody := map[string]interface{}{
		"action":   "password_change",
		"wallet":   wallet.ID.String(),
		"password": "",
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

	// Make sure not changed
	assert.Equal(t, "0", respJson["changed"].(string))

	// Make sure wallet not locked
	seed, _ := pw.GetDecryptedKeyFromStorage(wallet, "seed")
	assert.Equal(t, newSeed, seed)

	// Wallet non-encrypted and empty password should return error
	reqBody = map[string]interface{}{
		"action":   "password_change",
		"wallet":   wallet.ID.String(),
		"password": "mypassword",
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

	// Make sure changed
	assert.Equal(t, "1", respJson["changed"].(string))

	// // Make sure wallet is locked now
	nWallet, _ := MockController.Wallet.GetWallet(wallet.ID.String())
	_, err := pw.GetDecryptedKeyFromStorage(nWallet, "seed")
	assert.ErrorIs(t, err, pw.ErrWalletLocked)
}

func TestPasswordEnter(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("eeeeee04c76978f47e6630bb97e6fc169dd734d25ddcb323609a5699789b104"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Lock wallet
	MockController.Wallet.EncryptWallet(wallet, "mypassword")
	// Request JSON
	// Wallet non-encrypted and empty password should return error
	reqBody := map[string]interface{}{
		"action":   "password_enter",
		"wallet":   wallet.ID.String(),
		"password": "badpassword",
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

	// Make sure not valid
	assert.Equal(t, "0", respJson["valid"].(string))

	// Make sure wallet not unlocked
	_, err := pw.GetDecryptedKeyFromStorage(wallet, "seed")
	assert.NotNil(t, err)

	// Wallet non-encrypted and empty password should return error
	reqBody = map[string]interface{}{
		"action":   "password_enter",
		"wallet":   wallet.ID.String(),
		"password": "mypassword",
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

	// Make sure changed
	assert.Equal(t, "1", respJson["valid"].(string))

	// Make sure wallet is unlocked
	nWallet, _ := MockController.Wallet.GetWallet(wallet.ID.String())
	seed, err := pw.GetDecryptedKeyFromStorage(nWallet, "seed")
	assert.Nil(t, err)
	assert.Equal(t, newSeed, seed)
}
