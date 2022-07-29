package config

import (
	"io/ioutil"
	"strings"

	"github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type DestinationInstanceConfig struct {
	XsAppName    string
	ClientID     string
	ClientSecret string
	URL          string
	AuthURL      string
	Cert         string
	Key          string
}

type DestinationsConfig struct {
	InstanceXsAppNamePath    string `envconfig:"APP_DESTINATION_INSTANCE_XSAPPNAME_PATH"`
	InstanceClientIDPath     string `envconfig:"APP_DESTINATION_INSTANCE_CLIENT_ID_PATH"`
	InstanceClientSecretPath string `envconfig:"APP_DESTINATION_INSTANCE_CLIENT_SECRET_PATH"`
	InstanceURLPath          string `envconfig:"APP_DESTINATION_INSTANCE_URL_PATH"`
	InstanceAuthURLPath      string `envconfig:"APP_DESTINATION_INSTANCE_AUTH_URL_PATH"`
	InstanceCertPath         string `envconfig:"APP_DESTINATION_INSTANCE_X509_CERT_PATH"`
	InstanceKeyPath          string `envconfig:"APP_DESTINATION_INSTANCE_X509_KEY_PATH"`

	DestinationSecretPath                string                               `envconfig:"APP_DESTINATION_SECRET_PATH"`
	RegionToDestinationCredentialsConfig map[string]DestinationInstanceConfig `envconfig:"-"`

	OAuthMode oauth.AuthMode `envconfig:"APP_DESTINATION_OAUTH_MODE,default=oauth-mtls"`
}

// TODO Code duplication
func (c *DestinationsConfig) getDestinationsSecret(path string) (string, error) {
	if path == "" {
		return "", errors.New("destinations secret path cannot be empty")
	}
	secret, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read destinations secret file")
	}

	return string(secret), nil
}

func (c *DestinationsConfig) MapInstanceConfigs() error {
	secretData, err := c.getDestinationsSecret(c.DestinationSecretPath)
	if err != nil {
		return errors.Wrapf(err, "while getting destinations secret")
	}

	if ok := gjson.Valid(secretData); !ok {
		return errors.New("failed to validate instance configs")
	}

	bindingsResult := gjson.Parse(secretData)
	bindingsMap := bindingsResult.Map()
	c.RegionToDestinationCredentialsConfig = make(map[string]DestinationInstanceConfig)
	for region, config := range bindingsMap {
		i := DestinationInstanceConfig{
			XsAppName:    gjson.Get(config.String(), c.InstanceXsAppNamePath).String(),
			ClientID:     gjson.Get(config.String(), c.InstanceClientIDPath).String(),
			ClientSecret: gjson.Get(config.String(), c.InstanceClientSecretPath).String(),
			URL:          gjson.Get(config.String(), c.InstanceURLPath).String(),
			AuthURL:      gjson.Get(config.String(), c.InstanceAuthURLPath).String(),
			Cert:         gjson.Get(config.String(), c.InstanceCertPath).String(),
			Key:          gjson.Get(config.String(), c.InstanceKeyPath).String(),
		}

		if err := i.validate(c.OAuthMode); err != nil {
			c.RegionToDestinationCredentialsConfig = nil
			return errors.Wrapf(err, "while validating instance for region: %q", region)
		}
		c.RegionToDestinationCredentialsConfig[region] = i
	}

	return nil
}

func (c *DestinationInstanceConfig) validate(oauthMode oauth.AuthMode) error {
	errorMessages := make([]string, 0)

	if c.XsAppName == "" {
		errorMessages = append(errorMessages, "XsAppName is missing")
	}

	if c.ClientID == "" {
		errorMessages = append(errorMessages, "Client ID is missing")
	}
	if c.AuthURL == "" {
		errorMessages = append(errorMessages, "Token URL is missing")
	}
	if c.URL == "" {
		errorMessages = append(errorMessages, "URL is missing")
	}

	switch oauthMode {
	case oauth.Standard:
		if c.ClientSecret == "" {
			errorMessages = append(errorMessages, "Client Secret is missing")
		}
	case oauth.Mtls:
		if c.Cert == "" {
			errorMessages = append(errorMessages, "Certificate is missing")
		}
		if c.Key == "" {
			errorMessages = append(errorMessages, "Key is missing")
		}
	}

	errorMsg := strings.Join(errorMessages, ", ")
	if errorMsg != "" {
		return errors.New(errorMsg)
	}

	return nil
}
