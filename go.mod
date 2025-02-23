module github.com/layer5io/meshery

go 1.23.4

// replace github.com/layer5io/meshkit v0.8.15 => ../meshkit

replace github.com/vektah/gqlparser/v2 => github.com/vektah/gqlparser/v2 v2.5.11

require (
	cuelang.org/go v0.11.2
	fortio.org/fortio v1.66.2
	github.com/99designs/gqlgen v0.17.45
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/briandowns/spinner v1.23.1
	github.com/docker/cli v27.5.1+incompatible
	github.com/docker/docker v27.5.1+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/eiannone/keyboard v0.0.0-20220611211555-0d226195f203
	github.com/envoyproxy/go-control-plane/envoy v1.32.3
	github.com/fatih/color v1.18.0
	github.com/fsnotify/fsnotify v1.8.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.28.0
	github.com/go-openapi/strfmt v0.23.0
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/gorilla/websocket v1.5.3
	github.com/gosimple/slug v1.15.0
	github.com/grafana-tools/sdk v0.0.0-20220919052116-6562121319fc
	github.com/jarcoal/httpmock v1.3.1
	github.com/jinzhu/copier v0.4.0
	github.com/layer5io/gowrk2 v0.6.1
	github.com/layer5io/meshery-operator v0.8.1
	github.com/layer5io/meshkit v0.8.17
	github.com/layer5io/meshsync v0.8.5
	github.com/layer5io/nighthawk-go v1.0.3
	github.com/layer5io/service-mesh-performance v0.6.1
	github.com/lib/pq v1.10.9
	github.com/manifoldco/promptui v0.9.0
	github.com/meshery/schemas v0.7.45
	github.com/nsf/termbox-go v1.1.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/common v0.62.0
	github.com/qri-io/jsonschema v0.2.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.19.0
	github.com/stretchr/testify v1.10.0
	github.com/vektah/gqlparser/v2 v2.5.21
	github.com/vmihailenco/taskq/v3 v3.2.9
	golang.org/x/oauth2 v0.25.0
	golang.org/x/sync v0.11.0
	golang.org/x/text v0.22.0
	gonum.org/v1/gonum v0.15.1
	google.golang.org/api v0.218.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.4
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/gorm v1.25.12
	k8s.io/api v0.32.1
	k8s.io/apiextensions-apiserver v0.32.1
	k8s.io/apimachinery v0.32.1
	k8s.io/client-go v0.32.1
	sigs.k8s.io/controller-runtime v0.20.1
)

require (
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/theupdateframework/notary v0.7.0 // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	golang.org/x/lint v0.0.0-20241112194109-818c5a804067 // indirect
	honnef.co/go/tools v0.6.0 // indirect
)
