package open_resource_discovery

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type CsrRequest struct {
	Csr CsrPayload `json:"certificate-signing-request"`
	//AuthorityRequest string     `json:"authority-request"`
}

type CsrPayload struct {
	Value string `json:"value"`
	Type  string `json:"type"`
	//Validity time.Duration `json:"validity"`
}

type CsrResponse struct {
	CrtResponse CrtResponse `json:"certificate-response"`
}

type CrtResponse struct {
	Crt string `json:"value"`
}

type CertificateClient struct {
	lock        sync.RWMutex
	client      *http.Client
	csrEndpoint string
	subject     string
	key         *rsa.PrivateKey
	certificate *tls.Certificate
}

func NewCertificateClient(client *http.Client, csrEndpoint, subject string) (*CertificateClient, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	return &CertificateClient{
		lock:        sync.RWMutex{},
		client:      client,
		csrEndpoint: csrEndpoint,
		// subject:     ParseSubject(subject),
		subject:     subject,
		key:         privateKey,
		certificate: nil,
	}, nil
}

func (cc *CertificateClient) getCert(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	if cc.readCert() != nil {
		fmt.Println("Cert already exists...")
		return cc.readCert(), nil
	}

	// The order of the attributes within the subject is mandatory...
	//csrTemplate := x509.CertificateRequest{
	//	Subject: cc.subject,
	//}
	//
	//csr, err := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, cc.key)
	//if err != nil {
	//	return nil, errors.Wrap(err, "Failed to create CSR")
	//}
	//
	//pemEncodedCSR := pem.EncodeToMemory(&pem.Block{
	//	Type: "CERTIFICATE REQUEST", Bytes: csr,
	//})

	priv := PrivateKey{
		Type: Rsa{
			Bits: 2048,
		},
		Password: "S3CuRe_PRIVATE_KEY_pA5SpHrAs3=",
	}
	csr := Csr{
		Subject:    cc.subject,
		PrivateKey: priv.CsrStruct(),
	}

	csrRequest := CsrRequest{
		Csr: CsrPayload{
			Value: csr.PemStr(),
			Type:  "pkcs10-pem",
			//Validity: time.Minute * 100,
		},
		//AuthorityRequest: "<uuid>>",
	}
	body, err := json.Marshal(csrRequest)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to marshal")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, cc.csrEndpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create csr request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cc.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to call cert service")
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	csrResp := &CsrResponse{}
	err = json.Unmarshal(bodyBytes, csrResp)

	crt := CertificateChain{
		PemStr: csrResp.CrtResponse.Crt,
	}

	if err != nil {
		return nil, err
	}

	clientCert := &tls.Certificate{
		Certificate: crt.CertificateWithChain(),
		PrivateKey:  privateKeyPemToType(privateKeyPemStrToPem(csr.PrivateKey.PemStr, csr.PrivateKey.Password)),
	}
	cc.writeCert(clientCert)

	return clientCert, nil
}

func (cc *CertificateClient) readCert() *tls.Certificate {
	cc.lock.RLock()
	defer cc.lock.RUnlock()
	return cc.certificate
}

func (cc *CertificateClient) writeCert(cert *tls.Certificate) {
	cc.lock.Lock()
	cc.certificate = cert
	cc.lock.Unlock()
}

//func ParseSubject(subject string) pkix.Name {
//	subjectInfo := extractSubject(subject)
//
//	return pkix.Name{
//		CommonName:         subjectInfo["CN"][0],
//		Country:            subjectInfo["C"],
//		Organization:       subjectInfo["O"],
//		OrganizationalUnit: subjectInfo["OU"],
//		Locality:           subjectInfo["L"],
//	}
//}
//
//func extractSubject(subject string) map[string][]string {
//	result := map[string][]string{}
//
//	segments := strings.Split(subject, ",")
//
//	for _, segment := range segments {
//		segmentTrimed := strings.Trim(segment, " ")
//		parts := strings.Split(segmentTrimed, "=")
//		result[parts[0]] = append(result[parts[0]], parts[1:]...)
//	}
//
//	return result
//}

type CertificateTransport struct {
	transport *http.Transport
}

func NewCertificateTransport(config CertificateConfig, timeout time.Duration) (*CertificateTransport, error) {
	baseClient := &http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
		Timeout:   timeout,
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, baseClient)
	ccConf := clientcredentials.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		TokenURL:     config.OAuthURL,
		AuthStyle:    oauth2.AuthStyleAutoDetect,
	}
	client := ccConf.Client(ctx)
	certClient, err := NewCertificateClient(client, config.CrsEndpoint, config.Subject)
	if err != nil {
		return nil, err
	}

	return &CertificateTransport{
		transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				GetClientCertificate: certClient.getCert,
				InsecureSkipVerify:   false,
			},
		},
	}, nil
}

func (ct *CertificateTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return ct.transport.RoundTrip(r)
}
