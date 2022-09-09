package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"k8s.io/klog/v2"
)

// This is called the "Gateway" because it's the entry point for all requests
// This API is intended to replace the nano node wallet RPCs
// https://docs.nano.org/commands/rpc-protocol/#wallet-rpcs
// It will:
// 1) Determine if the request is a supported wallet RPC, if so process it
// 2) If it's a wallet RPC we don't support, return error
// 3) Other requests with a correct signature go straight to the node
// The error messages and behavior are also intended to replace what the nano node returns
// The node isn't exactly great at returning errors, and the error messages are not very helpful
// But as we want to be a drop-in replacement we mimic the behavior
func (hc *HttpController) Gateway(w http.ResponseWriter, r *http.Request) {
	var baseRequest map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&baseRequest); err != nil {
		klog.Errorf("Error unmarshalling http base request %s", err)
		ErrUnableToParseJson(w, r)
		return
	}

	if _, ok := baseRequest["action"]; !ok {
		ErrUnableToParseJson(w, r)
		return
	}

	action := strings.ToLower(fmt.Sprintf("%v", baseRequest["action"]))

	switch action {
	case "wallet_create":
		hc.HandleWalletCreate(&baseRequest, w, r)
		return
	}
}
