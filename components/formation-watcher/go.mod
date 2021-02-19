module github.com/kyma-incubator/compass/components/formation-watcher

go 1.15

require (
	code.cloudfoundry.org/lager v2.0.0+incompatible
	github.com/drewolson/testflight v1.0.0 // indirect
	github.com/google/uuid v1.2.0
	github.com/kyma-incubator/compass/components/director v0.0.0-20210219113652-f03b8599af6a
	github.com/lib/pq v1.9.0
	github.com/onrik/logrus v0.8.0
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/pivotal-cf/brokerapi v6.4.2+incompatible
	github.com/pkg/errors v0.9.1
	github.com/sergi/go-diff v1.1.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/oauth2 v0.0.0-20210113205817-d3ed898aa8a3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.18.6 // indirect
)

// replace github.com/kyma-incubator/compass/components/director v0.0.0-20210218134213-6ee23979ecdd => ../director
// replace github.com/kyma-incubator/compass/components/director v0.0.0-20210219113652-f03b8599af6a => github.com/kyma-incubator/compass/components/director v0.0.0-20210218134213-6ee23979ecdd
