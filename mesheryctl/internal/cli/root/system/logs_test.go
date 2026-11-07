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

import (
	"context"
	"errors"
	"io"
	"testing"
)

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

var errUnexpectedRead = errors.New("unexpected additional read")

type mockReadResult struct {
	data []byte
	err  error
}

type mockReadCloser struct {
	results    []mockReadResult
	closed     bool
	closeCalls int
	readCalls  int
}

func (m *mockReadCloser) Read(p []byte) (int, error) {
	if m.readCalls >= len(m.results) {
		m.readCalls++
		return 0, errUnexpectedRead
	}

	res := m.results[m.readCalls]
	m.readCalls++

	n := 0
	if len(res.data) > 0 {
		n = copy(p, res.data)
	}

	return n, res.err
}

func (m *mockReadCloser) Close() error {
	m.closeCalls++
	m.closed = true
	return nil
}

func TestProcessLogStream(t *testing.T) {
	genericErr := errors.New("read error: connection reset")

	tests := []struct {
		name               string
		podName            string
		results            []mockReadResult
		expectedLogs       []string
		expectedPods       []string
		expectedErr        error
		expectedCloseCalls int
		expectedReadCalls  int
	}{
		{
			name:    "Normal stream completion with multiple chunks",
			podName: "meshery-server-1",
			results: []mockReadResult{
				{data: []byte("chunk 1\n"), err: nil},
				{data: []byte("chunk 2\n"), err: nil},
				{data: nil, err: io.EOF},
			},
			expectedLogs:       []string{"chunk 1\n", "chunk 2\n"},
			expectedPods:       []string{"meshery-server-1", "meshery-server-1"},
			expectedErr:        nil,
			expectedCloseCalls: 1,
			expectedReadCalls:  3,
		},
		{
			name:    "EOF with final bytes in single read",
			podName: "meshery-server-2",
			results: []mockReadResult{
				{data: []byte("final chunk\n"), err: io.EOF},
			},
			expectedLogs:       []string{"final chunk\n"},
			expectedPods:       []string{"meshery-server-2"},
			expectedErr:        nil,
			expectedCloseCalls: 1,
			expectedReadCalls:  1,
		},
		{
			name:    "Empty stream immediate EOF",
			podName: "meshery-server-3",
			results: []mockReadResult{
				{data: nil, err: io.EOF},
			},
			expectedLogs:       nil,
			expectedPods:       nil,
			expectedErr:        nil,
			expectedCloseCalls: 1,
			expectedReadCalls:  1,
		},
		{
			name:    "Unexpected read error accompanied by final bytes",
			podName: "meshery-server-4",
			results: []mockReadResult{
				{data: []byte("partial output before failure\n"), err: genericErr},
			},
			expectedLogs:       []string{"partial output before failure\n"},
			expectedPods:       []string{"meshery-server-4"},
			expectedErr:        genericErr,
			expectedCloseCalls: 1,
			expectedReadCalls:  1,
		},
		{
			name:    "Context cancellation terminates quietly",
			podName: "meshery-server-5",
			results: []mockReadResult{
				{data: []byte("before cancel\n"), err: nil},
				{data: nil, err: context.Canceled},
			},
			expectedLogs:       []string{"before cancel\n"},
			expectedPods:       []string{"meshery-server-5"},
			expectedErr:        nil,
			expectedCloseCalls: 1,
			expectedReadCalls:  2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReader := &mockReadCloser{
				results: tt.results,
			}

			var recordedLogs []string
			var recordedPods []string
			captureLog := func(logMsg string, pod string) {
				recordedLogs = append(recordedLogs, logMsg)
				recordedPods = append(recordedPods, pod)
			}

			err := processLogStream(mockReader, tt.podName, captureLog)

			if tt.expectedErr != nil {
				if err == nil || !errors.Is(err, tt.expectedErr) {
					t.Fatalf("processLogStream() error = %v, expectedErr = %v", err, tt.expectedErr)
				}
			} else {
				if err != nil {
					t.Fatalf("processLogStream() unexpected error = %v", err)
				}
			}

			if mockReader.closeCalls != tt.expectedCloseCalls {
				t.Errorf("mockReader.closeCalls = %d, want %d", mockReader.closeCalls, tt.expectedCloseCalls)
			}
			if !mockReader.closed {
				t.Errorf("mockReader.closed = false, want true")
			}
			if mockReader.readCalls != tt.expectedReadCalls {
				t.Errorf("mockReader.readCalls = %d, want %d", mockReader.readCalls, tt.expectedReadCalls)
			}

			if len(recordedLogs) != len(tt.expectedLogs) {
				t.Fatalf("recordedLogs length = %d (%v), want %d (%v)",
					len(recordedLogs), recordedLogs, len(tt.expectedLogs), tt.expectedLogs)
			}
			for i := range recordedLogs {
				if recordedLogs[i] != tt.expectedLogs[i] {
					t.Errorf("recordedLogs[%d] = %q, want %q", i, recordedLogs[i], tt.expectedLogs[i])
				}
				if recordedPods[i] != tt.expectedPods[i] {
					t.Errorf("recordedPods[%d] = %q, want %q", i, recordedPods[i], tt.expectedPods[i])
				}
			}
		})
	}
}
