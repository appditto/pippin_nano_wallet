package pow

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var PPow *PippinPow

func TestMain(m *testing.M) {
	os.Setenv("BPOW_KEY", "1234")
	defer os.Unsetenv("BPOW_KEY")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := io.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		if !strings.Contains(bodyString, "boompowaccepted") {
			w.WriteHeader(400)
			return
		}
		if r.Header.Get("Authorization") == "overriden" {
			w.Header().Add("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": map[string]interface{}{
					"workGenerate": "newworkfromboompow",
				},
			})
			return
		}
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"workGenerate": "boompowwork",
			},
		})
	}))
	defer server.Close()
	os.Setenv("BPOW_URL", server.URL)
	defer os.Unsetenv("BPOW_URL")

	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	PPow = NewPippinPow(
		[]string{
			"https://workerurl1.com",
			"https://workerurl2.com",
		},
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

	httpmock.RegisterResponder("POST", "http://boompowurl.notreal.com/graphql",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "abcd1234",
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

	result, err := PPow.WorkGenerateMeta("abc", 1, false, true, "")
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

	result, err = PPow.WorkGenerateMeta("abcdef", 1, false, true, "")
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

	result, err = PPow.WorkGenerateMeta("abcdef", 1, false, true, "")
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

	result, err = PPow.WorkGenerateMeta("abcdef", 1, false, true, "")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Unable to generate work")

	// BoomPoW only one that works with this request

	result, err = PPow.WorkGenerateMeta("boompowaccepted", 1, false, true, "")
	assert.Nil(t, err)
	assert.Equal(t, "boompowwork", result)

	// Test with local pow (no peers, no boompow configured)
	ppow := &PippinPow{
		WorkPeers: []string{},
	}

	result, err = ppow.WorkGenerateMeta("09263b65752d05ce4df5aeed849ffc2be5bf47026abb4fa5879359ae571ba9c8", 1, true, true, "")
	assert.Nil(t, err)
	assert.Len(t, result, 16)
}
