package destinationfetchersvc_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kyma-incubator/compass/components/director/internal/destinationfetchersvc"
	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"
	"github.com/kyma-incubator/compass/components/director/pkg/config"
	"github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/kyma-incubator/compass/components/director/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

const (
	sensitiveEndpoint  = "/destination-configuration/v1/destinations"
	subaccountEndpoint = "/destination-configuration/v1/subaccountDestinations"
	subdomain          = "test"
	tokenPath          = "/test"
	noPageCountHeader  = "noPageCount"
)

func TestClient_SubaccountEndpoint(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	mockClient, mockServerCloseFn, endpoint := fixHTTPClientSubaccount(t)
	defer mockServerCloseFn()

	apiConfig := destinationfetchersvc.APIConfig{
		GoroutineLimit:                    10,
		RetryInterval:                     time.Duration(100 * time.Millisecond),
		RetryAttempts:                     3,
		EndpointGetSubbacountDestinations: endpoint + subaccountEndpoint,
		EndpointFindDestination:           "",
		Timeout:                           time.Duration(100 * time.Millisecond),
		PageSize:                          100,
		PagingPageParam:                   "$page",
		PagingSizeParam:                   "$pageSize",
		PagingCountParam:                  "$pageCount",
		PagingCountHeader:                 "Page-Count",
	}

	cert, key := generateTestCertAndKey(t, "test")
	instanceCfg := config.InstanceConfig{}
	instanceCfg.Cert = string(cert)
	instanceCfg.Key = string(key)
	client, err := destinationfetchersvc.NewClient(instanceCfg, apiConfig, tokenPath, subdomain)

	require.NoError(t, err)
	client.SetHTTPClient(mockClient)

	t.Run("Success fetching data page 3", func(t *testing.T) {
		// WHEN
		res, err := client.FetchSubbacountDestinationsPage(ctx, "3")
		// THEN
		require.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("Success fetching data page but no Page-Count header is in response", func(t *testing.T) {
		// WHEN
		_, err := client.FetchSubbacountDestinationsPage(ctx, noPageCountHeader)
		// THEN
		require.ErrorContains(t, err, "failed to extract header")
	})

	t.Run("Fetch should fail with status code 500, but do three attempts", func(t *testing.T) {
		// WHEN
		_, err := client.FetchSubbacountDestinationsPage(ctx, "internalServerError")
		// THEN
		require.ErrorContains(t, err, "#3")
		require.ErrorContains(t, err, "status code 500")
	})

	t.Run("Fetch should fail with status code 4xx", func(t *testing.T) {
		// WHEN
		_, err := client.FetchSubbacountDestinationsPage(ctx, "forbidden")
		// THEN
		require.ErrorContains(t, err, "status code 403")
	})
}

func TestClient_SenstiveDataEndpoint(t *testing.T) {
	// GIVEN
	ctx := context.Background()
	mockClient, mockServerCloseFn, endpoint := fixHTTPClientSensitive(t)
	defer mockServerCloseFn()

	apiConfig := destinationfetchersvc.APIConfig{}
	apiConfig.EndpointFindDestination = endpoint + sensitiveEndpoint
	apiConfig.EndpointGetSubbacountDestinations = endpoint + subaccountEndpoint
	apiConfig.RetryAttempts = 3
	apiConfig.RetryInterval = time.Duration(100 * time.Millisecond)

	cert, key := generateTestCertAndKey(t, "test")
	instanceCfg := config.InstanceConfig{}
	instanceCfg.Cert = string(cert)
	instanceCfg.Key = string(key)
	client, err := destinationfetchersvc.NewClient(instanceCfg, apiConfig, tokenPath, subdomain)

	require.NoError(t, err)
	client.SetHTTPClient(mockClient)

	t.Run("Success fetching sensitive data", func(t *testing.T) {
		// WHEN
		res, err := client.FetchDestinationSensitiveData(ctx, "s4ext")
		// THEN
		require.NoError(t, err)
		assert.NotEmpty(t, res)
	})

	t.Run("Fetch should fail with status code 500, but do three attempts", func(t *testing.T) {
		// WHEN
		_, err := client.FetchDestinationSensitiveData(ctx, "internalServerError")
		// THEN
		require.ErrorContains(t, err, "#3")
		require.ErrorContains(t, err, "status code 500")
	})

	t.Run("NewNotFoundError should be returned for status 404", func(t *testing.T) {
		// WHEN
		_, err := client.FetchDestinationSensitiveData(ctx, "notFound")
		// THEN
		require.ErrorIs(t, err, apperrors.NewNotFoundError(resource.Destination, "notFound"))
	})

	t.Run("Error should be returned for status 400", func(t *testing.T) {
		// WHEN
		_, err := client.FetchDestinationSensitiveData(ctx, "badRequest")
		// THEN
		require.ErrorContains(t, err, "400")
	})
}

func fixHTTPClientSubaccount(t *testing.T) (*http.Client, func(), string) {
	mux := http.NewServeMux()

	mux.HandleFunc(subaccountEndpoint, func(w http.ResponseWriter, r *http.Request) {
		pageCount := r.URL.Query().Get("$pageCount")
		page := r.URL.Query().Get("$page")
		pageSize := r.URL.Query().Get("$pageSize")

		if page == "forbidden" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if page != "3" && page != noPageCountHeader {
			http.Error(w, "page number invalid", http.StatusInternalServerError)
			return
		}

		if pageSize != "100" {
			http.Error(w, "pageSize invalid", http.StatusInternalServerError)
			return
		}

		if pageCount != "true" {
			http.Error(w, "pageCount invalid", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if page != noPageCountHeader {
			w.Header().Set("Page-Count", "3")
		}

		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, fixSubaccountDestinationsEndpoint())
		require.NoError(t, err)
	})

	ts := httptest.NewServer(mux)

	return ts.Client(), ts.Close, ts.URL
}

func fixHTTPClientSensitive(t *testing.T) (*http.Client, func(), string) {
	mux := http.NewServeMux()

	mux.HandleFunc(sensitiveEndpoint+"/s4ext", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, fixSesnitiveDataJSON())
		require.NoError(t, err)
	})
	mux.HandleFunc(sensitiveEndpoint+"/internalServerError", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	mux.HandleFunc(sensitiveEndpoint+"/badRequest", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	})

	ts := httptest.NewServer(mux)

	return ts.Client(), ts.Close, ts.URL
}

func TestNewClient(t *testing.T) {
	const clientID = "client"
	const clientSecret = "secret"

	cert, key := generateTestCertAndKey(t, "test")

	t.Run("mtls+client-secret mode", func(t *testing.T) {

		instanceCfg := config.InstanceConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			URL:          "url",
			TokenURL:     "tokenURL",
			Cert:         string(cert),
			Key:          string(key),
		}

		client, err := destinationfetchersvc.NewClient(instanceCfg, destinationfetchersvc.APIConfig{}, "/oauth/token", "subdomain")
		require.NoError(t, err)

		httpClient := client.GetHTTPClient()
		tr, ok := httpClient.Transport.(*oauth2.Transport)
		require.True(t, ok, "expected *oauth2.Transport")

		certCfg := oauth.X509Config{
			Cert: string(cert),
			Key:  string(key),
		}

		tlsCert, err := certCfg.ParseCertificate()
		require.NoError(t, err)

		expectedTransport := &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{*tlsCert},
			},
		}
		require.Equal(t, tr.Base, expectedTransport)
	})
}

func fixSesnitiveDataJSON() string {
	return `{
		"s4ext": {
      "owner": {
        "SubaccountId": "8fb6ac72-124e-11ed-861d-0242ac120002",
        "InstanceId": null
      },
      "destinationConfiguration": {
        "Name": "s4ext",
        "Type": "HTTP",
        "URL": "https://kaladin.bg",
        "Authentication": "BasicAuthentication",
        "ProxyType": "Internet",
        "XFSystemName": "Rock",
        "HTML5.DynamicDestination": "true",
        "User": "Kaladin",
        "product.name": "SAP S/4HANA Cloud",
        "Password": "securePass",
      },
      "authTokens": [
        {
          "type": "Basic",
          "value": "blJhbHQ1==",
          "http_header": {
            "key": "Authorization",
            "value": "Basic blJhbHQ1=="
          }
        }
      ]
    }
  }`
}

func fixSubaccountDestinationsEndpoint() string {
	return `
  [
    {
      "Name": "string",
      "Type": "HTTP",
      "PropertyName": "string"
    },
    {
      "Name": "string",
      "Type": "HTTP",
      "PropertyName": "string"
    }
  ]`
}

func generateTestCertAndKey(t *testing.T, commonName string) (crtPem, keyPem []byte) {
	clientKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	template := &x509.Certificate{
		IsCA:         true,
		SerialNumber: big.NewInt(1234),
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	parent := template
	certRaw, err := x509.CreateCertificate(rand.Reader, template, parent, &clientKey.PublicKey, clientKey)
	require.NoError(t, err)

	crtPem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certRaw})
	keyPem = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(clientKey)})

	return
}
