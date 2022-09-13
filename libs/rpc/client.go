package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/requests"
	"github.com/appditto/pippin_nano_wallet/libs/rpc/models/responses"
	"k8s.io/klog/v2"
)

type RPCClient struct {
	Url string
}

// Base request
func (client *RPCClient) MakeRequest(request interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(request)
	// HTTP post
	httpRequest, err := http.NewRequest(http.MethodPost, client.Url, bytes.NewBuffer(requestBody))
	if err != nil {
		klog.Errorf("Error building request %s", err)
		return nil, err
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	httpC := &http.Client{}
	resp, err := httpC.Do(httpRequest)
	if err != nil {
		klog.Errorf("Error making RPC request %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Try to decode+deserialize
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		klog.Errorf("Error decoding response body %s", err)
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
		klog.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp responses.AccountsBalancesResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		klog.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// Check that it'
	if resp.Balances == nil {
		return nil, errors.New("No balances returned")
	}

	return &resp, nil
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
		klog.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp responses.AccountsFrontiersResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		klog.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// Check that it'
	if resp.Frontiers == nil {
		return nil, errors.New("No frontiers returned")
	}

	return &resp, nil
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
		klog.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp responses.AccountsPendingResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		klog.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// Check that it'
	if resp.Blocks == nil {
		return nil, errors.New("No blocks returned")
	}

	return &resp, nil
}
