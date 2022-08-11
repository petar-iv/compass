package destinationfetchersvc_test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	exampleDestination1 = `{
        "Name": "mys4_1",
        "Type": "HTTP",
        "URL": "https://localhost:443",
        "Authentication": "BasicAuthentication",
        "ProxyType": "Internet",
        "XFSystemName": "Test S4HANA system",
        "HTML5.DynamicDestination": "true",
        "User": "SOME_USER",
        "product.name": "SAP S/4HANA Cloud",
        "WebIDEEnabled": "true",
        "communicationScenarioId": "SAP_COM_0108",
        "Password": "SecretPassword",
        "WebIDEUsage": "odata_gen"
    }`
	exampleDestination2 = `{
        "Name": "mys4_2",
        "Type": "HTTP",
        "URL": "https://localhost:443",
        "Authentication": "BasicAuthentication",
        "ProxyType": "Internet",
        "XFSystemName": "Test S4HANA system",
        "HTML5.DynamicDestination": "true",
        "User": "SOME_USER",
        "product.name": "SAP S/4HANA Cloud",
        "WebIDEEnabled": "true",
        "communicationScenarioId": "SAP_COM_0109",
        "Password": "SecretPassword",
        "WebIDEUsage": "odata_gen"
    }`
)

type destinationHandler struct {
	t                        *testing.T
	tenantDestinationHandler func(w http.ResponseWriter, req *http.Request)
	fetchDestinationHandler  func(w http.ResponseWriter, req *http.Request)
	tokenHandler             func(w http.ResponseWriter, req *http.Request)
}

func (dh *destinationHandler) defaultTenantDestinationHandler(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	page := query.Get("$page")

	w.Header().Set("Content-Type", "application/json")
	// assuming pageSize is always 1
	switch page {
	case "1":
		w.Header().Set("Page-Count", "2")
		w.Write([]byte(fmt.Sprintf("[%s]", exampleDestination1)))
	case "2":
		w.Header().Set("Page-Count", "2")
		w.Write([]byte(fmt.Sprintf("[%s]", exampleDestination2)))
	default:
		dh.t.Logf("Expected page size to be 1 or 2, got '%s'", page)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (dh *destinationHandler) defaultFetchDestinationHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(500)
}

func (dh *destinationHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if strings.HasPrefix(path, "/subaccountDestinations") {
		if dh.tenantDestinationHandler == nil {
			dh.defaultTenantDestinationHandler(w, req)
		} else {
			dh.tenantDestinationHandler(w, req)
		}
		return
	}
	if strings.HasPrefix(path, "/destinations") {
		if dh.tenantDestinationHandler == nil {
			dh.defaultTenantDestinationHandler(w, req)
		} else {
			dh.tenantDestinationHandler(w, req)
		}
		return
	}
	if strings.HasPrefix(path, "/oauth/token") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write([]byte(`{
			"access_token": "accesstoken",
			"token_type": "tokentype",
			"refresh_token": "refreshtoken",
			"expires_in": 100
		}`))
		return
	}
	w.WriteHeader(500)
}

type destinationServer struct {
	server  *httptest.Server
	handler *destinationHandler
}

func newDestinationServer(t *testing.T) destinationServer {
	destinationHandler := &destinationHandler{t: t}
	httpServer := httptest.NewUnstartedServer(destinationHandler)
	var err error
	httpServer.Listener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	return destinationServer{
		server:  httpServer,
		handler: destinationHandler,
	}
}
