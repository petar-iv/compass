package config

import "github.com/spf13/viper"

type ControllerManagerSettings struct {
	Webhooks  WebhooksSettings `json:"webhooks" mapstructure:"webhooks"`
	APIConfig string           `json:"apiconfig" mapstructure:"apiconfig"`
	Log       LogSettings      `json:"log" mapstructure:"log"`
}

type WebhooksSettings struct {
	Admission WebServerConfiguration `json:"admission" mapstructure:"admission"`
	X509      X509Settings           `json:"x509" mapstructure:"x509"`
}

type WebServerConfiguration struct {
	Port       int  `json:"port" mapstructure:"port"`
	TLSEnabled bool `json:"tls" mapstructure:"tls"`
}

type X509Settings struct {
	Path        string `json:"path" mapstructure:"path"`
	Certificate string `json:"certificate" mapstructure:"certificate"`
	Key         string `json:"key" mapstructure:"key"`
}

type LogSettings struct {
	Enabled bool
	Level   string
}

func NewControllerManagerSettings(configPath string) (*ControllerManagerSettings, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("application")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	viper.AutomaticEnv()

	var settings *ControllerManagerSettings
	if err := viper.Unmarshal(&settings); err != nil {
		return nil, err
	}

	return settings, nil
}
