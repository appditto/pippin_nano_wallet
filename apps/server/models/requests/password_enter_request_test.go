package requests

import (
	"encoding/json"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestDecodePasswordEnterRequest(t *testing.T) {
	encoded := `{"action":"password_enter","password":"1234","wallet":"1234"}`
	var decoded PasswordEnterRequest
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "password_enter", decoded.Action)
	assert.Equal(t, "1234", decoded.Password)
	assert.Equal(t, "1234", decoded.Wallet)

	encoded = `{"action":"password_enter","password":"","wallet":"1234"}`
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "password_enter", decoded.Action)
	assert.Equal(t, "", decoded.Password)
	assert.Equal(t, "1234", decoded.Wallet)
}

func TestMapStructureDecodePasswordEnterRequest(t *testing.T) {
	request := map[string]interface{}{
		"action":   "password_enter",
		"password": "1234",
		"wallet":   "1234",
	}
	var decoded PasswordEnterRequest
	mapstructure.Decode(request, &decoded)
	assert.Equal(t, "password_enter", decoded.Action)
	assert.Equal(t, "1234", decoded.Password)
	assert.Equal(t, "1234", decoded.Wallet)
}
