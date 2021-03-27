module github.com/layer5io/meshery

go 1.13

replace (
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)

require (
	fortio.org/fortio v1.14.1
	github.com/99designs/gqlgen v0.13.0
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/docker v1.13.1
	github.com/ghodss/yaml v1.0.0
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/gosimple/slug v1.9.0
	github.com/grafana-tools/sdk v0.0.0-20200411085644-f7626bef00b3
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/layer5io/gowrk2 v0.0.0-20191111234958-a4c9071c0f87
	github.com/layer5io/meshery-operator v0.2.11
	github.com/layer5io/meshkit v0.2.6
	github.com/layer5io/meshsync v0.1.17
	github.com/layer5io/nighthawk-go v0.1.5
	github.com/layer5io/service-mesh-performance v0.3.3
	github.com/lib/pq v1.7.0
	github.com/nats-io/nats.go v1.10.0
	github.com/olekukonko/tablewriter v0.0.2
	github.com/pkg/errors v0.9.1
	github.com/prologic/bitcask v0.3.9
	github.com/prometheus/client_golang v1.9.0
	github.com/prometheus/common v0.18.0
	github.com/qri-io/jsonschema v0.2.0
	github.com/sirupsen/logrus v1.8.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/tcnksm/go-latest v0.0.0-20170313132115-e3007ae9052e
	github.com/vektah/gqlparser/v2 v2.1.0
	github.com/vmihailenco/taskq/v3 v3.0.0-beta.9.0.20200519124923-d9823546b85a
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9 // indirect
	gonum.org/v1/gonum v0.8.2
	google.golang.org/grpc v1.36.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/gorm v1.20.10
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v0.18.12
)
