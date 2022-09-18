package wallet

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/appditto/pippin_nano_wallet/libs/database/ent"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/mocks"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/requests"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/appditto/pippin_nano_wallet/libs/utils/ed25519"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetBlockFromDatabase(t *testing.T) {
	// Predictable seed
	seed, _ := utils.GenerateSeed(strings.NewReader("42f567f406715877a8d60c6890b4655e8e0fe64f3f82089fe76899f180cd4021"))

	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	assert.Equal(t, false, wallet.Encrypted)
	assert.Equal(t, seed, wallet.Seed)

	// Create an account and an adhoc account
	account, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)

	_, priv, _ := ed25519.GenerateKey(strings.NewReader("1f729340e07eee69abac049c2fdd4a3c4b50e4672a2fabdf1ae295f2b4f3040d"))
	adhocAcct, err := MockWallet.AdhocAccountCreate(wallet, priv)
	assert.Nil(t, err)

	// Create some blocks for each
	// Create a block object
	_, err = MockWallet.DB.Block.Create().SetAccount(account).SetBlock(map[string]interface{}{
		"block": "hello",
	}).SetBlockHash("abc").SetSubtype("send").SetSendID("abc").Save(MockWallet.Ctx)
	assert.Nil(t, err)
	_, err = MockWallet.DB.Block.Create().SetAccount(adhocAcct).SetBlock(map[string]interface{}{
		"block": "world",
	}).SetBlockHash("def").SetSubtype("change").SetSendID("def").Save(MockWallet.Ctx)
	assert.Nil(t, err)

	// Retrieve block
	block, err := MockWallet.GetBlockFromDatabase(wallet, account.Address, "abc")
	assert.Nil(t, err)
	assert.Equal(t, "hello", block.Block["block"])
	block, err = MockWallet.GetBlockFromDatabase(wallet, adhocAcct.Address, "def")
	assert.Nil(t, err)
	assert.Equal(t, "world", block.Block["block"])
	block, err = MockWallet.GetBlockFromDatabase(wallet, account.Address, "nonexistent")
	assert.ErrorIs(t, err, ErrBlockNotFound)
}

func TestReceiveBlockCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/mockrpcendpoint",
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
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	_, err := MockWallet.createReceiveBlock(nil, nil, "", nil, nil)
	assert.ErrorIs(t, err, ErrInvalidWallet)
	_, err = MockWallet.createReceiveBlock(&ent.Wallet{}, nil, "", nil, nil)
	assert.ErrorIs(t, err, ErrInvalidAccount)

	// Create a wallet
	seed, err := utils.GenerateSeed(strings.NewReader("CA21ACED8B10297F40A3001BB97FC220B6E96AB236D36DC91D7B40A6852E05D1"))
	assert.Nil(t, err)
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create an account
	acc, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)

	// Receive a block
	work := "0000000000000000"
	block, err := MockWallet.createReceiveBlock(wallet, acc, "FB20236176F12827E71FD1F2928C8ABCBE6D2D9EE02E8BE8AC13F13AAF5575AE", &work, nil)
	assert.Nil(t, err)
	assert.Equal(t, "84ee43f56904a239e4bdd9f3e0835b0bc233416d7122e69fadddc1dba3e82cbe", block.Hash)
	assert.Equal(t, "state", block.Type)
	assert.Equal(t, acc.Address, block.Account)
	assert.Equal(t, "FB20236176F12827E71FD1F2928C8ABCBE6D2D9EE02E8BE8AC13F13AAF5575AE", block.Link)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", block.Previous)
	assert.Equal(t, "0000000000000000", block.Work)
}

func TestSendBlockCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/mockrpcendpoint",
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
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	_, err := MockWallet.createSendBlock(nil, nil, "", "", nil, nil)
	assert.ErrorIs(t, err, ErrInvalidWallet)
	_, err = MockWallet.createSendBlock(&ent.Wallet{}, nil, "", "", nil, nil)
	assert.ErrorIs(t, err, ErrInvalidAccount)

	// Create a wallet
	seed, err := utils.GenerateSeed(strings.NewReader("5F91A5BCC65ECE6912EBED8C33EE88048A0BD417B28553B513BC45C90D0FF1AB"))
	assert.Nil(t, err)
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create an account
	acc, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)

	// Receive a block
	work := "0000000000000000"
	block, err := MockWallet.createSendBlock(wallet, acc, "1", "nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj", &work, nil)
	assert.Nil(t, err)
	assert.Equal(t, "2d68dfea9879bffffb5b6dece2dd6e78221d01eab75d3726dc6791a244f073e3", block.Hash)
	assert.Equal(t, "state", block.Type)
	assert.Equal(t, acc.Address, block.Account)
	assert.Equal(t, "d4bbfa50649d80e5f63fc396c6f4cf6321cabd7c1480e964c2701d56aafeb5e3", block.Link)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", block.Previous)
	assert.Equal(t, "0000000000000000", block.Work)
}

func TestChangeBlockCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/mockrpcendpoint",
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
			}
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	_, err := MockWallet.createChangeBlock(nil, nil, "", nil, nil, true)
	assert.ErrorIs(t, err, ErrInvalidWallet)
	_, err = MockWallet.createChangeBlock(&ent.Wallet{}, nil, "", nil, nil, true)
	assert.ErrorIs(t, err, ErrInvalidAccount)

	// Create a wallet
	seed, err := utils.GenerateSeed(strings.NewReader("9c0784e354217e282e8b0177db3bf18d74768d8adb1de88412c6c4d9fd33f407"))
	assert.Nil(t, err)
	wallet, err := MockWallet.WalletCreate(seed)
	assert.Nil(t, err)

	// Create an account
	acc, err := MockWallet.AccountCreate(wallet, nil)
	assert.Nil(t, err)

	// Receive a block
	work := "0000000000000000"
	block, err := MockWallet.createChangeBlock(wallet, acc, "nano_3o7uzba8b9e1wqu5ziwpruteyrs3scyqr761x7ke6w1xctohxfh5du75qgaj", &work, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, "61595310547a16b7b7240eb65b09eb1b6994143ba74596f6e64883c5e3342150", block.Hash)
	assert.Equal(t, "state", block.Type)
	assert.Equal(t, acc.Address, block.Account)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000000", block.Link)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", block.Previous)
	assert.Equal(t, "0000000000000000", block.Work)

	// Test only if different
	// nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5
	block, err = MockWallet.createChangeBlock(wallet, acc, "nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5", &work, nil, true)
	assert.ErrorIs(t, err, ErrSameRepresentative)
}
