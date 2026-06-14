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

package system

import "testing"

func TestIsPodRequired(t *testing.T) {
	type args struct {
		requiredPods []string
		pod          string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Pod '" + "test'" + " is required",
			args: args{
				requiredPods: []string{"test", "test1", "test2"},
				pod:          "test",
			},
			want: true,
		},
		{
			name: "Pod '" + "test'" + " is not required",
			args: args{
				requiredPods: []string{"test1", "test2"},
				pod:          "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPodRequired(tt.args.requiredPods, tt.args.pod); got != tt.want {
				t.Fatalf("IsPodRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
