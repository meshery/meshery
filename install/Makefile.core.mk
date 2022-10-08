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

GOVERSION = 1.19.1
GOPATH = $(shell go env GOPATH)
GOBIN  = $(GOPATH)/bin

SHELL :=/bin/bash -o pipefail

#-----------------------------------------------------------------------------
# Components
#-----------------------------------------------------------------------------
ADAPTER_URLS := "localhost:10000 localhost:10001 localhost:10002 localhost:10004 localhost:10005 localhost:10006 localhost:10007 localhost:10009 localhost:10010 localhost:10012"

#-----------------------------------------------------------------------------
# Providers
#-----------------------------------------------------------------------------
REMOTE_PROVIDER_LOCAL="http://localhost:9876"
MESHERY_CLOUD_DEV="http://localhost:9876"
MESHERY_CLOUD_PROD="https://meshery.layer5.io"
MESHERY_CLOUD_STAGING="https://staging-meshery.layer5.io"

#-----------------------------------------------------------------------------
# Server
#-----------------------------------------------------------------------------
MESHERY_K8S_SKIP_COMP_GEN ?= TRUE
APPLICATIONCONFIGPATH="./apps.json"

#-----------------------------------------------------------------------------
# Build
#-----------------------------------------------------------------------------
RELEASE_CHANNEL="edge"
