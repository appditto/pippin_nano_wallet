package responses

//	{
//	  "frontiers" : {
//	    "nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3": "791AF413173EEE674A6FCF633B5DFC0F3C33F397F0DA08E987D9E0741D40D81A",
//	    "nano_3i1aq1cchnmbn9x5rsbap8b15akfh7wj7pwskuzi7ahz8oq6cobd99d4r3b7": "6A32397F4E95AF025DE29D9BF1ACE864D5404362258E06489FABDBA9DCCC046F"
//	  }
//	}
type AccountsFrontiersResponse struct {
	Frontiers *map[string]string `json:"frontiers,omitempty" mapstructure:"frontiers,omitempty"`
}
