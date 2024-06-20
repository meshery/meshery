# Meshery Proposal for CNCF Incubation

## About Meshery

[Meshery](https://meshery.io/) is a self-service engineering platform, designed to manage the complete lifecycle of modern infrastructure.

The project encompasses tools that can be utilized through a command-line interface, integrated into clusters, or incorporated into various other tools. It enables users to assess cloud native configurations during development, integration, and post-deployment phases.

## Meshery in the CNCF Sandbox

Since its entry into the CNCF Sandbox on June 22, 2021, Meshery has undergone significant enhancements and expansions. 

* Supports every CNCF project and public cloud services that have native resource / service integration with Kubernetes.
* Introduction of advanced load generation and performance characterization features.
* Enhanced support for multi-cluster and hybrid cloud environments.
* Meshery has evolved into a full-blown self-service engineering platform with a robust extensibility framework for easy integration of new extensions whether they are backend adapters, frontend plugins, or separately hosted providers of additional functionality (e.g. identity providers).

### Community & Growth

Since joining the CNCF Sandbox, Meshery has experienced notable community growth and engagement:

* [Is the 10th fastest growing project in the CNCF](https://www.cncf.io/blog/2023/10/27/october-2023-where-we-are-with-velocity-of-cncf-lf-and-top-30-open-source-projects/).
* Accumulated a total of 4,000 new GitHub stars
* Closed 9,000 [Pull Requests (PRs)](https://github.com/pulls?q=is%3Aissue+created%3A%3E%3D2021-06-21+org%3Ameshery), showcasing an active contributor community.
* Addressed and resolved 3,000 [issues](https://github.com/issues?q=is%3Aissue+created%3A%3E%3D2021-06-21+org%3Ameshery+), ensuring ongoing project improvement.
* Participated heavily in mentoring initiatives and is the #1 most popular LFX internship.

### Neutrality

Meshery, committed to maintaining vendor neutrality, ensuring a diverse collection of maintainers representing 10 different organizations. As an extensible platform, Meshery goes to great lengths to support third-party plugins, elevating them within its ecosystem of integrations. Meshery places emphasis on fostering an open and inclusive community of contributors, maintainers, and integrators. 

## Incubation stage requirements

### Used successfully in production by at least three independent direct adopters

Meshery's non-exhaustive, public list of [adaptors](https://github.com/meshery/meshery/blob/master/ADOPTERS.md) is available for reference. With thousands of users already leveraging its features, Meshery has proven its production-readiness and garnered positive feedback from adopters and the broader community.

### Have a healthy number of committers
Meshery boasts a diverse and active community of maintainers, with [15 maintainers](https://github.com/meshery/meshery/blob/main/MAINTAINERS.md) overseeing the ~30 repositories included in the `meshery` GitHub organization.  and [extension modules](https://github.com/meshery/meshery/blob/main/EXTENSIONS.md). 

### Demonstrates a substantial ongoing flow of commits and merged contributions

Meshery is the 10th fastest growing CNCF project.

Over the past year, over [300 committers have had PRs merged](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=merged_prs&var-repogroup_name=All&var-country_name=All&var-companies=All), and [800 contributors actively participating](https://meshery.devstats.cncf.io/d/66/developer-activity-counts-by-companies?orgId=1&var-period_name=Last%20year&var-metric=contributions&var-repogroup_name=All&var-country_name=All&var-companies=All).

### Clear versioning scheme

Meshery follows a clear and [well-documented](https://docs.meshery.io/project/contributing/build-and-release) build and release process. All components adhere to the principles of [semantic versioning](https://semver.org/). The project releases page documents [Meshery Server, UI, and CLI releases](https://docs.meshery.io/project/releases) with Meshery Adapters, Meshery Operator and custom controllers, centralized API schemas, and so on all releasing independently. See each repo for a release history.

### Clearly documented security processes explaining how to report security issues to the project, and describing how the project provides updated releases or patches to resolve security vulnerabilities

Meshery's [security reporting process](https://docs.meshery.io/project/security-vulnerabilities) is well-documented and has [previous and current CVEs](https://docs.meshery.io/project/security-vulnerabilities) published along with the reporting process and expectation setting of how reports are handled. This process is followed successfully with GitHub engineers being the most recent to report vulneraabilities and a Meshery maintainer (from Intel) to provide patches. Mainatiners strive to acknowledge reports within 24 hours of being received and offer an analysis and plan for a full response withing 10 days.

## Alignment with the CNCF Mission

Meshery aligns closely with the CNCF mission by contributing to the security and efficiency of cloud-native computing. Focused solely on cloud native management, Meshery integrates deeply and seamlessly with _many_ CNCF projects. The full list of project integrations can be seen on https://meshery.io/integrations. Meshery's commitment to interoperability and collaboration is one of a number of reasons that it is a valuable asset to the CNCF ecosystem.

## Future plans

Looking ahead to 2024, Meshery's roadmap includes the introduction of advanced features in collaborative workflows of developer-defined infrasturcture and the engineering teams that operate this infrastructure. These additions will be crafted through a combination of new AI-centric developments, integration of workflow engine, and potential launch as the [CNCF's official project playground](https://docs.google.com/document/d/1Cr0MxlOxWq70d-BUisfKXF_VPGss0hK8XqRxy4xm4YE/edit#heading=h.58lqw93jp55u). Meshery remains dedicated to evolving as a key player in the cloud native landscape, responding to user needs and industry trends.


## Bill of Materials and Licenses
- name: cloud.google.com/go/compute
  version: v1.23.3
  license: Apache-2.0
- name: cloud.google.com/go/compute/metadata
  version: v0.2.3
  license: Apache-2.0
- name: cuelang.org/go
  version: v0.6.0
  license: Apache-2.0
- name: dario.cat/mergo
  version: v1.0.0
  license: BSD-3-Clause
- name: fortio.org/dflag
  version: v1.7.0
  license: Apache-2.0
- name: fortio.org/fortio
  version: v1.63.2
  license: Apache-2.0
- name: fortio.org/log
  version: v1.12.0
  license: Apache-2.0
- name: fortio.org/sets
  version: v1.0.3
  license: Apache-2.0
- name: fortio.org/struct2env
  version: v0.4.0
  license: Apache-2.0
- name: fortio.org/version
  version: v1.0.3
  license: Apache-2.0
- name: github.com/99designs/gqlgen
  version: v0.17.42
  license: MIT
- name: github.com/AdaLogics/go-fuzz-headers
  version: v0.0.0-20230811130428-ced1acdcaa24
  license: Apache-2.0
- name: github.com/Azure/azure-sdk-for-go/sdk/azcore
  version: v1.9.0
  license: MIT
- name: github.com/Azure/azure-sdk-for-go/sdk/azidentity
  version: v1.4.0
  license: MIT
- name: github.com/Azure/azure-sdk-for-go/sdk/internal
  version: v1.5.0
  license: MIT
- name: github.com/Azure/go-ansiterm
  version: v0.0.0-20230124172434-306776ec8161
  license: MIT
- name: github.com/AzureAD/microsoft-authentication-library-for-go
  version: v1.1.1
  license: MIT
- name: github.com/BurntSushi/toml
  version: v1.3.2
  license: MIT
- name: github.com/MakeNowJust/heredoc
  version: v1.0.0
  license: MIT
- name: github.com/Masterminds/goutils
  version: v1.1.1
  license: Apache-2.0
- name: github.com/Masterminds/semver/v3
  version: v3.2.1
  license: MIT
- name: github.com/Masterminds/sprig/v3
  version: v3.2.3
  license: MIT
- name: github.com/Masterminds/squirrel
  version: v1.5.4
  license: MIT
- name: github.com/Microsoft/go-winio
  version: v0.6.1
  license: MIT
- name: github.com/Microsoft/hcsshim
  version: v0.11.4
  license: MIT
- name: github.com/OneOfOne/xxhash
  version: v1.2.8
  license: Apache-2.0
- name: github.com/ProtonMail/go-crypto
  version: v0.0.0-20230828082145-3c4c8a2d2371
  license: BSD-3-Clause
- name: github.com/agnivade/levenshtein
  version: v1.1.1
  license: MIT
- name: github.com/asaskevich/govalidator
  version: v0.0.0-20230301143203-a9d515a09cc2
  license: MIT
- name: github.com/aws/aws-sdk-go-v2
  version: v1.24.0
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/config
  version: v1.26.1
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/credentials
  version: v1.16.12
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/feature/ec2/imds
  version: v1.14.10
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/internal/configsources
  version: v1.2.9
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/internal/endpoints/v2
  version: v2.5.9
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/internal/ini
  version: v1.7.2
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/service/ecr
  version: v1.24.5
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding
  version: v1.10.4
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/service/internal/presigned-url
  version: v1.10.9
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/service/sso
  version: v1.18.5
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/service/ssooidc
  version: v1.21.5
  license: Apache-2.0
- name: github.com/aws/aws-sdk-go-v2/service/sts
  version: v1.26.5
  license: Apache-2.0
- name: github.com/aws/smithy-go
  version: v1.19.0
  license: Apache-2.0
- name: github.com/beorn7/perks
  version: v1.0.1
  license: MIT
- name: github.com/briandowns/spinner
  version: v1.23.0
  license: Apache-2.0
- name: github.com/bsm/redislock
  version: v0.7.2
  license: Apache-2.0
- name: github.com/buger/jsonparser
  version: v1.1.1
  license: MIT
- name: github.com/capnm/sysinfo
  version: v0.0.0-20130621111458-5909a53897f3
  license: BSD-2-Clause
- name: github.com/cespare/xxhash/v2
  version: v2.2.0
  license: MIT
- name: github.com/chai2010/gettext-go
  version: v1.0.2
  license: BSD-3-Clause
- name: github.com/chzyer/readline
  version: v1.5.1
  license: MIT
- name: github.com/cloudflare/circl
  version: v1.3.7
  license: BSD-3-Clause
- name: github.com/cncf/xds/go
  version: v0.0.0-20230607035331-e9ce68804cb4
  license: Apache-2.0
- name: github.com/cockroachdb/apd/v3
  version: v3.2.1
  license: Apache-2.0
- name: github.com/compose-spec/compose-go
  version: v1.19.0
  license: Apache-2.0
- name: github.com/containerd/containerd
  version: v1.7.11
  license: Apache-2.0
- name: github.com/containerd/log
  version: v0.1.0
  license: Apache-2.0
- name: github.com/containerd/stargz-snapshotter/estargz
  version: v0.14.3
  license: Apache-2.0
- name: github.com/cyphar/filepath-securejoin
  version: v0.2.4
  license: BSD-3-Clause
- name: github.com/deckarep/golang-set
  version: v1.8.0
  license: MIT
- name: github.com/dgryski/go-rendezvous
  version: v0.0.0-20200823014737-9f7001d12a5f
  license: MIT
- name: github.com/distribution/reference
  version: v0.5.0
  license: Apache-2.0
- name: github.com/docker/cli
  version: v23.0.6+incompatible
  license: Apache-2.0
- name: github.com/docker/distribution
  version: v2.8.2+incompatible
  license: Apache-2.0
- name: github.com/docker/docker
  version: v23.0.6+incompatible
  license: Apache-2.0
- name: github.com/docker/docker-credential-helpers
  version: v0.8.0
  license: MIT
- name: github.com/docker/go
  version: v1.5.1-1.0.20160303222718-d30aec9fd63c
  license: BSD-3-Clause
- name: github.com/docker/go-connections
  version: v0.5.0
  license: Apache-2.0
- name: github.com/docker/go-metrics
  version: v0.0.1
  license: CC-BY-SA-4.0
- name: github.com/docker/go-units
  version: v0.5.0
  license: Apache-2.0
- name: github.com/eiannone/keyboard
  version: v0.0.0-20220611211555-0d226195f203
  license: MIT
- name: github.com/emicklei/go-restful/v3
  version: v3.11.0
  license: MIT
- name: github.com/emirpasic/gods
  version: v1.18.1
  license: BSD-2-Clause
- name: github.com/envoyproxy/go-control-plane
  version: v0.11.1
  license: Apache-2.0
- name: github.com/envoyproxy/protoc-gen-validate
  version: v1.0.2
  license: Apache-2.0
- name: github.com/evanphx/json-patch
  version: v5.7.0+incompatible
  license: BSD-3-Clause
- name: github.com/exponent-io/jsonpath
  version: v0.0.0-20210407135951-1de76d718b3f
  license: MIT
- name: github.com/fatih/camelcase
  version: v1.0.0
  license: MIT
- name: github.com/fatih/color
  version: v1.15.0
  license: MIT
- name: github.com/fatih/structs
  version: v1.1.0
  license: MIT
- name: github.com/felixge/httpsnoop
  version: v1.0.4
  license: MIT
- name: github.com/fluxcd/pkg/oci
  version: v0.34.0
  license: Apache-2.0
- name: github.com/fluxcd/pkg/sourceignore
  version: v0.4.0
  license: Apache-2.0
- name: github.com/fluxcd/pkg/tar
  version: v0.4.0
  license: BSD-3-Clause
- name: github.com/fluxcd/pkg/version
  version: v0.2.2
  license: Apache-2.0
- name: github.com/fsnotify/fsnotify
  version: v1.7.0
  license: BSD-3-Clause
- name: github.com/fsouza/go-dockerclient
  version: v1.9.8
  license: BSD-2-Clause
- name: github.com/fvbommel/sortorder
  version: v1.1.0
  license: MIT
- name: github.com/ghodss/yaml
  version: v1.0.0
  license: null
- name: github.com/go-errors/errors
  version: v1.5.1
  license: MIT
- name: github.com/go-git/gcfg
  version: v1.5.1-0.20230307220236-3a3c6141e376
  license: BSD-3-Clause
- name: github.com/go-git/go-billy/v5
  version: v5.5.0
  license: Apache-2.0
- name: github.com/go-git/go-git/v5
  version: v5.11.0
  license: Apache-2.0
- name: github.com/go-gorp/gorp/v3
  version: v3.1.0
  license: MIT
- name: github.com/go-ini/ini
  version: v1.67.0
  license: Apache-2.0
- name: github.com/go-logr/logr
  version: v1.3.0
  license: Apache-2.0
- name: github.com/go-logr/stdr
  version: v1.2.2
  license: Apache-2.0
- name: github.com/go-openapi/analysis
  version: v0.21.4
  license: Apache-2.0
- name: github.com/go-openapi/errors
  version: v0.20.4
  license: Apache-2.0
- name: github.com/go-openapi/jsonpointer
  version: v0.20.0
  license: Apache-2.0
- name: github.com/go-openapi/jsonreference
  version: v0.20.2
  license: Apache-2.0
- name: github.com/go-openapi/loads
  version: v0.21.2
  license: Apache-2.0
- name: github.com/go-openapi/runtime
  version: v0.26.2
  license: Apache-2.0
- name: github.com/go-openapi/spec
  version: v0.20.11
  license: Apache-2.0
- name: github.com/go-openapi/strfmt
  version: v0.21.8
  license: Apache-2.0
- name: github.com/go-openapi/swag
  version: v0.22.4
  license: Apache-2.0
- name: github.com/go-openapi/validate
  version: v0.22.3
  license: Apache-2.0
- name: github.com/go-redis/redis/v8
  version: v8.11.5
  license: BSD-2-Clause
- name: github.com/go-redis/redis_rate/v9
  version: v9.1.2
  license: BSD-2-Clause
- name: github.com/gobwas/glob
  version: v0.2.3
  license: MIT
- name: github.com/gocarina/gocsv
  version: v0.0.0-20231116093920-b87c2d0e983a
  license: MIT
- name: github.com/gofrs/uuid
  version: v4.4.0+incompatible
  license: MIT
- name: github.com/gogo/protobuf
  version: v1.3.2
  license: 0BSD
- name: github.com/golang-jwt/jwt
  version: v3.2.2+incompatible
  license: MIT
- name: github.com/golang-jwt/jwt/v5
  version: v5.0.0
  license: MIT
- name: github.com/golang/groupcache
  version: v0.0.0-20210331224755-41bb18bfe9da
  license: Apache-2.0
- name: github.com/golang/protobuf
  version: v1.5.3
  license: BSD-3-Clause
- name: github.com/google/btree
  version: v1.1.2
  license: Apache-2.0
- name: github.com/google/gnostic-models
  version: v0.6.8
  license: Apache-2.0
- name: github.com/google/go-cmp
  version: v0.6.0
  license: BSD-3-Clause
- name: github.com/google/go-containerregistry
  version: v0.17.0
  license: Apache-2.0
- name: github.com/google/gofuzz
  version: v1.2.0
  license: Apache-2.0
- name: github.com/google/s2a-go
  version: v0.1.7
  license: Apache-2.0
- name: github.com/google/shlex
  version: v0.0.0-20191202100458-e7afc7fbc510
  license: Apache-2.0
- name: github.com/google/uuid
  version: v1.5.0
  license: BSD-3-Clause
- name: github.com/googleapis/enterprise-certificate-proxy
  version: v0.3.2
  license: Apache-2.0
- name: github.com/googleapis/gax-go/v2
  version: v2.12.0
  license: BSD-3-Clause
- name: github.com/gorilla/mux
  version: v1.8.1
  license: BSD-3-Clause
- name: github.com/gorilla/websocket
  version: v1.5.0
  license: BSD-2-Clause
- name: github.com/gosimple/slug
  version: v1.13.1
  license: MPL-2.0
- name: github.com/gosimple/unidecode
  version: v1.0.1
  license: Apache-2.0
- name: github.com/gosuri/uitable
  version: v0.0.4
  license: MIT
- name: github.com/grafana-tools/sdk
  version: v0.0.0-20220919052116-6562121319fc
  license: Apache-2.0
- name: github.com/gregjones/httpcache
  version: v0.0.0-20190611155906-901d90724c79
  license: MIT
- name: github.com/hashicorp/errwrap
  version: v1.1.0
  license: MPL-2.0
- name: github.com/hashicorp/go-multierror
  version: v1.1.1
  license: MPL-2.0
- name: github.com/hashicorp/golang-lru
  version: v0.5.4
  license: MPL-2.0
- name: github.com/hashicorp/golang-lru/v2
  version: v2.0.7
  license: MPL-2.0
- name: github.com/hashicorp/hcl
  version: v1.0.0
  license: MPL-2.0
- name: github.com/huandu/xstrings
  version: v1.4.0
  license: MIT
- name: github.com/imdario/mergo
  version: v0.3.16
  license: BSD-3-Clause
- name: github.com/inconshreveable/mousetrap
  version: v1.1.0
  license: Apache-2.0
- name: github.com/jackc/pgpassfile
  version: v1.0.0
  license: MIT
- name: github.com/jackc/pgservicefile
  version: v0.0.0-20221227161230-091c0ba34f0a
  license: MIT
- name: github.com/jackc/pgx/v5
  version: v5.4.3
  license: MIT
- name: github.com/jarcoal/httpmock
  version: v1.3.1
  license: MIT
- name: github.com/jbenet/go-context
  version: v0.0.0-20150711004518-d14ea06fba99
  license: MIT
- name: github.com/jinzhu/copier
  version: v0.4.0
  license: MIT
- name: github.com/jinzhu/inflection
  version: v1.0.0
  license: MIT
- name: github.com/jinzhu/now
  version: v1.1.5
  license: MIT
- name: github.com/jmespath/go-jmespath
  version: v0.4.0
  license: Apache-2.0
- name: github.com/jmoiron/sqlx
  version: v1.3.5
  license: MIT
- name: github.com/joho/godotenv
  version: v1.5.1
  license: MIT
- name: github.com/josharian/intern
  version: v1.0.0
  license: MIT
- name: github.com/json-iterator/go
  version: v1.1.12
  license: MIT
- name: github.com/kevinburke/ssh_config
  version: v1.2.0
  license: MIT
- name: github.com/klauspost/compress
  version: v1.17.0
  license: Apache-2.0
- name: github.com/meshery/kompose
  version: v1.26.2-0.20230425025309-3bb778d54007
  license: Apache-2.0
- name: github.com/kylelemons/godebug
  version: v1.1.0
  license: Apache-2.0
- name: github.com/lann/builder
  version: v0.0.0-20180802200727-47ae307949d0
  license: MIT
- name: github.com/lann/ps
  version: v0.0.0-20150810152359-62de8c46ede0
  license: MIT
- name: github.com/layer5io/gowrk2
  version: v0.6.1
  license: Apache-2.0
- name: github.com/layer5io/meshery-operator
  version: v0.7.0
  license: Apache-2.0
- name: github.com/layer5io/meshkit
  version: v0.7.14
  license: Apache-2.0
- name: github.com/layer5io/meshsync
  version: v0.6.24
  license: Apache-2.0
- name: github.com/layer5io/nighthawk-go
  version: v1.0.3
  license: Apache-2.0
- name: github.com/layer5io/service-mesh-performance
  version: v0.6.1
  license: Apache-2.0
- name: github.com/lib/pq
  version: v1.10.9
  license: MIT
- name: github.com/liggitt/tabwriter
  version: v0.0.0-20181228230101-89fcab3d43de
  license: BSD-3-Clause
- name: github.com/magiconair/properties
  version: v1.8.7
  license: BSD-2-Clause
- name: github.com/mailru/easyjson
  version: v0.7.7
  license: MIT
- name: github.com/manifoldco/promptui
  version: v0.9.0
  license: BSD-3-Clause
- name: github.com/mattn/go-colorable
  version: v0.1.13
  license: MIT
- name: github.com/mattn/go-isatty
  version: v0.0.19
  license: MIT
- name: github.com/mattn/go-runewidth
  version: v0.0.15
  license: MIT
- name: github.com/mattn/go-shellwords
  version: v1.0.12
  license: MIT
- name: github.com/mattn/go-sqlite3
  version: v1.14.17
  license: MIT
- name: github.com/matttproud/golang_protobuf_extensions/v2
  version: v2.0.0
  license: Apache-2.0
- name: github.com/miekg/pkcs11
  version: v1.1.1
  license: BSD-3-Clause
- name: github.com/mitchellh/copystructure
  version: v1.2.0
  license: MIT
- name: github.com/mitchellh/go-homedir
  version: v1.1.0
  license: MIT
- name: github.com/mitchellh/go-wordwrap
  version: v1.0.1
  license: MIT
- name: github.com/mitchellh/mapstructure
  version: v1.5.0
  license: MIT
- name: github.com/mitchellh/reflectwalk
  version: v1.0.2
  license: MIT
- name: github.com/moby/locker
  version: v1.0.1
  license: Apache-2.0
- name: github.com/moby/patternmatcher
  version: v0.6.0
  license: Apache-2.0
- name: github.com/moby/spdystream
  version: v0.2.0
  license: Apache-2.0
- name: github.com/moby/sys/sequential
  version: v0.5.0
  license: Apache-2.0
- name: github.com/moby/term
  version: v0.5.0
  license: Apache-2.0
- name: github.com/modern-go/concurrent
  version: v0.0.0-20180306012644-bacd9c7ef1dd
  license: Apache-2.0
- name: github.com/modern-go/reflect2
  version: v1.0.2
  license: Apache-2.0
- name: github.com/monochromegane/go-gitignore
  version: v0.0.0-20200626010858-205db1a8cc00
  license: MIT
- name: github.com/morikuni/aec
  version: v1.0.0
  license: MIT
- name: github.com/mpvl/unique
  version: v0.0.0-20150818121801-cbe035fff7de
  license: MIT
- name: github.com/munnerz/goautoneg
  version: v0.0.0-20191010083416-a7dc8b61c822
  license: BSD-3-Clause
- name: github.com/nats-io/nats.go
  version: v1.31.0
  license: Apache-2.0
- name: github.com/nats-io/nkeys
  version: v0.4.6
  license: Apache-2.0
- name: github.com/nats-io/nuid
  version: v1.0.1
  license: Apache-2.0
- name: github.com/novln/docker-parser
  version: v1.0.0
  license: Apache-2.0
- name: github.com/nsf/termbox-go
  version: v1.1.1
  license: MIT
- name: github.com/oklog/ulid
  version: v1.3.1
  license: Apache-2.0
- name: github.com/olekukonko/tablewriter
  version: v0.0.5
  license: MIT
- name: github.com/open-policy-agent/opa
  version: v0.57.1
  license: Apache-2.0
- name: github.com/opencontainers/go-digest
  version: v1.0.0
  license: CC-BY-SA-4.0
- name: github.com/opencontainers/image-spec
  version: v1.1.0-rc5
  license: Apache-2.0
- name: github.com/opencontainers/runc
  version: v1.1.12
  license: Apache-2.0
- name: github.com/openshift/api
  version: v3.9.0+incompatible
  license: Apache-2.0
- name: github.com/pelletier/go-toml/v2
  version: v2.1.0
  license: MIT
- name: github.com/peterbourgon/diskv
  version: v2.0.1+incompatible
  license: MIT
- name: github.com/pjbgf/sha1cd
  version: v0.3.0
  license: Apache-2.0
- name: github.com/pkg/browser
  version: v0.0.0-20210911075715-681adbf594b8
  license: BSD-2-Clause
- name: github.com/pkg/errors
  version: v0.9.1
  license: BSD-2-Clause
- name: github.com/prometheus/client_golang
  version: v1.18.0
  license: Apache-2.0
- name: github.com/prometheus/client_model
  version: v0.5.0
  license: Apache-2.0
- name: github.com/prometheus/common
  version: v0.45.0
  license: Apache-2.0
- name: github.com/prometheus/procfs
  version: v0.12.0
  license: Apache-2.0
- name: github.com/qri-io/jsonpointer
  version: v0.1.1
  license: MIT
- name: github.com/qri-io/jsonschema
  version: v0.2.1
  license: MIT
- name: github.com/rcrowley/go-metrics
  version: v0.0.0-20201227073835-cf1acfcdf475
  license: null
- name: github.com/rivo/uniseg
  version: v0.4.4
  license: MIT
- name: github.com/rubenv/sql-migrate
  version: v1.5.2
  license: MIT
- name: github.com/russross/blackfriday/v2
  version: v2.1.0
  license: BSD-2-Clause
- name: github.com/sagikazarmark/locafero
  version: v0.4.0
  license: MIT
- name: github.com/sagikazarmark/slog-shim
  version: v0.1.0
  license: BSD-3-Clause
- name: github.com/sergi/go-diff
  version: v1.3.1
  license: MIT
- name: github.com/shopspring/decimal
  version: v1.3.1
  license: MIT
- name: github.com/sirupsen/logrus
  version: v1.9.3
  license: MIT
- name: github.com/skeema/knownhosts
  version: v1.2.1
  license: Apache-2.0
- name: github.com/sosodev/duration
  version: v1.2.0
  license: MIT
- name: github.com/sourcegraph/conc
  version: v0.3.0
  license: MIT
- name: github.com/spf13/afero
  version: v1.11.0
  license: Apache-2.0
- name: github.com/spf13/cast
  version: v1.6.0
  license: MIT
- name: github.com/spf13/cobra
  version: v1.8.0
  license: Apache-2.0
- name: github.com/spf13/pflag
  version: v1.0.5
  license: BSD-3-Clause
- name: github.com/spf13/viper
  version: v1.17.0
  license: MIT
- name: github.com/subosito/gotenv
  version: v1.6.0
  license: MIT
- name: github.com/tchap/go-patricia/v2
  version: v2.3.1
  license: MIT
- name: github.com/theupdateframework/notary
  version: v0.7.0
  license: Apache-2.0
- name: github.com/vbatts/tar-split
  version: v0.11.3
  license: BSD-3-Clause
- name: github.com/vektah/gqlparser/v2
  version: v2.5.10
  license: MIT
- name: github.com/vmihailenco/msgpack/v5
  version: v5.3.5
  license: BSD-2-Clause
- name: github.com/vmihailenco/tagparser/v2
  version: v2.0.0
  license: BSD-2-Clause
- name: github.com/vmihailenco/taskq/v3
  version: v3.2.9
  license: BSD-2-Clause
- name: github.com/xanzy/ssh-agent
  version: v0.3.3
  license: Apache-2.0
- name: github.com/xeipuuv/gojsonpointer
  version: v0.0.0-20190905194746-02993c407bfb
  license: Apache-2.0
- name: github.com/xeipuuv/gojsonreference
  version: v0.0.0-20180127040603-bd5ef7bd5415
  license: Apache-2.0
- name: github.com/xeipuuv/gojsonschema
  version: v1.2.0
  license: Apache-2.0
- name: github.com/xlab/treeprint
  version: v1.2.0
  license: MIT
- name: github.com/yashtewari/glob-intersection
  version: v0.2.0
  license: Apache-2.0
- name: go.mongodb.org/mongo-driver
  version: v1.13.1
  license: Apache-2.0
- name: go.opencensus.io
  version: v0.24.0
  license: Apache-2.0
- name: go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp
  version: v0.46.1
  license: Apache-2.0
- name: go.opentelemetry.io/otel
  version: v1.21.0
  license: Apache-2.0
- name: go.opentelemetry.io/otel/metric
  version: v1.21.0
  license: Apache-2.0
- name: go.opentelemetry.io/otel/sdk
  version: v1.21.0
  license: Apache-2.0
- name: go.opentelemetry.io/otel/trace
  version: v1.21.0
  license: Apache-2.0
- name: go.starlark.net
  version: v0.0.0-20230814145427-12f4cb8177e4
  license: BSD-3-Clause
- name: go.uber.org/multierr
  version: v1.11.0
  license: MIT
- name: golang.org/x/crypto
  version: v0.18.0
  license: BSD-3-Clause
- name: golang.org/x/exp
  version: v0.0.0-20231206192017-f3f8817b8deb
  license: BSD-3-Clause
- name: golang.org/x/mod
  version: v0.14.0
  license: BSD-3-Clause
- name: golang.org/x/net
  version: v0.20.0
  license: BSD-3-Clause
- name: golang.org/x/oauth2
  version: v0.15.0
  license: BSD-3-Clause
- name: golang.org/x/sync
  version: v0.5.0
  license: BSD-3-Clause
- name: golang.org/x/sys
  version: v0.16.0
  license: BSD-3-Clause
- name: golang.org/x/term
  version: v0.16.0
  license: BSD-3-Clause
- name: golang.org/x/text
  version: v0.14.0
  license: BSD-3-Clause
- name: golang.org/x/time
  version: v0.5.0
  license: BSD-3-Clause
- name: golang.org/x/tools
  version: v0.16.1
  license: BSD-3-Clause
- name: gonum.org/v1/gonum
  version: v0.14.0
  license: BSD-3-Clause
- name: google.golang.org/api
  version: v0.152.0
  license: BSD-3-Clause
- name: google.golang.org/appengine
  version: v1.6.8
  license: Apache-2.0
- name: google.golang.org/genproto/googleapis/rpc
  version: v0.0.0-20231127180814-3a041ad873d4
  license: Apache-2.0
- name: google.golang.org/grpc
  version: v1.60.1
  license: Apache-2.0
- name: google.golang.org/protobuf
  version: v1.31.0
  license: BSD-3-Clause
- name: gopkg.in/inf.v0
  version: v0.9.1
  license: BSD-3-Clause
- name: gopkg.in/ini.v1
  version: v1.67.0
  license: Apache-2.0
- name: gopkg.in/warnings.v0
  version: v0.1.2
  license: BSD-2-Clause
- name: gopkg.in/yaml.v2
  version: v2.4.0
  license: Apache-2.0
- name: gopkg.in/yaml.v3
  version: v3.0.1
  license: Apache-2.0
- name: gorm.io/driver/postgres
  version: v1.5.3
  license: MIT
- name: gorm.io/driver/sqlite
  version: v1.5.4
  license: MIT
- name: gorm.io/gorm
  version: v1.25.5
  license: MIT
- name: helm.sh/helm/v3
  version: v3.13.2
  license: Apache-2.0
- name: k8s.io/api
  version: v0.28.4
  license: Apache-2.0
- name: k8s.io/apiextensions-apiserver
  version: v0.28.4
  license: Apache-2.0
- name: k8s.io/apimachinery
  version: v0.28.4
  license: Apache-2.0
- name: k8s.io/apiserver
  version: v0.28.4
  license: Apache-2.0
- name: k8s.io/cli-runtime
  version: v0.28.4
  license: Apache-2.0
- name: k8s.io/client-go
  version: v0.28.3
  license: Apache-2.0
- name: k8s.io/component-base
  version: v0.28.4
  license: Apache-2.0
- name: k8s.io/klog/v2
  version: v2.110.1
  license: Apache-2.0
- name: k8s.io/kube-openapi
  version: v0.0.0-20231113174909-778a5567bc1e
  license: Apache-2.0
- name: k8s.io/kubectl
  version: v0.28.3
  license: Apache-2.0
- name: k8s.io/utils
  version: v0.0.0-20230726121419-3b25d923346b
  license: Apache-2.0
- name: oras.land/oras-go
  version: v1.2.3
  license: Apache-2.0
- name: sigs.k8s.io/controller-runtime
  version: v0.16.3
  license: Apache-2.0
- name: sigs.k8s.io/json
  version: v0.0.0-20221116044647-bc3834ca7abd
  license: Apache-2.0
- name: sigs.k8s.io/kustomize/api
  version: v0.13.5-0.20230601165947-6ce0bf390ce3
  license: Apache-2.0
- name: sigs.k8s.io/kustomize/kyaml
  version: v0.14.3-0.20230601165947-6ce0bf390ce3
  license: Apache-2.0
- name: sigs.k8s.io/structured-merge-diff/v4
  version: v4.4.1
  license: Apache-2.0
- name: sigs.k8s.io/yaml
  version: v1.4.0
  license: MIT
