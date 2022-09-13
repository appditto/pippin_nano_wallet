package mocks

import (
	"bytes"
	"io"
)

var AccountBalancesResponseStr = "{\n  \"balances\": {\n    \"nano_1gyeqc6u5j3oaxbe5qy1hyz3q745a318kh8h9ocnpan7fuxnq85cxqboapu5\": {\n      \"balance\": \"11999999999999999918751838129509869131\",\n      \"pending\": \"0\",\n      \"receivable\": \"0\"\n    }\n  }\n}"
var AccountsBalancesResponse = io.NopCloser(bytes.NewReader([]byte(AccountBalancesResponseStr)))
var AccountsFrontiersResponseStr = "{\n  \"frontiers\" : {\n    \"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3\": \"791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A\",\n    \"nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7\": \"6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F\"\n  }\n}"
var AccountsFrontiersResponse = io.NopCloser(bytes.NewReader([]byte(AccountsFrontiersResponseStr)))
var AccountsPendingResponseStr = "{\n  \"blocks\" : {\n    \"nano_1111111111111111111111111111111111111111111111111117353trpda\": [\"142A538F36833D1CC78B94E11C766F75818F8B940771335C6C1B8AB880C5BB1D\"],\n    \"nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3\": [\"4C1FEEF0BEA7F50BE35489A1233FE002B212DEA554B55B1B470D78BD8F210C74\"]\n  }\n}"
var AccountsPendingResponse = io.NopCloser(bytes.NewReader([]byte(AccountsPendingResponseStr)))
