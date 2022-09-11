package pow

import (
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

// We want to test that we can invoke to numerous peers and get the first response
func TestWorkGeneratePeers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

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
				"work": "abcd1234",
			})
			return resp, err
		},
	)

	result, err := PPow.WorkGeneratePeers("abc", "ffffffffffff")
	assert.Nil(t, err)
	assert.Equal(t, "abcd1234", result)

	// In this case #1 is lazy
	httpmock.RegisterResponder("POST", "https://workerurl1.com",
		func(req *http.Request) (*http.Response, error) {
			time.Sleep(10 * time.Second)
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "abcd1234",
			})
			return resp, err
		},
	)

	// This responder is slow/lazt
	httpmock.RegisterResponder("POST", "https://workerurl2.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "abcd1234",
			})
			return resp, err
		},
	)

	result, err = PPow.WorkGeneratePeers("abcdef", "ffffffffffff")
	assert.Nil(t, err)
	assert.Equal(t, "abcd1234", result)

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
				"work": "abcd1234",
			})
			return resp, err
		},
	)

	result, err = PPow.WorkGeneratePeers("abcdef", "ffffffffffff")
	assert.Nil(t, err)
	assert.Equal(t, "abcd1234", result)

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

	result, err = PPow.WorkGeneratePeers("abcdef", "ffffffffffff")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Unable to generate work")
}
