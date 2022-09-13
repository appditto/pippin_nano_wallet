package net

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(testMainWrapper(m))
}

func testMainWrapper(m *testing.M) int {
	// Mock HTTP client
	return m.Run()
}

func TestWorkGenerate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://workurl.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"work": "abcd1234",
			})
			return resp, err
		},
	)

	resp, err := MakeWorkGenerateRequest(context.TODO(), "https://workurl.com", "nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3", "ffffffffffff")
	assert.Nil(t, err)
	assert.Equal(t, "abcd1234", resp.Work)

	// Error
	httpmock.RegisterResponder("POST", "https://workurl.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"error": "abcd1234",
			})
			return resp, err
		},
	)

	resp, err = MakeWorkGenerateRequest(context.TODO(), "https://workurl.com", "nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3", "ffffffffffff")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "Unable to generate work")
}

func TestWorkCancel(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://workurl.com",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
				"success": "",
			})
			return resp, err
		},
	)

	err := MakeWorkCancelRequest(context.TODO(), "https://workurl.com", "nano_3t6k35gi95xu6tergt6p69ck76ogmitsa8mnijtpxm9fkcm736xtoncuohr3")
	assert.Nil(t, err)
}

func TestBoompowWorkGenerate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "/mockboompow",
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("Authorization") == "valid" {
				resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
					"data": map[string]interface{}{
						"workGenerate": "abcd1234",
					},
				})
				return resp, err
			}
			resp, err := httpmock.NewJsonResponse(403, map[string]interface{}{
				"error": "unauthorized",
			})
			return resp, err
		},
	)

	resp, err := MakeBoompowWorkGenerateRequest(context.TODO(), "/mockboompow", "invalidkey", "74681bdf45345dcef2f13fbaf9cedf5b30412ce9fcbe437d09677529cb0dbe9e", 1, true)
	assert.NotNil(t, err)
	assert.Len(t, resp, 0)

	resp, err = MakeBoompowWorkGenerateRequest(context.TODO(), "/mockboompow", "valid", "74681bdf45345dcef2f13fbaf9cedf5b30412ce9fcbe437d09677529cb0dbe9e", 1, true)
	assert.Nil(t, err)
	assert.Equal(t, "abcd1234", resp)
}
