package mocks

import (
	"bytes"
	"io"
)

var AccountsBalancesResponse = io.NopCloser(bytes.NewReader([]byte("{\n  \"balances\": {\n    \"nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5\": {\n      \"balance\": \"11999999999999999918751838129509869131\",\n      \"pending\": \"0\",\n      \"receivable\": \"0\"\n    }\n  }\n}")))
