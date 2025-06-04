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
GOVERSION = 1.23
GOPATH = $(shell go env GOPATH)
GOBIN  = $(GOPATH)/bin
KEYS_PATH="../../server/permissions/keys.csv"

SHELL := /usr/bin/env bash -o pipefail

#-----------------------------------------------------------------------------
# Components
#-----------------------------------------------------------------------------
# All Adapters
# ADAPTER_URLS := "localhost:10000 localhost:10001 localhost:10002 localhost:10004 localhost:10005 localhost:10006 localhost:10007 localhost:10009 localhost:10010 localhost:10012"
# No Adapters
ADAPTER_URLS := "localhost:10000 localhost:10001 localhost:10012 localhost:10013"

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

#-----------------------------------------------------------------------------
# Server
#-----------------------------------------------------------------------------
MESHERY_K8S_SKIP_COMP_GEN ?= TRUE
APPLICATIONCONFIGPATH="./apps.json"
PORT:=9081

#-----------------------------------------------------------------------------
# Build
#-----------------------------------------------------------------------------
RELEASE_CHANNEL="edge"
