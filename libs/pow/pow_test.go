package pow

import (
	"os"
	"testing"
)

var PPow *PippinPow

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	PPow = &PippinPow{
		WorkPeers: []string{
			"https://workerurl1.com",
			"https://workerurl2.com",
		},
	}
	return m.Run()
}
