package net

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/pow/models"
)

// Base request
func MakeRequest(ctx context.Context, url string, request interface{}, authorization string) ([]byte, error) {
	requestBody, _ := json.Marshal(request)
	// HTTP post
	httpRequest, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Errorf("Error building request %s", err)
		return nil, err
	}
	httpRequest.Header.Add("Content-Type", "application/json")
	if authorization != "" {
		httpRequest.Header.Add("Authorization", authorization)
	}
	httpRequest.WithContext(ctx)
	client := &http.Client{}
	resp, err := client.Do(httpRequest)
	if err != nil {
		log.Errorf("Error making RPC request %s", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Try to decode+deserialize
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error decoding response body %s", err)
		return nil, err
	}
	return body, nil
}

func MakeWorkGenerateRequest(ctx context.Context, url string, hash string, difficulty string) (*models.WorkGenerateResponse, error) {
	request := models.WorkGenerateRequest{
		WorkBaseRequest: models.WorkBaseRequest{
			Action: "work_generate",
			Hash:   hash,
		},
		Difficulty: difficulty,
	}
	response, err := MakeRequest(ctx, url, request, "")
	if err != nil {
		log.Errorf("Error making request %s", err)
		return nil, err
	}
	var resp models.WorkGenerateResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return nil, err
	}
	// Check that it's not empty
	if resp.Work == "" {
		return nil, errors.New("Unable to generate work")
	}

	return &resp, nil
}

// We don't care about the response for work cancel
func MakeWorkCancelRequest(ctx context.Context, url string, hash string) error {
	request := models.WorkBaseRequest{
		Action: "work_cancel",
		Hash:   hash,
	}
	_, err := MakeRequest(ctx, url, request, "")
	if err != nil {
		log.Errorf("Error making request %s", err)
		return err
	}
	return nil
}

func MakeBoompowWorkGenerateRequest(ctx context.Context, url string, bpowKey string, hash string, difficultyMultiplier int, blockAward bool) (string, error) {
	request := models.BoompowWorkGenerateRequest{
		Query: "mutation($hash:String!, $difficultyMultiplier: Int!, $blockAward: Boolean) { workGenerate(input:{hash:$hash, difficultyMultiplier:$difficultyMultiplier, blockAward: $blockAward}) }",
		Variables: models.BoompowWorkGenerateRequestVariables{
			Hash:                 hash,
			DifficultyMultiplier: difficultyMultiplier,
			BlockAward:           blockAward,
		},
	}
	response, err := MakeRequest(ctx, url, request, bpowKey)
	if err != nil {
		log.Errorf("Error making request %s", err)
		return "", err
	}
	var resp models.BoompowResponse
	err = json.Unmarshal(response, &resp)
	if err != nil {
		log.Errorf("Error unmarshalling response %s", err)
		return "", err
	}
	// Check that it's not empty
	if resp.Data.WorkGenerate == "" {
		return "", errors.New("Unable to generate work")
	}

	return resp.Data.WorkGenerate, nil
}
