package responses

//	{
//	  "blocks" : {
//	    "000D1BAEC8EC208142C99059B393051BAC8380F9B5A2E6B2489A277D81789F3F": "6000000000000000000000000000000"
//	  }
//	}
type ReceivableResponse struct {
	Blocks map[string]string `json:"blocks" mapstructure:"blocks"`
}
