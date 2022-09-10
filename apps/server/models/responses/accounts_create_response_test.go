package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeAccountsCreateResponse(t *testing.T) {
	response := AccountsCreateResponse{
		Accounts: []string{"account", "account2"},
	}
	encoded, err := json.Marshal(response)
	assert.Nil(t, err)
	assert.Equal(t, "{\"accounts\":[\"account\",\"account2\"]}", string(encoded))
}
