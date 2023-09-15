module github.com/layer5io/meshery

go 1.21

replace (
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.9.3
	github.com/docker/docker => github.com/moby/moby v20.10.17+incompatible
	github.com/docker/libcompose => github.com/docker/libcompose v0.4.1-0.20190808084053-143e0f3f1ab9
	github.com/go-openapi/runtime => github.com/go-openapi/runtime v0.19.31
	github.com/go-openapi/strfmt => github.com/go-openapi/strfmt v0.20.2
	github.com/kubernetes/kompose => github.com/meshery/kompose v1.0.1 // do not upgrade until https://github.com/meshery/kompose/pull/3 has been merged
	github.com/kudobuilder/kuttl => github.com/layer5io/kuttl v0.4.1-0.20200806180306-b7e46afd657f
	gonum.org/v1/plot v0.9.0 => github.com/gonum/plot v0.9.0
	vbom.ml/util => github.com/fvbommel/util v0.0.0-20180919145318-efcd4e0f9787
)

require (
	cuelang.org/go v0.5.0
	fortio.org/fortio v1.58.0
	github.com/99designs/gqlgen v0.17.36
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/briandowns/spinner v1.23.0
	github.com/docker/cli v20.10.21+incompatible
	github.com/docker/docker v20.10.23+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/envoyproxy/go-control-plane v0.11.1
	github.com/ghodss/yaml v1.0.0
	github.com/go-errors/errors v1.4.2
	github.com/go-openapi/runtime v0.19.15
	github.com/go-openapi/strfmt v0.19.5
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.5.0
	github.com/gosimple/slug v1.12.0
	github.com/grafana-tools/sdk v0.0.0-20220919052116-6562121319fc
	github.com/jarcoal/httpmock v1.3.0
	github.com/jinzhu/copier v0.3.5
	github.com/layer5io/gowrk2 v0.6.1
	github.com/layer5io/meshery-operator v0.6.10
	github.com/layer5io/meshkit v0.6.65
	github.com/layer5io/meshsync v0.6.14
	github.com/layer5io/nighthawk-go v1.0.6
	github.com/layer5io/service-mesh-performance v0.6.1
	github.com/lib/pq v1.10.7
	github.com/manifoldco/promptui v0.9.0
	github.com/nsf/termbox-go v1.1.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo/v2 v2.11.0
	github.com/onsi/gomega v1.27.8
	github.com/pkg/browser v0.0.0-20210706143420-7d21f8c997e2
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.16.0
	github.com/prometheus/common v0.44.0
	github.com/qri-io/jsonschema v0.2.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.16.0
	github.com/stretchr/testify v1.8.4
	github.com/vektah/gqlparser/v2 v2.5.8
	github.com/vmihailenco/taskq/v3 v3.2.9
	golang.org/x/oauth2 v0.10.0
	golang.org/x/text v0.12.0
	gonum.org/v1/gonum v0.14.0
	google.golang.org/grpc v1.57.0
	google.golang.org/protobuf v1.31.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/gorm v1.25.4
	k8s.io/api v0.26.1
	k8s.io/apiextensions-apiserver v0.26.1
	k8s.io/apimachinery v0.26.1
	k8s.io/client-go v0.26.1
	sigs.k8s.io/controller-runtime v0.13.0
)

require (
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
)

require (
	fortio.org/dflag v1.5.3 // indirect
	fortio.org/log v1.10.0 // indirect
	fortio.org/sets v1.0.3 // indirect
	fortio.org/version v1.0.2 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.0 // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/Masterminds/squirrel v1.5.3 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/Microsoft/hcsshim v0.9.7 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20210428141323-04723f9f07d7 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bsm/redislock v0.7.2 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/capnm/sysinfo v0.0.0-20130621111458-5909a53897f3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/cncf/xds/go v0.0.0-20230607035331-e9ce68804cb4 // indirect
	github.com/cockroachdb/apd/v2 v2.0.2 // indirect
	github.com/containerd/cgroups v1.1.0 // indirect
	github.com/containerd/containerd v1.6.19 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/docker/go v1.5.1-1.0.20160303222718-d30aec9fd63c // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/docker/libcompose v0.4.0 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/envoyproxy/protoc-gen-validate v1.0.2 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fsouza/go-dockerclient v1.9.3 // indirect
	github.com/fvbommel/sortorder v1.0.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.4.0 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-gorp/gorp/v3 v3.0.2 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/analysis v0.19.10 // indirect
	github.com/go-openapi/errors v0.19.8 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.1 // indirect
	github.com/go-openapi/loads v0.19.5 // indirect
	github.com/go-openapi/spec v0.19.8 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-openapi/validate v0.19.10 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-redis/redis_rate/v9 v9.1.2 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/gosuri/uitable v0.0.4 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.3 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/huandu/xstrings v1.3.3 // indirect
	github.com/imdario/mergo v0.3.15 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.2.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/kubernetes/kompose v1.28.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/pkcs11 v1.1.1 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/sys/mount v0.3.3 // indirect
	github.com/moby/sys/mountinfo v0.6.2 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/mpvl/unique v0.0.0-20150818121801-cbe035fff7de // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nats-io/nats.go v1.22.1 // indirect
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/novln/docker-parser v1.0.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/open-policy-agent/opa v0.52.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc2 // indirect
	github.com/opencontainers/runc v1.1.5 // indirect
	github.com/openshift/api v3.9.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.4.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/qri-io/jsonpointer v0.1.1 // indirect
	github.com/rubenv/sql-migrate v1.2.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sergi/go-diff v1.3.1 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/spf13/afero v1.9.5 // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/subosito/gotenv v1.4.2 // indirect
	github.com/theupdateframework/notary v0.7.0 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	github.com/zmap/zcrypto v0.0.0-20230422215203-9a665e1e9968 // indirect
	github.com/zmap/zlint/v3 v3.4.1 // indirect
	go.mongodb.org/mongo-driver v1.5.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	golang.org/x/crypto v0.12.0 // indirect
	golang.org/x/exp v0.0.0-20230801115018-d63ba01acd4b // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/net v0.14.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/term v0.11.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.12.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230803162519-f966b187b2e5 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/postgres v1.4.6 // indirect
	gorm.io/driver/sqlite v1.4.4 // indirect
	helm.sh/helm/v3 v3.11.1 // indirect
	k8s.io/apiserver v0.26.1 // indirect
	k8s.io/cli-runtime v0.26.0 // indirect
	k8s.io/component-base v0.26.1 // indirect
	k8s.io/klog/v2 v2.90.0 // indirect
	k8s.io/kube-openapi v0.0.0-20230109183929-3758b55a6596 // indirect
	k8s.io/kubectl v0.26.0 // indirect
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448 // indirect
	oras.land/oras-go v1.2.2 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/kustomize/api v0.12.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)
