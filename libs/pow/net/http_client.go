package net

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"k8s.io/klog/v2"
)

// Base request
func MakeRequest(url string, request interface{}) ([]byte, error) {
	requestBody, _ := json.Marshal(request)
	// HTTP post
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		klog.Errorf("Error building request %s", err)
		return nil, err
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	resp, err := Client.Do(httpRequest)
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
