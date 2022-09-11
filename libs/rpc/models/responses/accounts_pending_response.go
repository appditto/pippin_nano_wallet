package responses

//	{
//	  "blocks" : {
//	    "nano_1111111111111111111111111111111111111111111111111117353trpda": ["142A538F36833D1CC78B94E11C766F75818F8B940771335C6C1B8AB880C5BB1D"],
//	    "nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3": ["4C1FEEF0BEA7F50BE35489A1233FE002B212DEA554B55B1B470D78BD8F210C74"]
//	  }
//	}
type AccountsPendingResponse struct {
	Blocks *map[string][]string `json:"blocks" mapstructure:"blocks"`
}
