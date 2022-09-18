package rpc

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/rpc/mocks"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/requests"
	walletmodels "github.com/appditto/pippin_nano_wallet/libs/wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var MockRpcClient *RPCClient

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	// Mock HTTP client
	MockRpcClient = &RPCClient{Url: "http://localhost:123456"}
	return m.Run()
}

func TestGetAccountsBalances(t *testing.T) {
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

	resp, err := MockRpcClient.MakeAccountsBalancesRequest([]string{"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"})
	assert.Nil(t, err)
	balances := *resp.Balances
	assert.Equal(t, "11999999999999999918751838129509869131", balances["nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5"].Balance)
}

func TestGetAccountBalance(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.AccountBalanceResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)

	resp, err := MockRpcClient.MakeAccountBalanceRequest("nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3")
	assert.Nil(t, err)
	assert.Equal(t, "10000", resp.Balance)
	assert.Equal(t, "10000", resp.Pending)
	assert.Equal(t, "10000", resp.Receivable)
}

func TestGetAccountsFrontiers(t *testing.T) {
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

	resp, err := MockRpcClient.MakeAccountsFrontiersRequest([]string{"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"})
	assert.Nil(t, err)
	frontiers := *resp.Frontiers
	assert.Equal(t, "791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A", frontiers["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"])
	assert.Equal(t, "6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F", frontiers["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"])
}

func TestGetAccountsPending(t *testing.T) {
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

	resp, err := MockRpcClient.MakeAccountsPendingRequest([]string{"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"})
	assert.Nil(t, err)
	blocks := *resp.Blocks
	assert.Equal(t, "4C1FEEF0BEA7F50BE35489A1233FE002B212DEA554B55B1B470D78BD8F210C74", blocks["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"][0])
	assert.Equal(t, "142A538F36833D1CC78B94E11C766F75818F8B940771335C6C1B8AB880C5BB1D", blocks["nano_1111111111111111111111111111111111111111111111111117353trpda"][0])
}

func TestMakeProcessRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var pr requests.ProcessRequest
			json.NewDecoder(req.Body).Decode(&pr)
			if pr.Block.Hash == "abcd1234" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.ProcessResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			}
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.ErrorResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)

	resp, err := MockRpcClient.MakeProcessRequest(requests.ProcessRequest{
		BaseRequest: requests.BaseRequest{
			Action: "process",
		},
		JsonBlock: true,
		Block: walletmodels.StateBlock{
			Hash: "abcd1234",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, "E2FB233EF4554077A7BF1AA85851D5BF0B36965D2B0FB504B2BC778AB89917D3", resp.Hash)

	// Make an error req
	resp, err = MockRpcClient.MakeProcessRequest(requests.ProcessRequest{
		BaseRequest: requests.BaseRequest{
			Action: "process",
		},
		JsonBlock: true,
		Block: walletmodels.StateBlock{
			Hash: "notabcd1234",
		},
	})
	assert.NotNil(t, err)
	assert.Equal(t, "bad input", err.Error())
}

func TestMakeBlockInfoRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var pr requests.BlockInfoRequest
			json.NewDecoder(req.Body).Decode(&pr)
			if pr.Hash == "abcd1234" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.BlockInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			}
			var js map[string]interface{}
			json.Unmarshal([]byte(mocks.ErrorResponseStr), &js)
			resp, err := httpmock.NewJsonResponse(200, js)
			return resp, err
		},
	)

	resp, err := MockRpcClient.MakeBlockInfoRequest("abcd1234")

	assert.Nil(t, err)
	assert.Equal(t, "nano_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est", resp.BlockAccount)
	assert.Equal(t, "30000000000000000000000000000000000", resp.Amount)
	assert.Equal(t, "5606157000000000000000000000000000000", resp.Balance)
	assert.Equal(t, "58", resp.Height)
	assert.Equal(t, "0", resp.LocalTimestamp)
	assert.Equal(t, "8D3AB98B301224253750D448B4BD997132400CEDD0A8432F775724F2D9821C72", resp.Successor)
	assert.Equal(t, "true", resp.Confirmed)
	assert.Equal(t, "state", resp.Contents.Type)
	assert.Equal(t, "nano_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est", resp.Contents.Account)
	assert.Equal(t, "CE898C131AAEE25E05362F247760F8A3ACF34A9796A5AE0D9204E86B0637965E", resp.Contents.Previous)
	assert.Equal(t, "nano_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou", resp.Contents.Representative)
	assert.Equal(t, "5606157000000000000000000000000000000", resp.Contents.Balance)
	assert.Equal(t, "5D1AA8A45F8736519D707FCB375976A7F9AF795091021D7E9C7548D6F45DD8D5", resp.Contents.Link)
	assert.Equal(t, "nano_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z", resp.Contents.LinkAsAccount)
	assert.Equal(t, "82D41BC16F313E4B2243D14DFFA2FB04679C540C2095FEE7EAE0F2F26880AD56DD48D87A7CC5DD760C5B2D76EE2C205506AA557BF00B60D8DEE312EC7343A501", resp.Contents.Signature)
	assert.Equal(t, "8a142e07a10996d5", resp.Contents.Work)
	assert.Equal(t, "send", resp.Subtype)

	// Make an error req
	resp, err = MockRpcClient.MakeBlockInfoRequest("def")
	assert.NotNil(t, err)
	assert.Equal(t, "bad input", err.Error())
}

func TestMakeAccountInfoRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "http://localhost:123456",
		func(req *http.Request) (*http.Response, error) {
			var pr requests.AccountInfoRequest
			json.NewDecoder(req.Body).Decode(&pr)
			if pr.Account == "abcd1234" {
				var js map[string]interface{}
				json.Unmarshal([]byte(mocks.AccountInfoResponseStr), &js)
				resp, err := httpmock.NewJsonResponse(200, js)
				return resp, err
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "Account not found",
			})
			return resp, err
		},
	)

	resp, err := MockRpcClient.MakeAccountInfoRequest("abcd1234")

	assert.Nil(t, err)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", resp.Frontier)
	assert.Equal(t, "0E3F07F7F2B8AEDEA4A984E29BFE1E3933BA473DD3E27C662EC041F6EA3917A0", resp.OpenBlock)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", resp.RepresentativeBlock)
	assert.Equal(t, "11999999999999999918751838129509869131", resp.Balance)
	assert.Equal(t, "1606934662", resp.ModifiedTimestamp)
	assert.Equal(t, "22966", resp.BlockCount)
	assert.Equal(t, "1", resp.AccountVersion)
	assert.Equal(t, "22966", resp.ConfirmedHeight)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", resp.ConfirmedFrontier)
	assert.Equal(t, "nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5", resp.Representative)
	assert.Equal(t, "11999999999999999918751838129509869131", resp.Weight)
	assert.Equal(t, "0", resp.Pending)
	assert.Equal(t, "0", resp.Receivable)
	assert.Equal(t, "0", resp.ConfirmedPending)
	assert.Equal(t, "0", resp.ConfirmedReceivable)

	// Make an error req
	resp, err = MockRpcClient.MakeAccountInfoRequest("def")
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrAccountNotFound)
}
