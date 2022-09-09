package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeWalletCreateResponse(t *testing.T) {
	encoded := `{"wallet":"1"}`
	var decoded WalletCreateResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "1", decoded.Wallet)
}
