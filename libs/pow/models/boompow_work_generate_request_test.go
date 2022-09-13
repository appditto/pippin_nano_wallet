package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBoompowWorkGenerateRequest(t *testing.T) {
	encoded := "{\"query\":\"mutation($hash:String!, $difficultyMultiplier: Int!, $blockAward: Boolean) { workGenerate(input:{hash:$hash, difficultyMultiplier:$difficultyMultiplier, blockAward: $blockAward}) }\",\"variables\":{\"hash\":\"abc1234\",\"difficultyMultiplier\":1,\"blockAward\":true}}"
	req := BoompowWorkGenerateRequest{
		Query: "mutation($hash:String!, $difficultyMultiplier: Int!, $blockAward: Boolean) { workGenerate(input:{hash:$hash, difficultyMultiplier:$difficultyMultiplier, blockAward: $blockAward}) }",
		Variables: BoompowWorkGenerateRequestVariables{
			Hash:                 "abc1234",
			DifficultyMultiplier: 1,
			BlockAward:           true,
		}}
	encodedActual, _ := json.Marshal(&req)
	assert.Equal(t, encoded, string(encodedActual))
}
