package pow

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/appditto/pippin_nano_wallet/libs/utils"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var PPow *PippinPow

func TestMain(m *testing.M) {
	os.Setenv("BPOW_KEY", "1234")
	defer os.Unsetenv("BPOW_KEY")
	os.Setenv("BPOW_URL", "/fakeboompow")
	defer os.Unsetenv("BPOW_URL")

	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	PPow = NewPippinPow(
		[]string{
			"https://workerurl1.com",
			"https://workerurl2.com",
		},
		utils.GetEnv("BPOW_KEY", ""),
		utils.GetEnv("BPOW_URL", ""),
	)
	return m.Run()
}

func TestWorkGenerateLocal(t *testing.T) {
	// Test local pow generation
	work, err := PPow.generateWorkLocally("09263b65752d05ce4df5aeed849ffc2be5bf47026abb4fa5879359ae571ba9c8", 1)
	assert.Nil(t, err)
	assert.Len(t, work, 16)

	work, err = PPow.generateWorkLocally("abcdefg", 1)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid hash")
}

// We want to test that we can invoke to numerous peers and get the first response
func TestWorkGenerateMeta(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/fakeboompow",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "overridden" {
				resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
					"data": map[string]interface{}{
						"workGenerate": "customkey",
					},
				})
				return resp, err
			}
			if req.Header.Get("Authorization") == "1234" {
				resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
					"data": map[string]interface{}{
						"workGenerate": "boompowwork",
					},
				})
				return resp, err
			}

			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("POST", "https://workerurl1.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "abcd1234",
			})
			return resp, err
		},
	)

	// This responder is slow/lazt
	httpmock.RegisterResponder("POST", "https://workerurl2.com",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(10 * time.Second)
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "99999",
			})
			return resp, err
		},
	)

	result, err := PPow.WorkGenerateMeta("abc", 1, false, true, "bad")
	assert.Nil(t, err)
	assert.Equal(t, "abcd1234", result)

	// In this case #1 is lazy
	httpmock.RegisterResponder("POST", "https://workerurl1.com",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(10 * time.Second)
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "4444",
			})
			return resp, err
		},
	)

	// This responder is slow/lazt
	httpmock.RegisterResponder("POST", "https://workerurl2.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "5555",
			})
			return resp, err
		},
	)

	result, err = PPow.WorkGenerateMeta("abcdef", 1, false, true, "bad")
	assert.Nil(t, err)
	assert.Equal(t, "5555", result)

	// In this case one returns an error
	httpmock.RegisterResponder("POST", "https://workerurl1.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	httpmock.RegisterResponder("POST", "https://workerurl2.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "8888",
			})
			return resp, err
		},
	)

	result, err = PPow.WorkGenerateMeta("abcdef", 1, false, true, "bad")
	assert.Nil(t, err)
	assert.Equal(t, "8888", result)

	// In this case both return an error

	httpmock.RegisterResponder("POST", "https://workerurl1.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	// This responder is slow/lazt
	httpmock.RegisterResponder("POST", "https://workerurl2.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "error",
			})
			return resp, err
		},
	)

	result, err = PPow.WorkGenerateMeta("abcdef", 1, false, true, "bad")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Unable to generate work")

	// BoomPoW only one that works with this request
	result, err = PPow.WorkGenerateMeta("boompowaccepted", 1, false, true, "")
	assert.Nil(t, err)
	assert.Equal(t, "boompowwork", result)

	// Test BoomPoW overriding key
	result, err = PPow.WorkGenerateMeta("boompowaccepted", 1, false, true, "overridden")
	assert.Nil(t, err)
	assert.Equal(t, "customkey", result)

	// Test with local pow (no peers, no boompow configured)
	ppow := &PippinPow{
		WorkPeers: []string{},
	}

	result, err = ppow.WorkGenerateMeta("09263b65752d05ce4df5aeed849ffc2be5bf47026abb4fa5879359ae571ba9c8", 1, true, true, "")
	assert.Nil(t, err)
	assert.Len(t, result, 16)
}
