package open_resource_discovery

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"io"
	"strings"

	"github.com/fullsailor/pkcs7"
	"github.com/youmark/pkcs8"
)

// KEY
// This type defines the interface to switch between real and mocked crypto/rsa implementation.
type cryptoRsaPkg interface {
	generateKey(random io.Reader, bits int) (*rsa.PrivateKey, error)
}

// This type implements the interface for real crypto/rsa.
type defaultCryptoRsa struct{}

// This function implements the interface for real crypto/rsa.GenerateKey.
func (d *defaultCryptoRsa) generateKey(random io.Reader, bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(random, bits)
}

// This variable block defines package-wide variables.
var (
	// This variable is the real and mocked crypto/rsa switch.
	cryptoRsa cryptoRsaPkg
)

// This function is a package-wide initialization.
func init() {
	// This sets the default to the real crypto/rsa.
	cryptoRsa = &defaultCryptoRsa{}
}

// This type is used to configure the RSA private key.
type Rsa struct {
	Bits int
}

// This type is used to collect all parameters for the private key.
type PrivateKey struct {
	Type     interface{}
	Password string
	key      *rsa.PrivateKey
}

// This function returns a RSA private key.
func (k *PrivateKey) rsa() *rsa.PrivateKey {
	attr, _ := k.Type.(Rsa)
	key, _ := cryptoRsa.generateKey(rand.Reader, attr.Bits)
	k.key = key
	return key
}

// This function returns the DER encoded private key.
func (k *PrivateKey) der(key interface{}) []byte {
	var (
		der []byte
	)
	if len(k.Password) > 0 {
		opts := pkcs8.Opts{
			Cipher: pkcs8.AES256CBC,
			KDFOpts: pkcs8.PBKDF2Opts{
				HMACHash: crypto.SHA256,
			},
		}
		der, _ = pkcs8.MarshalPrivateKey(key, []byte(k.Password), &opts)
	} else {
		der, _ = x509.MarshalPKCS8PrivateKey(key)
	}
	return der
}

// This function returns the PEM encoded private key.
func (k *PrivateKey) pem(key interface{}) []byte {
	var (
		pemBlock *pem.Block
	)
	if len(k.Password) > 0 {
		pemBlock = &pem.Block{Type: "ENCRYPTED PRIVATE KEY", Bytes: k.der(key)}
	} else {
		pemBlock = &pem.Block{Type: "PRIVATE KEY", Bytes: k.der(key)}
	}
	return pem.EncodeToMemory(pemBlock)
}

// This function returns a PEM encoded private key string.
func (k *PrivateKey) RsaPemStr() string {
	return string(k.pem(k.rsa()))
}

// This function returns a PEM encoded private key for CSR generation.
func (k *PrivateKey) CsrStruct() *CsrPrivateKey {
	return &CsrPrivateKey{
		PemStr:   k.RsaPemStr(),
		Password: k.Password,
	}
}

// This function returns the parsed PEM encoded private key.
func privateKeyPemStrToPem(pemStr string, password string) []byte {
	pemBlock, _ := pem.Decode([]byte(pemStr))
	_, err := x509.ParsePKCS8PrivateKey(pemBlock.Bytes)
	_, isEncrypted := err.(asn1.StructuralError)
	if isEncrypted {
		decryptedKey, _ := pkcs8.ParsePKCS8PrivateKey(pemBlock.Bytes, []byte(password))
		decryptedBytes, _ := x509.MarshalPKCS8PrivateKey(decryptedKey)
		pemBlock = &pem.Block{Type: "PRIVATE KEY", Bytes: decryptedBytes}
	}
	return pemBlock.Bytes
}

// This function returns the parsed private key.
func privateKeyPemToType(pem []byte) interface{} {
	der, _ := x509.ParsePKCS8PrivateKey(pem)
	return der
}

// CSR

// This type is used to prepare the private key for the CSR generation.
type CsrPrivateKey struct {
	PemStr   string
	Password string
}

// This type is used to collect all parameters for the CSR.
type Csr struct {
	Subject    string
	PrivateKey *CsrPrivateKey
}

// This function returns the DER encoded ordered subject RDN.
func (c *Csr) subjectRdnSequenceDer() []byte {
	var (
		rdnSeq pkix.RDNSequence
	)
	splits := strings.Split(c.Subject, ",")
	for _, s := range splits {
		keyValue := strings.Split(s, "=")
		if len(keyValue) == 2 {
			key := strings.TrimSpace(keyValue[0])
			value := keyValue[1]
			var (
				attrType asn1.ObjectIdentifier
			)
			switch key {
			case "C":
				attrType = asn1.ObjectIdentifier{2, 5, 4, 6}
			case "O":
				attrType = asn1.ObjectIdentifier{2, 5, 4, 10}
			case "OU":
				attrType = asn1.ObjectIdentifier{2, 5, 4, 11}
			case "L":
				attrType = asn1.ObjectIdentifier{2, 5, 4, 7}
			case "CN":
				attrType = asn1.ObjectIdentifier{2, 5, 4, 3}
			default:
				continue
			}
			rdnSet := pkix.RelativeDistinguishedNameSET{pkix.AttributeTypeAndValue{Type: attrType, Value: value}}
			rdnSeq = append(rdnSeq, rdnSet)
		}
	}
	rdnsDer, _ := asn1.Marshal(rdnSeq)
	return rdnsDer
}

// This function returns the DER encoded CSR.
func (c *Csr) der() []byte {
	csrTemplate := x509.CertificateRequest{
		RawSubject:         c.subjectRdnSequenceDer(),
		SignatureAlgorithm: x509.SHA256WithRSA,
	}
	der, _ := x509.CreateCertificateRequest(rand.Reader, &csrTemplate, privateKeyPemToType(privateKeyPemStrToPem(c.PrivateKey.PemStr, c.PrivateKey.Password)))
	return der
}

// This function returns the PEM encoded CSR.
func (c *Csr) pem() []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: c.der()})
}

// This function returns the PEM encoded CSR string.
func (c *Csr) PemStr() string {
	return string(c.pem())
}

// CRT

// This type is used to store the Certificate Service response for further use.
type CertificateChain struct {
	PemStr string
}

// This function returns the string array of PEM encoded certificates.
func (c *CertificateChain) convert() [][]byte {
	var (
		res [][]byte
	)
	pemBlock, _ := pem.Decode([]byte(c.PemStr))
	pkcs, _ := pkcs7.Parse(pemBlock.Bytes)
	for _, c := range pkcs.Certificates {
		res = append(res, c.Raw)
	}
	return res
}

// This function returns the first certificate of the array.
func (c *CertificateChain) Certificate() []byte {
	res := c.convert()[0]
	return res
}

// This function returns all certificates of the array.
func (c *CertificateChain) CertificateWithChain() [][]byte {
	var (
		res [][]byte
	)
	for _, c := range c.convert() {
		res = append(res, c)
	}
	return res
}
