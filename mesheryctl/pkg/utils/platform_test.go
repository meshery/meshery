// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"
)

// The function are related to download should be test in meshkit package, please do not add test here.

func TestListManifests(t *testing.T) {
	t.Run("ListManifests with empty manifest", func(t *testing.T) {
		url := "https://api.github.com/repos/meshery/meshery/git/trees/47c634a49e6d143a54d734437a26ad233146ddf5"

		_, err := ListManifests(url)
		if err != nil {
			t.Errorf("ListManifests failed: %v", err)
		}
	})
}
