package controller

import (
	"net/http"
	"strconv"

	"github.com/appditto/pippin_nano_wallet/apps/server/models/requests"
	"github.com/appditto/pippin_nano_wallet/apps/server/models/responses"
	"github.com/appditto/pippin_nano_wallet/libs/log"
	"github.com/appditto/pippin_nano_wallet/libs/pow"
	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
)

func (hc *HttpController) HandleWorkGenerate(rawRequest *map[string]interface{}, w http.ResponseWriter, r *http.Request) {
	var workRequest requests.WorkGenerateRequest
	if err := mapstructure.Decode(rawRequest, &workRequest); err != nil {
		log.Errorf("Error unmarshalling work_generate request %s", err)
		ErrUnableToParseJson(w, r)
		return
	} else if workRequest.Action == "" || workRequest.Hash == "" {
		ErrUnableToParseJson(w, r)
		return
	}

	// Validate the blcok hash
	if !utils.Validate64HexHash(workRequest.Hash) {
		ErrInvalidHash(w, r)
		return
	}

	// Difficulty is optional, we default to 64 if in nano mode, 1 if banano mode
	// Banano is always 1 no matter what
	difficulty := 1
	if !hc.Wallet.Banano {
		// If they override difficulty, convert it to a multiplier
		if workRequest.Difficulty != "" {
			difficultyUint, err := strconv.ParseUint(workRequest.Difficulty, 16, 64)
			// Their difficulty is invalid
			if err != nil {
				ErrUnableToParseJson(w, r)
			}
			difficulty = pow.MultiplierFromDifficulty(difficultyUint)
		} else if workRequest.Subtype == "receive" {
			// Receives always get 1
			difficulty = 1
		} else {
			// Default to 64 (nano send)
			difficulty = 64
		}
	}

	blockAward := true
	if workRequest.BlockAward != nil {
		blockAward = *workRequest.BlockAward
	}

	work, err := hc.PowClient.WorkGenerateMeta(workRequest.Hash, difficulty, true, blockAward, workRequest.BpowKey)
	if err != nil {
		log.Errorf("Error generating work %s", err)
		ErrWorkFailed(w, r)
		return
	}

	resp := responses.WorkGenerateResponse{
		Work: work,
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}
