# Copyright Meshery Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#-----------------------------------------------------------------------------
# Global Variables
#-----------------------------------------------------------------------------
GIT_VERSION	= $(shell git describe --tags `git rev-list --tags --max-count=1`)
GIT_COMMITSHA = $(shell git rev-list -1 HEAD)
GIT_STRIPPED_VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1` | cut -c 2-)

# Extension Point for remote provider . Add your provider here.
# Empty by default so installs do not enforce a specific provider; users see
# the provider-selection UI on first launch. Set to e.g. "Meshery" or "Layer5"
# to enforce a single provider via the PROVIDER env var.
REMOTE_PROVIDER=""

LOCAL_PROVIDER="Local"
GOVERSION = 1.25
GOPATH = $(shell go env GOPATH)
GOBIN  = $(GOPATH)/bin
KEYS_PATH="../../server/permissions/keys.csv"

SHELL := /usr/bin/env bash -o pipefail

#-----------------------------------------------------------------------------
# Components
#-----------------------------------------------------------------------------
# Adapters to dial on start. Default is empty so developers who do not run
# adapters locally don't see meshery-server-1047 "connection refused" errors
# every sync tick. To use adapters, export ADAPTER_URLS in your shell, e.g.:
#   export ADAPTER_URLS='localhost:10000 localhost:10001 localhost:10012 localhost:10013'
# or override on the make command line:
#   make server-local ADAPTER_URLS='localhost:10000 localhost:10001'
ADAPTER_URLS ?= ""

#-----------------------------------------------------------------------------
# Providers (Add your provider here. See https://docs.meshery.io/extensibility/providers)
#-----------------------------------------------------------------------------
# Dev / staging / local provider hosts. These are intentional overrides used
# by specific make targets (server-stg, server-local-provider, etc.) and are
# NOT part of the canonical default provider list.
REMOTE_PROVIDER_LOCAL="http://localhost:9876"
EQUINIX_DEV="http://meshery.console.equinix.com"
EQUINIX_DEV2="http://meshery-2.console.equinix.com"
MESHERY_CLOUD_DEV="http://localhost:9876"
MESHERY_CLOUD_STAGING="https://staging-cloud.meshery.io"
EXOSCALE_STG="https://stg-sks.exoscale.com"
EXOSCALE_DEV="https://dev-sks.exoscale.com"
PROVIDER_CAPABILITIES_FILEPATH="" # Path to capabilities file for remote provider. If empty, capabilities will be fetched from remote provider.

# --- AUTO-SYNC SOURCE: edit here, then run `make sync-provider-defaults` ---
# Canonical list of default remote providers. This is the single source of
# truth - every consumer (helm values.yaml, docker-compose, k8s manifests,
# mesheryctl Services map, UI constants, docker-extension chooser, the
# server's viper SetDefault, etc.) is regenerated from these variables by
# scripts/sync-provider-defaults.sh. To add or change a default provider,
# add it to REMOTE_PROVIDER_URLS (and to the PAIRS array in
# scripts/sync-provider-defaults.sh) and run `make sync-provider-defaults`.
#
# Active default providers. These MUST resolve in DNS - the server retries
# each unreachable URL 10x with a 3s sleep at startup (see
# server/models/remote_provider.go loadCapabilities), so adding an
# unreachable host adds ~30s to startup before the HTTP listener accepts
# traffic, which breaks the e2e and meshsync integration tests.
MESHERY_CLOUD_PROD="https://cloud.meshery.io"
LAYER5_CLOUD_PROD="https://cloud.layer5.io"
# Declared-but-not-yet-active providers. Their hostnames do not currently
# resolve in DNS, so registering them at startup would block the server
# behind retry timeouts. Keep the declarations here so re-enabling is a
# one-line change once DNS is set up; move into REMOTE_PROVIDER_URLS
# below (and add to PAIRS in scripts/sync-provider-defaults.sh) when ready.
MESHERY_DIGITALOCEAN_PROD="https://meshery.digitalocean.com"
CLEVERLUCK_PROD="https://idp.cleverluck.com"
EXOSCALE_PROD="https://designer.exoscale.com"
INTEL_PROD="https://perf.platform.intel.com"
UTAUSTIN_PROD="https://ppf.research.utexas.edu"
TCSLABS_PROD="https://tcs-labs.in"
REMOTE_PROVIDER_URLS=$(MESHERY_CLOUD_PROD),$(LAYER5_CLOUD_PROD)
PRIMARY_PROVIDER_URL=$(MESHERY_CLOUD_PROD)
# Display names paired with each URL (consumed by the docker-extension
# chooser generator to produce {name, url} pairs). Keep parallel with
# REMOTE_PROVIDER_URLS - only the active providers are rendered.
MESHERY_NAME="Meshery"
LAYER5_NAME="Layer5"
MESHERY_DIGITALOCEAN_NAME="DigitalOcean"
CLEVERLUCK_NAME="CleverLuck"
EXOSCALE_NAME="Exoscale"
INTEL_NAME="Intel"
UTAUSTIN_NAME="UT Austin"
TCSLABS_NAME="TCS Labs"
# --- END AUTO-SYNC SOURCE ---

# Helper target for the sync script: `make -s print-VARNAME` echoes the
# resolved value of any Make variable. Used by scripts/sync-provider-defaults.sh
# to read the canonical block above.
print-%:
	@echo $($*)

#-----------------------------------------------------------------------------
# Server
#-----------------------------------------------------------------------------
MESHERY_K8S_SKIP_COMP_GEN ?= TRUE
APPLICATIONCONFIGPATH="./apps.json"
PORT ?= 9081
USE_GO_POLICY_ENGINE ?= true
# OpenTelemetry Config (Ansi-C string format). Defaults to empty so tracing
# is disabled unless a developer explicitly points it at a live collector.
# Previously this defaulted to localhost:4317 which floods logs with
# "traces export: ... connection refused" every ~10s when no collector is
# running. To enable tracing from bash, override on the command line, e.g.:
#   make server-local OTEL_CONFIG=$'service_name: meshery-server\nservice_version: 1.0.0\nendpoint: localhost:4317\ninsecure: true'
OTEL_CONFIG ?=
#-----------------------------------------------------------------------------
# Build
#-----------------------------------------------------------------------------
RELEASE_CHANNEL="edge"

