package responses

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAccountsBalancesResponse(t *testing.T) {
	encoded := "{\"balances\" : {\"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3\": {\"balance\": \"325586539664609129644855132177\",\"pending\": \"2309372032769300000000000000000000\",\"receivable\": \"2309372032769300000000000000000000\"},\"nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7\":{\"balance\": \"10000000\",\"pending\": \"0\",\"receivable\": \"0\" }}}"
	var decoded AccountsBalancesResponse
	json.Unmarshal([]byte(encoded), &decoded)
	balances := *decoded.Balances
	assert.Equal(t, "325586539664609129644855132177", balances["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"].Balance)
	assert.Equal(t, "2309372032769300000000000000000000", balances["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"].Pending)
	assert.Equal(t, "2309372032769300000000000000000000", balances["nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3"].Receivable)
	assert.Equal(t, "10000000", balances["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"].Balance)
	assert.Equal(t, "0", balances["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"].Pending)
	assert.Equal(t, "0", balances["nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7"].Receivable)
}

func TestDecodeAccountsBalancesResponseError(t *testing.T) {
	encoded := "{\"error\": \"Account not found\"}"
	var decoded AccountsBalancesResponse
	json.Unmarshal([]byte(encoded), &decoded)
	assert.Nil(t, decoded.Balances)
}
