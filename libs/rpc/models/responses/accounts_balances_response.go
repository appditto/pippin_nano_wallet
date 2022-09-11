package responses

type AccountBalanceItem struct {
	Balance    string `json:"balance" mapstructure:"balance"`
	Pending    string `json:"pending" mapstructure:"pending"`
	Receivable string `json:"receivable" mapstructure:"receivable"`
}

//	{
//	  "balances" : {
//	    "nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3": {
//	        "balance": "325586539664609129644855132177",
//	        "pending": "2309372032769300000000000000000000",
//	        "receivable": "2309372032769300000000000000000000"
//	    },
//	    "nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7":
//	    {
//	      "balance": "10000000",
//	      "pending": "0",
//	      "receivable": "0"
//	    }
//	  }
//	}
type AccountsBalancesResponse struct {
	Balances *map[string]AccountBalanceItem `json:"balances,omitempty" mapstructure:"balances,omitempty"`
}
