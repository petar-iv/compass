package config

import (
	"bytes"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type ControllerManagerSettings struct {
	RMConfig             string                      `yaml:"rmconfig" mapstructure:"rmconfig"`
	OAuth20              OAuth20Config               `yaml:"oauth20" mapstructure:"oauth20"`
	SkipSSLValidation    bool                        `yaml:"skip_ssl_validation" mapstructure:"skip_ssl_validation"`
	TenantFetcherURL     string                      `yaml:"tenant_fetcher_url" mapstructure:"tenant_fetcher_url"`
	Log                  LogSettings                 `yaml:"log" mapstructure:"log"`
	MaxConcurrentThreads ControllerConcurrentThreads `yaml:"max_concurrent_threads" mapstructure:"max_concurrent_threads"`
}

type OAuth20Config struct {
	ClientKey    string `yaml:"client_key" mapstructure:"client_key"`
	ClientSecret string `yaml:"client_secret" mapstructure:"client_secret"`
	TokenURL     string `yaml:"token_url" mapstructure:"token_url"`
}

type LogSettings struct {
	Enabled bool
	Level   string // debug|info|warn|error
}

type ControllerConcurrentThreads struct {
	OrganizationReconcilerThreads  int `yaml:"organization_controller" mapstructure:"organization_controller"`
	ResourceGroupReconcilerThreads int `yaml:"resource_group_controller" mapstructure:"resource_group_controller"`
	FolderReconcilerThreads        int `yaml:"folder_controller" mapstructure:"folder_controller"`
}

func NewControllerManagerSettings() (*ControllerManagerSettings, error) {
	config := new(ControllerManagerSettings)
	bs, err := yaml.Marshal(config)
	if err != nil {
		return nil, errors.New("unable to marshal config to YAML")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	if err = viper.ReadConfig(bytes.NewBuffer(bs)); err != nil {
		return nil, err
	}

	if err = viper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}
