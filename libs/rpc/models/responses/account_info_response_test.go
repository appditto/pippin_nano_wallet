package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountInfoResponse(t *testing.T) {
	encoded := "{\n    \"frontier\": \"80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F\",\n    \"open_block\": \"0E3F07F7F2B8AEDEA4A984E29BFE1E3933BA473DD3E27C662EC041F6EA3917A0\",\n    \"representative_block\": \"80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F\",\n    \"balance\": \"11999999999999999918751838129509869131\",\n    \"confirmed_balance\": \"11999999999999999918751838129509869131\",\n    \"modified_timestamp\": \"1606934662\",\n    \"block_count\": \"22966\",\n    \"account_version\": \"1\",\n    \"confirmed_height\": \"22966\",\n    \"confirmed_frontier\": \"80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F\",\n    \"representative\": \"nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5\",\n    \"confirmed_representative\": \"nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5\",\n    \"weight\": \"11999999999999999918751838129509869131\",\n    \"pending\": \"0\",\n    \"receivable\": \"0\",\n    \"confirmed_pending\": \"0\",\n    \"confirmed_receivable\": \"0\"\n}"
	var decoded AccountInfoResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", decoded.Frontier)
	assert.Equal(t, "0E3F07F7F2B8AEDEA4A984E29BFE1E3933BA473DD3E27C662EC041F6EA3917A0", decoded.OpenBlock)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", decoded.RepresentativeBlock)
	assert.Equal(t, "11999999999999999918751838129509869131", decoded.Balance)
	assert.Equal(t, "1606934662", decoded.ModifiedTimestamp)
	assert.Equal(t, "22966", decoded.BlockCount)
	assert.Equal(t, "1", decoded.AccountVersion)
	assert.Equal(t, "22966", decoded.ConfirmedHeight)
	assert.Equal(t, "80A6745762493FA21A22718ABFA4F635656A707B48B3324198AC7F3938DE6D4F", decoded.ConfirmedFrontier)
	assert.Equal(t, "nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5", decoded.Representative)
	assert.Equal(t, "11999999999999999918751838129509869131", decoded.Weight)
	assert.Equal(t, "0", decoded.Pending)
	assert.Equal(t, "0", decoded.Receivable)
	assert.Equal(t, "0", decoded.ConfirmedPending)
	assert.Equal(t, "0", decoded.ConfirmedReceivable)
}
