package config

import (
	"github.com/kyma-incubator/compass/components/director/pkg/oauth"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type DestinationsConfig struct {
	InstanceClientIDPath     string `envconfig:"APP_DESTINATION_INSTANCE_CLIENT_ID_PATH"`
	InstanceClientSecretPath string `envconfig:"APP_DESTINATION_INSTANCE_CLIENT_SECRET_PATH"`
	InstanceURLPath          string `envconfig:"APP_DESTINATION_INSTANCE_URL_PATH"`
	InstanceTokenURLPath     string `envconfig:"APP_DESTINATION_INSTANCE_TOKEN_URL_PATH"`
	InstanceCertPath         string `envconfig:"APP_DESTINATION_INSTANCE_X509_CERT_PATH"`
	InstanceKeyPath          string `envconfig:"APP_DESTINATION_INSTANCE_X509_KEY_PATH"`

	DestinationSecretPath                string                    `envconfig:"APP_DESTINATION_SECRET_PATH"`
	RegionToDestinationCredentialsConfig map[string]InstanceConfig `envconfig:"-"`

	OAuthMode oauth.AuthMode `envconfig:"APP_DESTINATION_OAUTH_MODE,default=oauth-mtls"`
}

func (c *DestinationsConfig) MapInstanceConfigs() error {
	secretData, err := ReadConfigFile(c.DestinationSecretPath)
	if err != nil {
		return errors.Wrapf(err, "while getting destinations secret")
	}

	bindingsMap, err := ParseConfigToJsonMap(secretData)
	if err != nil {
		return err
	}

	c.RegionToDestinationCredentialsConfig = make(map[string]InstanceConfig)
	for region, config := range bindingsMap {
		i := InstanceConfig{
			ClientID:     gjson.Get(config.String(), c.InstanceClientIDPath).String(),
			ClientSecret: gjson.Get(config.String(), c.InstanceClientSecretPath).String(),
			URL:          gjson.Get(config.String(), c.InstanceURLPath).String(),
			TokenURL:     gjson.Get(config.String(), c.InstanceTokenURLPath).String(),
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
