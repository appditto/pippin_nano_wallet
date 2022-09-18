package controller

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/mocks"
	rpcresp "github.com/appditto/pippin_nano_wallet/libs/rpc/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.AccountBalancesResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)
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
func TestWalletFrontiers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.AccountsFrontiersResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)
	newSeed, _ := utils.GenerateSeed(strings.NewReader("597bce5b7950e33ad0c7755bf374b2a4c1a59a43e89b53cca1e4871bf8683c5e"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_frontiers",
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

	var respJson rpcresp.AccountsFrontiersResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	frontiers := *respJson.Frontiers
	assert.Len(t, frontiers, 2)
	assert.Equal(t, "791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A", frontiers["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"])
	assert.Equal(t, "6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F", frontiers["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"])
}

func TestWalletPending(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.AccountsPendingResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)
	newSeed, _ := utils.GenerateSeed(strings.NewReader("c53d6d42e94a41f55d8c7df53ebfceb7dffa468be271ab84e3d1a8221d385e0d"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_pending",
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

	var respJson rpcresp.AccountsPendingResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	blocks := *respJson.Blocks
	assert.Len(t, blocks, 2)
	assert.Equal(t, "4C1FEEF0BEA7F50BE35489A1233FE002B212DEA554B55B1B470D78BD8F210C74", blocks["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"][0])
	assert.Equal(t, "142A538F36833D1CC78B94E11C766F75818F8B940771335C6C1B8AB880C5BB1D", blocks["nano_1111111111111111111111111111111111111111111111111117353trpda"][0])
}

func TestWalletInfo(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.AccountBalancesResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)
	newSeed, _ := utils.GenerateSeed(strings.NewReader("34c77ad1c1fbf4fd026c7c8c80069c0e84c636fcd9f0e4f14a69dacf1f9d67d0"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Create some accounts
	MockController.Wallet.AccountsCreate(wallet, 2)
	// Request JSON
	reqBody := map[string]interface{}{
		"action": "wallet_info",
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

	var respJson responses.WalletInfoResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, 3, respJson.AccountsCount)
	assert.Equal(t, 0, respJson.AdhocCount)
	assert.Equal(t, 3, respJson.DeterministicCount)
	assert.Equal(t, 2, respJson.DeterministicIndex)
	assert.Equal(t, "11999999999999999918751838129509869131", respJson.Balance)
	assert.Equal(t, "0", respJson.Pending)
	assert.Equal(t, "0", respJson.Receivable)
}

func TestWalletContains(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("0dc54dc735bb2daad287a1bebc330f78a77af6097b9b45ac9a6144d4e57174ee"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Create some accounts
	accs, _ := MockController.Wallet.AccountsCreate(wallet, 2)
	// Request JSON
	reqBody := map[string]interface{}{
		"action":  "wallet_contains",
		"wallet":  wallet.ID.String(),
		"account": accs[0].Address,
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

	var respJson responses.WalletContainsResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "1", respJson.Exists)

	// Bad request
	reqBody = map[string]interface{}{
		"action":  "wallet_contains",
		"wallet":  wallet.ID.String(),
		"account": accs[0].Address + "1234",
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

	var errEsp map[string]interface{}
	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &errEsp)

	assert.Equal(t, "Invalid account", errEsp["error"])

	// Valid account that doesnt exist in wallet
	reqBody = map[string]interface{}{
		"action":  "wallet_contains",
		"wallet":  wallet.ID.String(),
		"account": "nano_1jtx5p8141zjtukz4msp1x93st7nh475f74odj8673qqm96xczmtcnanos1o",
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

	assert.Equal(t, "0", respJson.Exists)
}

func TestWalletRepresentativeSet(t *testing.T) {
	newSeed, _ := utils.GenerateSeed(strings.NewReader("b40ee7fa4a110bb17e7706e30eeed3bc360f571ddde6b436d7926ed3e77449f2"))
	wallet, _ := MockController.Wallet.WalletCreate(newSeed)
	// Request JSON
	reqBody := map[string]interface{}{
		"action":         "wallet_representative_set",
		"wallet":         wallet.ID.String(),
		"representative": "nano_1jtx5p8141zjtukz4msp1x93st7nh475f74odj8673qqm96xczmtcnanos1o",
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

	var respJson responses.SetResponse
	respBody, _ := io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &respJson)

	assert.Equal(t, "1", respJson.Set)

	// Bad request
	reqBody = map[string]interface{}{
		"action":         "wallet_representative_set",
		"wallet":         wallet.ID.String(),
		"representative": "1234",
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

	var errEsp map[string]interface{}
	respBody, _ = io.ReadAll(resp.Body)
	json.Unmarshal(respBody, &errEsp)

	assert.Equal(t, "Invalid representative account", errEsp["error"])
}
