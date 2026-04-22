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
REMOTE_PROVIDER="Layer5"

LOCAL_PROVIDER="None"
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
REMOTE_PROVIDER_LOCAL="http://localhost:9876"
EQUINIX_DEV="http://meshery.console.equinix.com"
EQUINIX_DEV2="http://meshery-2.console.equinix.com"
MESHERY_CLOUD_DEV="http://localhost:9876"
MESHERY_CLOUD_PROD="https://cloud.layer5.io"
MESHERY_CLOUD_STAGING="https://staging-cloud.layer5.io"
EXOSCALE_PROD="https://sks.exoscale.com"
EXOSCALE_STG="https://stg-sks.exoscale.com"
EXOSCALE_DEV="https://dev-sks.exoscale.com"
PROVIDER_CAPABILITIES_FILEPATH="" # Path to capabilities file for remote provider. If empty, capabilities will be fetched from remote provider.

#-----------------------------------------------------------------------------
# Server
#-----------------------------------------------------------------------------
MESHERY_K8S_SKIP_COMP_GEN ?= TRUE
APPLICATIONCONFIGPATH="./apps.json"
PORT:=9081
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

