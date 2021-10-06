module github.com/layer5io/meshery

go 1.16

replace (
	//golang.org/x/sys => golang.org/x/sys v0.0.0-20200826173525-f9321e4c35a6
	github.com/docker/docker => github.com/moby/moby v17.12.0-ce-rc1.0.20200618181300-9dc6525e6118+incompatible
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	// github.com/layer5io/meshkit => github.com/dhruv0000/meshkit v0.1.27-0.20210904132807-ad7b7db80267
	github.com/spf13/afero => github.com/spf13/afero v1.5.1 // Until viper bug is resolved #1161
	go.mongodb.org/mongo-driver v1.5.1 => github.com/mongodb/mongo-go-driver v1.5.1
	gonum.org/v1/plot v0.9.0 => github.com/gonum/plot v0.9.0
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)

require (
	fortio.org/fortio v1.17.0
	git.mills.io/prologic/bitcask v0.3.14
	github.com/99designs/gqlgen v0.14.0
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/aws/aws-sdk-go v1.36.30 // indirect
	github.com/briandowns/spinner v1.16.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/docker v1.13.1
	github.com/envoyproxy/go-control-plane v0.9.10-0.20210907150352-cf90f659a021
	github.com/ghodss/yaml v1.0.0
	github.com/go-openapi/runtime v0.19.31
	github.com/go-openapi/strfmt v0.20.2
	github.com/gofrs/uuid v3.4.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/google/go-github v17.0.0+incompatible // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/gosimple/slug v1.10.0
	github.com/grafana-tools/sdk v0.0.0-20210630212345-db1192e93802
	github.com/jarcoal/httpmock v1.0.8
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/klauspost/compress v1.11.0 // indirect
	github.com/layer5io/gowrk2 v0.0.0-20191111234958-a4c9071c0f87
	github.com/layer5io/meshery-operator v0.2.21
	github.com/layer5io/meshkit v0.2.29
	github.com/layer5io/meshsync v0.1.25
	github.com/layer5io/nighthawk-go v1.0.3
	github.com/layer5io/service-mesh-performance v0.3.3
	github.com/lib/pq v1.10.2
	github.com/lunixbochs/vtclean v1.0.0 // indirect
	github.com/manifoldco/promptui v0.8.0
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/nsf/termbox-go v1.1.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.11.0
	github.com/prometheus/common v0.31.1
	github.com/qri-io/jsonschema v0.2.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/tcnksm/go-latest v0.0.0-20170313132115-e3007ae9052e
	github.com/vektah/gqlparser/v2 v2.2.0
	github.com/vmihailenco/taskq/v3 v3.0.0-beta.9.0.20200519124923-d9823546b85a
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/exp v0.0.0-20200331195152-e8c3332aa8e5 // indirect
	golang.org/x/oauth2 v0.0.0-20210514164344-f6687ab2804c
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	gonum.org/v1/gonum v0.9.3
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.21.14
	k8s.io/api v0.18.12
	k8s.io/apimachinery v0.18.12
	k8s.io/client-go v0.18.12
)
