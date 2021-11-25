package http

import (
	"fmt"
	"github.com/kyma-incubator/compass/components/director/pkg/log"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

var (
	cachedToken []byte
)

// DefaultServiceAccountTokenPath missing godoc
const DefaultServiceAccountTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"

// InternalAuthorizationHeader missing godoc
const InternalAuthorizationHeader = "X-Authorization"

// NewServiceAccountTokenTransport constructs an serviceAccountTokenTransport
func NewServiceAccountTokenTransport(roundTripper HTTPRoundTripper) *serviceAccountTokenTransport {
	return &serviceAccountTokenTransport{
		roundTripper: roundTripper,
	}
}

// NewServiceAccountTokenTransportWithHeader constructs an serviceAccountTokenTransport with configurable header name
func NewServiceAccountTokenTransportWithHeader(roundTripper HTTPRoundTripper, headerName string) *serviceAccountTokenTransport {
	return &serviceAccountTokenTransport{
		roundTripper: roundTripper,
		headerName:   headerName,
	}
}

// NewServiceAccountTokenTransportWithPath constructs an serviceAccountTokenTransport with a given path
func NewServiceAccountTokenTransportWithPath(roundTripper HTTPRoundTripper, path string) *serviceAccountTokenTransport {
	return &serviceAccountTokenTransport{
		roundTripper: roundTripper,
		path:         path,
	}
}

// serviceAccountTokenTransport is transport that attaches a kubernetes service account token in the X-Authorization header for internal authentication.
type serviceAccountTokenTransport struct {
	roundTripper HTTPRoundTripper
	path         string
	headerName   string
}

// RoundTrip attaches a kubernetes service account token in the X-Authorization header for internal authentication.
func (tr *serviceAccountTokenTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	path := tr.path
	if len(path) == 0 {
		path = DefaultServiceAccountTokenPath
	}

	log.C(r.Context()).Info("[LOGGER] Inside Service Account Token Transport")
	fmt.Println("[FMT PRINTLN] Inside Service Account Token Transport")

	var token []byte
	if len(cachedToken) != 0 {
		log.C(r.Context()).Info("Will reuse Service Account token from cache")
		token = cachedToken
	} else {
		log.C(r.Context()).Info("Service Account token missing in cache")
		tkn, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read service account token file")
		}
		token = tkn
		cachedToken = tkn
	}

	headerName := InternalAuthorizationHeader
	if tr.headerName != "" {
		headerName = tr.headerName
	}
	r.Header.Set(headerName, "Bearer "+string(token))

	return tr.roundTripper.RoundTrip(r)
}
