// Copyright 2019 The Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root"
)

<<<<<<< HEAD
var (
	version   = "Unavailable"
	commitsha = "Unavailable"
)

func main() {
	root.Build = version
	root.CommitSHA = commitsha
=======
var version = "dev"
var clientcommitsha = "SHA-Client"
var servercommitsha = "SHA-Server"

func main() {
	root.BuildClient = version
	root.BuildServer = version
	root.ClientCommitSHA = clientcommitsha
	root.ServerCommitSHA = servercommitsha
>>>>>>> This PR fixes #603
	root.Execute()
}
