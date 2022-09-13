package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/database"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	"github.com/appditto/pippin_nano_wallet/libs/rpc"
	"github.com/appditto/pippin_nano_wallet/libs/wallet"
	"github.com/stretchr/testify/assert"
	"k8s.io/klog/v2"
)

// ! These API tests are higher level integration tests that test the API as a whole
// ! Specific functionality done in the API is tested in lower-level unit tests
// ! e.g., we just test that the API returns the response we expect, but don't verify it exists in the database

var MockController *HttpController

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	os.Setenv("MOCK_REDIS", "true")
	defer os.Unsetenv("MOCK_REDIS")
	// We use an in-memory sqlite database for testing
	ctx := context.Background()
	dbconn, err := database.GetSqlDbConn(true)
	if err != nil {
		klog.Fatalf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	entClient, err := database.NewEntClient(dbconn)
	defer entClient.Close()
	if err != nil {
		klog.Fatalf("Failed to create ent client: %v", err)
		os.Exit(1)
	}

	//Create schema
	if err := entClient.Schema.Create(ctx); err != nil {
		klog.Fatalf("Failed to run migrations: %v", err)
		os.Exit(1)
	}

	// Setup nano wallet
	wallet := wallet.NanoWallet{
		DB:     entClient,
		Ctx:    ctx,
		Banano: false,
	}

	MockController = &HttpController{
		Wallet: &wallet,
		RpcClient: &rpc.RPCClient{
			Url: "http://localhost:123456",
		},
		PowClient: pow.NewPippinPow([]string{}, "", ""),
	}
	return m.Run()
}

func TestBadJson(t *testing.T) {
	// Request JSON
	reqBody := map[string]interface{}{
		"badjson": "badjson",
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

	assert.Equal(t, "Unable to parse json", respJson["error"])
}
