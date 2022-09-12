package boompow

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

type GQLError string

const (
	InvalidUsernamePasssword GQLError = "Invalid username or password"
	ServerError                       = "Unknown server error, try again later"
)

type BpowClient struct {
	client graphql.Client
	url    string
}
type authedTransport struct {
	wrapped http.RoundTripper
	token   string
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.token)
	return t.wrapped.RoundTrip(req)
}

func NewBpowClient(url string, token string) *BpowClient {
	var gqlClient graphql.Client
	gqlClient = graphql.NewClient(url, &http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport, token: token}})
	return &BpowClient{
		client: gqlClient,
		url:    url,
	}
}

func (c *BpowClient) WorkGenerate(ctx context.Context, hash string, difficultyMultipler int, blockAward bool, bpowKey string) (string, error) {
	var gqlClient graphql.Client
	if bpowKey != "" {
		gqlClient = graphql.NewClient(c.url, &http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport, token: bpowKey}})
	} else {
		gqlClient = c.client
	}

	resp, err := workGenerate(ctx, gqlClient, WorkGenerateInput{
		Hash:                 hash,
		DifficultyMultiplier: difficultyMultipler,
		BlockAward:           blockAward,
	})

	if err != nil {
		fmt.Printf("Error generating work %v", err)
		return "", err
	}

	return resp.WorkGenerate, nil
}
