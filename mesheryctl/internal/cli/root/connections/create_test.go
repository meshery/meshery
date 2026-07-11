// Copyright Meshery Authors
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

package connections

import (
	"fmt"
	"net"
	"net/url"
	"syscall"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
)

func TestIsMesheryServerUnreachable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "application error",
			err:  fmt.Errorf("no contexts found"),
			want: false,
		},
		{
			name: "windows connectex message",
			err:  fmt.Errorf(`Post "http://localhost:9081/api/system/kubernetes/contexts": dial tcp [::1]:9081: connectex: No connection could be made because the target machine actively refused it.`),
			want: true,
		},
		{
			name: "unix connection refused message",
			err:  fmt.Errorf(`Post "http://localhost:9081/api/system/kubernetes/contexts": dial tcp 127.0.0.1:9081: connect: connection refused`),
			want: true,
		},
		{
			name: "meshkit-style wrapped request error",
			err:  utils.ErrRequestResponse(fmt.Errorf(`Post "http://localhost:9081/api/system/kubernetes/contexts": dial tcp [::1]:9081: connectex: No connection could be made because the target machine actively refused it.`)),
			want: true,
		},
		{
			name: "url.Error with net.OpError",
			err: &url.Error{
				Op:  "Post",
				URL: "http://localhost:9081/api/system/kubernetes/contexts",
				Err: &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED},
			},
			want: true,
		},
		{
			name: "pkg/errors wrapped dial message",
			err:  errors.Wrap(fmt.Errorf("dial tcp 127.0.0.1:9081: connection refused"), "request failed"),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMesheryServerUnreachable(tt.err); got != tt.want {
				t.Fatalf("isMesheryServerUnreachable() = %v, want %v (err=%v)", got, tt.want, tt.err)
			}
		})
	}
}
