package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/requests"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/responses"
	"github.com/mitchellh/mapstructure"
)

var ErrAccountNotFound = errors.New("Account not found")

type RPCClient struct {
	Url        string
	httpClient *http.Client
}

func NewRPCClient(url string) *RPCClient {
	return &RPCClient{
		Url: url,
		httpClient: &http.Client{
			Timeout: time.Second * 30, // Set a timeout for all requests
		},
	}
}

// Base request
func (client *RPCClient) MakeRequest(request interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(request)
	if err != nil {
		log.Errorf("Error marshalling request %s", err)
		return nil, err
	}
	// HTTP post
	resp, err := client.httpClient.Post(client.Url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Errorf("Error making RPC request %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Try to decode+deserialize
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error decoding response body %s", err)
		return nil, err
	}
	return body, nil
}

func (client *RPCClient) MakeAccountsBalancesRequest(accounts []string) (*responses.AccountsBalancesResponse, error) {
	request := requests.AccountsRequest{
		BaseRequest: requests.BaseRequest{
			Action: "accounts_balances",
		},
		Accounts: accounts,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	var decoded responses.AccountsBalancesResponse
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error decoding response %s", err)
		return nil, err
	}

	if decoded.Balances == nil {
		return nil, errors.New("No balances returned")
	}

	return &decoded, nil
}

func (client *RPCClient) MakeAccountBalanceRequest(account string) (*responses.AccountBalanceItem, error) {
	request := requests.AccountRequest{
		BaseRequest: requests.BaseRequest{
			Action: "account_balance",
		},
		Account: account,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	var decoded responses.AccountBalanceItem
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error decoding response %s", err)
		return nil, err
	}

	if decoded.Balance == "" || decoded.Pending == "" {
		return nil, errors.New("No balance returned")
	}

	return &decoded, nil
}

func (client *RPCClient) MakeAccountsFrontiersRequest(accounts []string) (*responses.AccountsFrontiersResponse, error) {
	request := requests.AccountsRequest{
		BaseRequest: requests.BaseRequest{
			Action: "accounts_frontiers",
		},
		Accounts: accounts,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	var decoded responses.AccountsFrontiersResponse
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error decoding response %s", err)
		return nil, err
	}
	// Check that it'
	if decoded.Frontiers == nil {
		return nil, errors.New("No frontiers returned")
	}

	return &decoded, nil
}

func (client *RPCClient) MakeAccountsPendingRequest(accounts []string) (*responses.AccountsPendingResponse, error) {
	request := requests.AccountsRequest{
		BaseRequest: requests.BaseRequest{
			Action: "accounts_pending",
		},
		Accounts: accounts,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	var decoded responses.AccountsPendingResponse
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// Check that it'
	if decoded.Blocks == nil {
		return nil, errors.New("No blocks returned")
	}

	return &decoded, nil
}

func (client *RPCClient) MakeBlockInfoRequest(hash string) (*responses.BlockInfoResponse, error) {
	request := requests.BlockInfoRequest{
		BaseRequest: requests.BaseRequest{
			Action: "block_info",
		},
		Hash:      hash,
		JsonBlock: true,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	var decoded responses.BlockInfoResponse
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error decoding response %s", err)
		return nil, err
	}

	return &decoded, nil
}

func (client *RPCClient) MakeProcessRequest(request requests.ProcessRequest) (*responses.ProcessResponse, error) {
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	if val, ok := resp["hash"]; ok {
		hash, ok := val.(string)
		if !ok {
			return nil, errors.New("Hash response is not a string")
		}
		return &responses.ProcessResponse{
			Hash: hash,
		}, nil
	}
	if val, ok := resp["error"]; ok {
		err, ok := val.(string)
		if !ok {
			return nil, errors.New("Error response is not a string")
		}
		return nil, errors.New(err)
	}

	return nil, errors.New("No hash or error returned")
}

// {"error":"Account not found"}
func (client *RPCClient) MakeAccountInfoRequest(account string) (*responses.AccountInfoResponse, error) {
	includeAll := true
	request := requests.AccountInfoRequest{
		AccountRequest: requests.AccountRequest{
			BaseRequest: requests.BaseRequest{
				Action: "account_info",
			},
			Account: account,
		},
		Representative:   &includeAll,
		Weight:           &includeAll,
		Pending:          &includeAll,
		IncludeConfirmed: &includeAll,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			if strings.ToLower(errStr) == "account not found" {
				return nil, ErrAccountNotFound
			}
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	var decoded responses.AccountInfoResponse
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error decoding response %s", err)
		return nil, err
	}

	return &decoded, nil
}

func (client *RPCClient) MakeReceivableRequest(account string, threshold string) (*responses.ReceivableResponse, error) {
	request := requests.ReceivableRequest{
		BaseRequest: requests.BaseRequest{
			Action: "receivable",
		},
		Account:              account,
		Threshold:            threshold,
		IncludeOnlyConfirmed: true,
	}
	response, err := client.MakeRequest(request)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// See if contains an error
	if val, ok := resp["error"]; ok {
		errStr, ok := val.(string)
		if ok {
			if strings.ToLower(errStr) == "account not found" {
				return nil, ErrAccountNotFound
			}
			return nil, errors.New(errStr)
		}
		return nil, errors.New("Unknown error")
	}
	// The node is inconsistent and will give an empty string if nothing is here
	if val, ok := resp["blocks"]; ok {
		if val == "" {
			return &responses.ReceivableResponse{
				Blocks: make(map[string]string),
			}, nil
		}
	}
	// Decode properly
	var decoded responses.ReceivableResponse
	err = mapstructure.Decode(resp, &decoded)
	if err != nil {
		log.Errorf("Error decoding response %s", err)
		return nil, err
	}

	return &decoded, nil
}
