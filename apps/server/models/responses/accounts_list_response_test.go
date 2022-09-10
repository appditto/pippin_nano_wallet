package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAccountsListResponse(t *testing.T) {
	response := AccountsListResponse{
		Accounts: []string{"account", "account2"},
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"accounts\":[\"account\",\"account2\"]}", string(encoded))
}
