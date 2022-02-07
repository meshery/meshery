module github.com/layer5io/meshery

go 1.16

replace (
	//golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/go-openapi/runtime => github.com/go-openapi/runtime v0.19.31
	github.com/go-openapi/strfmt => github.com/go-openapi/strfmt v0.20.2
	github.com/go-redis/redis/v8 => github.com/go-redis/redis/v8 v8.11.4
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	github.com/spf13/afero => github.com/spf13/afero v1.5.1 // Until viper bug is resolved #1161
	go.mongodb.org/mongo-driver v1.5.1 => github.com/mongodb/mongo-go-driver v1.5.1
	gonum.org/v1/plot v0.9.0 => github.com/gonum/plot v0.9.0
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)

require (
	fortio.org/fortio v1.20.0
	github.com/99designs/gqlgen v0.14.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/briandowns/spinner v1.18.0
	github.com/docker/docker v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/envoyproxy/go-control-plane v0.10.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/strfmt v0.19.5
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/gosimple/slug v1.12.0
	github.com/grafana-tools/sdk v0.0.0-20210630212345-db1192e93802
	github.com/jarcoal/httpmock v1.1.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/layer5io/gowrk2 v0.0.0-20191111234958-a4c9071c0f87
	github.com/layer5io/meshery-operator v0.5.2
	github.com/layer5io/meshkit v0.5.5
	github.com/layer5io/meshsync v0.1.30
	github.com/layer5io/nighthawk-go v0.5.0
	github.com/layer5io/service-mesh-performance v0.3.3
	github.com/lib/pq v1.10.4
	github.com/manifoldco/promptui v0.9.0
	github.com/nats-io/nats-server/v2 v2.2.0 // indirect; https://advisories.nats.io/CVE/CVE-2021-3127.txt
	github.com/nsf/termbox-go v1.1.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opencontainers/runc v1.0.0-rc95 // indirect; CVE-2021-30465
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.31.1
	github.com/qri-io/jsonschema v0.2.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.0
	github.com/vektah/gqlparser/v2 v2.2.0
	github.com/vmihailenco/taskq/v3 v3.2.7
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	gonum.org/v1/gonum v0.9.3
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.22.2
	k8s.io/api v0.22.4
	k8s.io/apiextensions-apiserver v0.22.3
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	sigs.k8s.io/controller-runtime v0.10.3
)
