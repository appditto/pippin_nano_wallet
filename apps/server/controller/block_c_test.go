package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/mocks"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReceive(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var pr requests.BaseRequest
			json.NewDecoder(req.Body).Decode(&pr)
			if pr.Action == "block_info" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.BlockInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			} else if pr.Action == "account_info" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.AccountInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			} else if pr.Action == "process" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.ProcessResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)
	newSeed, _ := utils.GenerateSeed(strings.NewReader("D8622D6D3754151F881D19CDE6E243C58AA05426ABB2BD5430455E653200ACC6"))
	wallet, err := MockController.Wallet.WalletCreate(newSeed)
	assert.Nil(t, err)
	acc, err := MockController.Wallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	// Request JSON
	reqBody := map[string]interface{}{
		"action":  "receive",
		"wallet":  wallet.ID.String(),
		"account": acc.Address,
		"block":   "95D72CE5ECA6ABFDE45F77BD75F1C888223BCCA2D5178DF2A1D89533005C69DC",
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

	var respJson responses.BlockResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "E2FB233EF4554077A7BF1AA85851D5BF0B36965D2B0FB504B2BC778AB89917D3", respJson.Block)

	// errors

	// Request JSON
	reqBody = map[string]interface{}{
		"action":  "receive",
		"wallet":  wallet.ID.String(),
		"account": "ban_1234",
		"block":   "95D72CE5ECA6ABFDE45F77BD75F1C888223BCCA2D5178DF2A1D89533005C69DC",
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

	var rawResp map[string]interface{}
	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &rawResp)

	assert.Equal(t, "Invalid account", rawResp["error"])
}

func TestSend(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var pr requests.BaseRequest
			json.NewDecoder(req.Body).Decode(&pr)
			if pr.Action == "block_info" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.BlockInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			} else if pr.Action == "account_info" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.AccountInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			} else if pr.Action == "process" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.ProcessResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)
	newSeed, _ := utils.GenerateSeed(strings.NewReader("95D72CE5ECA6ABFDE45F77BD75F1C888223BCCA2D5178DF2A1D89533005C69DC"))
	wallet, err := MockController.Wallet.WalletCreate(newSeed)
	assert.Nil(t, err)
	acc, err := MockController.Wallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	// Request JSON
	reqBody := map[string]interface{}{
		"action":      "send",
		"wallet":      wallet.ID.String(),
		"source":      acc.Address,
		"destination": acc.Address,
		"amount":      "1000000000000000000000000000000",
		"work":        "0000000000000000",
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

	var respJson responses.BlockResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "E2FB233EF4554077A7BF1AA85851D5BF0B36965D2B0FB504B2BC778AB89917D3", respJson.Block)

	// errors

	// Request JSON
	reqBody = map[string]interface{}{
		"action":      "send",
		"wallet":      wallet.ID.String(),
		"source":      acc.Address,
		"destination": "ban_1234",
		"amount":      "1000000000000000000000000000000",
		"work":        "0000000000000000",
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

	var rawResp map[string]interface{}
	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &rawResp)

	assert.Equal(t, "Invalid destination account", rawResp["error"])

	// Request JSON
	reqBody = map[string]interface{}{
		"action":      "send",
		"wallet":      wallet.ID.String(),
		"source":      "ban_1234",
		"destination": acc.Address,
		"amount":      "1000000000000000000000000000000",
		"block":       "CE34A533C3D6BFB820FD28F132FFF3FDD48D5DA979F6063DD16440C5EA075465",
		"work":        "0000000000000000",
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
	json.Unmarshal(respBody, &rawResp)

	assert.Equal(t, "Invalid source account", rawResp["error"])
}

func TestAccountRepresentativeSet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var pr requests.BaseRequest
			json.NewDecoder(req.Body).Decode(&pr)
			if pr.Action == "block_info" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.BlockInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			} else if pr.Action == "account_info" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.AccountInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			} else if pr.Action == "process" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.ProcessResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)
	newSeed, _ := utils.GenerateSeed(strings.NewReader("5DD64B08829F71953CE37435E3CAB4113691FDDCBD949004188C10CD15535DD6"))
	wallet, err := MockController.Wallet.WalletCreate(newSeed)
	assert.Nil(t, err)
	acc, err := MockController.Wallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)
	// Request JSON
	reqBody := map[string]interface{}{
		"action":         "account_representative_set",
		"wallet":         wallet.ID.String(),
		"account":        acc.Address,
		"representative": acc.Address,
		"work":           "0000000000000000",
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

	var respJson responses.BlockResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "E2FB233EF4554077A7BF1AA85851D5BF0B36965D2B0FB504B2BC778AB89917D3", respJson.Block)

	// errors

	// Request JSON
	reqBody = map[string]interface{}{
		"action":         "account_representative_set",
		"wallet":         wallet.ID.String(),
		"account":        "ban_1234",
		"representative": acc.Address,
		"work":           "0000000000000000",
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

	var rawResp map[string]interface{}
	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &rawResp)

	assert.Equal(t, "Invalid account", rawResp["error"])

	// Request JSON
	reqBody = map[string]interface{}{
		"action":         "account_representative_set",
		"wallet":         wallet.ID.String(),
		"account":        acc.Address,
		"representative": "ban_1234",
		"work":           "0000000000000000",
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
	json.Unmarshal(respBody, &rawResp)

	assert.Equal(t, "Invalid representative account", rawResp["error"])
}
