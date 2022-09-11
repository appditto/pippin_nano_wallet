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

func TestGetAccountsFrontiers(t *testing.T) {
	// Simulate response
	mocks.GetDoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: mocks.AccountsFrontiersResponse,
		}, nil
	}

	resp, err := MockRpcClient.MakeAccountsFrontiersRequest([]string{"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"})
	assert.Nil(t, err)
	frontiers := *resp.Frontiers
	assert.Equal(t, "791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A", frontiers["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"])
	assert.Equal(t, "6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F", frontiers["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"])
}
