package pow

type PippinPow struct {
	WorkPeers []string
	BpowKey   string
	BpowUrl   string // Not required, defaults to boompow.banano.cc/graphql
}
