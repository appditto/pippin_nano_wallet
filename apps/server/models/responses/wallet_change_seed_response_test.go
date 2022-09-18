package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeWalletChangeSeedResponse(t *testing.T) {
	response := WalletChangeSeedResponse{
		Success:             "",
		LastRestoredAccount: "ban_1234",
		RestoredCount:       1,
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"success\":\"\",\"last_restored_account\":\"ban_1234\",\"restored_count\":1}", string(encoded))
}
