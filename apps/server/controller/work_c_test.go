package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkGenerate(t *testing.T) {
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "work_generate",
		"hash":   "3F93C5CD2E314FA16702189041E68E68C07B27961BF37F0B7705145BEFBA3AA3",
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
	assert.Contains(t, respJson, "work")
	assert.Len(t, respJson["work"], 16)
	assert.Equal(t, "205452237a9b01f4", respJson["work"])

	// Test invalid hash
	reqBody = map[string]interface{}{
		"action": "work_generate",
		"hash":   "F93C5CD2E314FA16702189041E68E68C07B27961BF37F0B7705145BEFBA3AA3",
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

	// Make sure uuid is valid
	assert.Contains(t, respJson, "error")
	assert.Equal(t, "Invalid hash", respJson["error"])

}
