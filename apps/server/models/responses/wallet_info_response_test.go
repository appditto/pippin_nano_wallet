package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWalletInfoResponse(t *testing.T) {
	response := WalletInfoResponse{
		Balance:            "1",
		Pending:            "2",
		Receivable:         "2",
		AccountsCount:      5,
		AdhocCount:         1,
		DeterministicCount: 14,
		DeterministicIndex: 5,
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"balance\":\"1\",\"pending\":\"2\",\"receivable\":\"2\",\"accounts_count\":5,\"adhoc_count\":1,\"deterministic_count\":14,\"deterministic_index\":5}", string(encoded))
}
