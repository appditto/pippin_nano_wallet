package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBlockInfoResponse(t *testing.T) {
	encoded := "{\n  \"block_account\": \"nano_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est\",\n  \"amount\": \"30000000000000000000000000000000000\",\n  \"balance\": \"5606157000000000000000000000000000000\",\n  \"height\": \"58\",\n  \"local_timestamp\": \"0\",\n  \"successor\": \"8D3AB98B301224253750D448B4BD997132400CEDD0A8432F775724F2D9821C72\",\n  \"confirmed\": \"true\",\n  \"contents\": {\n    \"type\": \"state\",\n    \"account\": \"nano_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est\",\n    \"previous\": \"CE898C131AAEE25E05362F247760F8A3ACF34A9796A5AE0D9204E86B0637965E\",\n    \"representative\": \"nano_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou\",\n    \"balance\": \"5606157000000000000000000000000000000\",\n    \"link\": \"5D1AA8A45F8736519D707FCB375976A7F9AF795091021D7E9C7548D6F45DD8D5\",\n    \"link_as_account\": \"nano_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z\",\n    \"signature\": \"82D41BC16F313E4B2243D14DFFA2FB04679C540C2095FEE7EAE0F2F26880AD56DD48D87A7CC5DD760C5B2D76EE2C205506AA557BF00B60D8DEE312EC7343A501\",\n    \"work\": \"8a142e07a10996d5\"\n  },\n  \"subtype\": \"send\"\n}"

	var decoded BlockInfoResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "nano_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est", decoded.BlockAccount)
	assert.Equal(t, "30000000000000000000000000000000000", decoded.Amount)
	assert.Equal(t, "5606157000000000000000000000000000000", decoded.Balance)
	assert.Equal(t, "58", decoded.Height)
	assert.Equal(t, "0", decoded.LocalTimestamp)
	assert.Equal(t, "8D3AB98B301224253750D448B4BD997132400CEDD0A8432F775724F2D9821C72", decoded.Successor)
	assert.Equal(t, "true", decoded.Confirmed)
	assert.Equal(t, "state", decoded.Contents.Type)
	assert.Equal(t, "nano_1ipx847tk8o46pwxt5qjdbncjqcbwcc1rrmqnkztrfjy5k7z4imsrata9est", decoded.Contents.Account)
	assert.Equal(t, "CE898C131AAEE25E05362F247760F8A3ACF34A9796A5AE0D9204E86B0637965E", decoded.Contents.Previous)
	assert.Equal(t, "nano_1stofnrxuz3cai7ze75o174bpm7scwj9jn3nxsn8ntzg784jf1gzn1jjdkou", decoded.Contents.Representative)
	assert.Equal(t, "5606157000000000000000000000000000000", decoded.Contents.Balance)
	assert.Equal(t, "5D1AA8A45F8736519D707FCB375976A7F9AF795091021D7E9C7548D6F45DD8D5", decoded.Contents.Link)
	assert.Equal(t, "nano_1qato4k7z3spc8gq1zyd8xeqfbzsoxwo36a45ozbrxcatut7up8ohyardu1z", decoded.Contents.LinkAsAccount)
	assert.Equal(t, "82D41BC16F313E4B2243D14DFFA2FB04679C540C2095FEE7EAE0F2F26880AD56DD48D87A7CC5DD760C5B2D76EE2C205506AA557BF00B60D8DEE312EC7343A501", decoded.Contents.Signature)
	assert.Equal(t, "8a142e07a10996d5", decoded.Contents.Work)
	assert.Equal(t, "send", decoded.Subtype)
}
