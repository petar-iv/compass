package main

import (
	"time"

	"github.com/kyma-incubator/compass/components/director/pkg/log"
)

type Config struct {
	Log log.Config

	HydratorAddress string `envconfig:"default=127.0.0.1:8080"`

	ServerTimeout time.Duration `envconfig:"default=100s"`

	CSRSubject struct {
		Country            string `envconfig:"default=PL"`
		Organization       string `envconfig:"default=Org"`
		OrganizationalUnit string `envconfig:"default=OrgUnit"`
		Locality           string `envconfig:"default=Locality"`
		Province           string `envconfig:"default=State"`
	}
	ExternalIssuerSubject struct {
		Country                   string `envconfig:"default=DE"`
		Organization              string `envconfig:"default=Org"`
		OrganizationalUnitPattern string `envconfig:"default=OrgUnit"`
	}
	CertificateDataHeader   string `envconfig:"default=Certificate-Data"`
	RevocationConfigMapName string `envconfig:"default=compass-system/revocations-Config"`

	KubernetesClient struct {
		PollInteval time.Duration `envconfig:"default=2s"`
		PollTimeout time.Duration `envconfig:"default=1m"`
		Timeout     time.Duration `envconfig:"default=95s"`
	}

	SubjectConsumerMappingConfig string `envconfig:"default=[]"`
}
