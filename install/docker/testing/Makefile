.PHONY: download-go-modules build-meshery-server ui-deps provider-ui-deps ui-build provider-ui-build install-wrk2-build-deps configure-wrk2-git clone-and-build-wrk2 install-seed-content-deps seed-filters seed-patterns seed-applications install-runtime-deps create-runtime-user prepare-runtime-dirs

download-go-modules:
	go mod download

build-meshery-server:
	GOPROXY=https://proxy.golang.org go build -ldflags="-w -s -X main.version=$(GIT_VERSION)" -tags draft -a -o /meshery ./server/cmd

ui-deps:
	npm ci --omit=dev --legacy-peer-deps --ignore-scripts

provider-ui-deps:
	npm ci --omit=dev --legacy-peer-deps

ui-build:
	npm run build

provider-ui-build:
	npm run build

install-wrk2-build-deps:
	apt-get update
	apt-get install -y --no-install-recommends build-essential ca-certificates curl git libssl-dev zlib1g-dev
	rm -rf /var/lib/apt/lists/*

configure-wrk2-git:
	git config --global user.email "meshery@meshery.io"
	git config --global user.name "meshery"

clone-and-build-wrk2:
	git clone --depth=1 https://github.com/layer5io/wrk2
	$(MAKE) -C wrk2

install-seed-content-deps:
	apt-get update
	apt-get install -y --no-install-recommends ca-certificates curl
	rm -rf /var/lib/apt/lists/*

seed-filters:
	lines=$$(curl -fsSL --retry 5 --retry-all-errors --retry-delay 2 https://api.github.com/repos/layer5io/wasm-filters/releases/latest | grep "browser_download_url.*wasm" | cut -d : -f 2,3 | sed 's/"//g'); \
	mkdir -p seed_content/filters/binaries; \
	cd seed_content/filters/binaries; \
	for line in $$lines; do curl -fsSLO --retry 5 --retry-all-errors --retry-delay 2 "$$line"; done

seed-patterns:
	curl -fsSL --retry 5 --retry-all-errors --retry-delay 2 https://github.com/service-mesh-patterns/service-mesh-patterns/tarball/master -o service-mesh-patterns.tgz
	mkdir service-mesh-patterns
	mkdir -p /seed_content/patterns
	tar xzf service-mesh-patterns.tgz --directory=service-mesh-patterns
	mv service-mesh-patterns/*/samples/* /seed_content/patterns/

seed-applications:
	mkdir -p /seed_content/applications
	cd /seed_content/applications && curl -fsSLO --retry 5 --retry-all-errors --retry-delay 2 https://raw.githubusercontent.com/istio/istio/master/samples/bookinfo/platform/kube/bookinfo.yaml
	cd /seed_content/applications && curl -fsSLO --retry 5 --retry-all-errors --retry-delay 2 https://raw.githubusercontent.com/istio/istio/master/samples/httpbin/httpbin.yaml
	curl -fsSL --retry 5 --retry-all-errors --retry-delay 2 https://raw.githubusercontent.com/layer5io/image-hub/master/deployment.yaml -o /seed_content/applications/imagehub.yaml
	mkdir /emojivoto
	cd /emojivoto && curl --remote-name-all -fsSL --retry 5 --retry-all-errors --retry-delay 2 https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/emoji.yml https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/vote-bot.yml https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/voting.yml https://raw.githubusercontent.com/BuoyantIO/emojivoto/main/kustomize/deployment/web.yml
	awk 'FNR==1 && NR>1 { printf("\n%s\n\n","---") } 1' /emojivoto/*.yml > /seed_content/applications/emojivoto.yml

install-runtime-deps:
	apt-get update
	apt-get install -y --no-install-recommends ca-certificates curl libstdc++6
	rm -rf /var/lib/apt/lists/*
	update-ca-certificates

create-runtime-user:
	getent group docker >/dev/null || groupadd --system docker
	id -u appuser >/dev/null 2>&1 || useradd --system --create-home --home-dir /home/appuser --shell /usr/sbin/nologin appuser

prepare-runtime-dirs:
	mkdir -p /var/run /home/appuser/.meshery/config
	chown -R appuser:docker /var/run
	chown -R appuser:appuser /home/appuser /app/ui/public
