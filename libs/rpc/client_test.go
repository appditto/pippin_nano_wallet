package rpc

import (
	"net/http"
	"os"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/rpc/mocks"
	"github.com/stretchr/testify/assert"
)

var MockRpcClient *RPCClient

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	// Mock HTTP client
	Client = &mocks.MockClient{}
	MockRpcClient = &RPCClient{Url: "http://localhost:123456"}
	return m.Run()
}

func TestGetAccountsBalances(t *testing.T) {
	// Simulate response
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: mocks.AccountsBalancesResponse,
		}, nil
	}

	resp, err := MockRpcClient.MakeAccountsBalancesRequest([]string{"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"})
	assert.Nil(t, err)
	balances := *resp.Balances
	assert.Equal(t, "11999999999999999918751838129509869131", balances["nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5"].Balance)
}
