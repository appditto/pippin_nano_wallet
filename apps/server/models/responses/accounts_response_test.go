package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAccountsResponse(t *testing.T) {
	response := AccountsResponse{
		Accounts: []string{"account", "account2"},
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"accounts\":[\"account\",\"account2\"]}", string(encoded))
}
