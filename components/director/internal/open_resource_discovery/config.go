package open_resource_discovery

type CertificateConfig struct {
	ClientID     string `envconfig:"APP_ORD_CERT_CLIENT_ID"`
	ClientSecret string `envconfig:"APP_ORD_CERT_CLIENT_SECRET"`
	OAuthURL     string `envconfig:"APP_ORD_CERT_OAUTH_URL"`

	CrsEndpoint string `envconfig:"APP_ORD_CERT_CSR_ENDPOINT"`
	Subject     string `envconfig:"APP_ORD_CERT_SUBJECT"`
}
