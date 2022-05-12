package fetcher_client

import (
	"net/http"

	"golang.org/x/oauth2/clientcredentials"
)

type HTTPRoundTripper interface {
	RoundTrip(*http.Request) (*http.Response, error)
}

func NewOAuth20Transport(roundTripper HTTPRoundTripper, provider clientcredentials.Config) *OAuth20Transport {
	return &OAuth20Transport{
		tokenProvider: provider,
		roundTripper:  roundTripper,
	}
}

type OAuth20Transport struct {
	roundTripper  HTTPRoundTripper
	tokenProvider clientcredentials.Config
}

func (c *OAuth20Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()
	token, err := c.tokenProvider.Token(ctx)
	if err != nil {
		return nil, err
	}
	token.SetAuthHeader(r)
	return c.roundTripper.RoundTrip(r)
}
