package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodePasswordChangeRequest(t *testing.T) {
	encoded := `{"action":"password_change","password":"1234","wallet":"1234"}`
	var decoded PasswordChangeRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "password_change", decoded.Action)
	assert.Equal(t, "1234", decoded.Password)
	assert.Equal(t, "1234", decoded.Wallet)

	encoded = `{"action":"password_change","password":"","wallet":"1234"}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "password_change", decoded.Action)
	assert.Equal(t, "", decoded.Password)
	assert.Equal(t, "1234", decoded.Wallet)
}

func TestMapStructureDecodePasswordChangeRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":   "password_change",
		"password": "1234",
		"wallet":   "1234",
	}
	var decoded PasswordChangeRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "password_change", decoded.Action)
	assert.Equal(t, "1234", decoded.Password)
	assert.Equal(t, "1234", decoded.Wallet)
}
