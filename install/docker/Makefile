.PHONY: build-meshery-server ui-deps provider-ui-deps ui-build provider-ui-build install-seed-content-deps seed-filters prepare-build-time-extensions install-runtime-deps create-runtime-user prepare-runtime-dirs

build-meshery-server:
	go clean -modcache
	cd server/cmd && GOPROXY=https://proxy.golang.org GOSUMDB=off go build -ldflags="-w -s -X main.globalTokenForAnonymousResults=$(TOKEN) -X main.version=$(GIT_VERSION) -X main.commitsha=$(GIT_COMMITSHA) -X main.releasechannel=$(RELEASE_CHANNEL)" -tags draft -a -o /meshery .

ui-deps:
	npm ci --omit=dev --legacy-peer-deps --ignore-scripts

provider-ui-deps:
	npm ci --omit=dev --legacy-peer-deps

ui-build:
	npm run build

provider-ui-build:
	npm run build

install-seed-content-deps:
	apt-get update
	apt-get install -y --no-install-recommends ca-certificates curl
	rm -rf /var/lib/apt/lists/*

seed-filters:
	lines=$$(curl -fsSL --retry 5 --retry-all-errors --retry-delay 2 https://api.github.com/repos/layer5io/wasm-filters/releases/latest | grep "browser_download_url.*wasm" | cut -d : -f 2,3 | sed 's/"//g'); \
	mkdir -p seed_content/filters/binaries; \
	cd seed_content/filters/binaries; \
	for line in $$lines; do curl -fsSLO --retry 5 --retry-all-errors --retry-delay 2 "$$line"; done

prepare-build-time-extensions:
	if [ -d /meshery-context/.meshery ]; then \
		cp -r /meshery-context/.meshery /build-ext/.meshery; \
	else \
		mkdir -p /build-ext/.meshery; \
	fi

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
