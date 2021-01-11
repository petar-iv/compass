package credentialsprovider

import "fmt"

// Config configuration for Credentials Provider client.
type Config struct {
	URL                   string `mapstructure:"url"`
	CredentialsIssuerPath string `mapstructure:"credentials_issuer_path"`
	IssuerBody            string `mapstructure:"issuer_body"`
	CredentialsPath       string `mapstructure:"credentials_path"`
	IssuerNameQueryParam  string `mapstructure:"issuer_name_query_param"`
	TenantIDQueryParam    string `mapstructure:"tenant_id_query_param"`
	TenantIDContextKey    string `mapstructure:"tenant_id_context_key"`
}

// DefaultConfig returns default configuration for credentials provider remote client. Mandatory fields are empty.
func DefaultConfig() *Config {
	return &Config{}
}

// Validate checks if the configuration has all the required fields.
func (c *Config) Validate() error {
	if c == nil {
		return fmt.Errorf("credentials provider config must be provided")
	}
	if len(c.URL) == 0 {
		return fmt.Errorf("URL cannot be empty")
	}
	if len(c.CredentialsIssuerPath) == 0 {
		return fmt.Errorf("CredentialsIssuerPath cannot be empty")
	}
	if len(c.IssuerBody) == 0 {
		return fmt.Errorf("IssuerBody cannot be empty")
	}
	if len(c.CredentialsPath) == 0 {
		return fmt.Errorf("CredentialsPath cannot be empty")
	}
	if len(c.IssuerNameQueryParam) == 0 {
		return fmt.Errorf("IssuerNameQueryParam cannot be empty")
	}
	if len(c.TenantIDQueryParam) == 0 {
		return fmt.Errorf("TenantIDQueryParam cannot be empty")
	}
	if len(c.TenantIDContextKey) == 0 {
		return fmt.Errorf("TenantIDContextKey cannot be empty")
	}
	return nil
}
