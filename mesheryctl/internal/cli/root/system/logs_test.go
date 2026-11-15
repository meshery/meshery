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
	"sync"
	"testing"
	"time"
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

// blockingContextReadCloser simulates an active log stream that signals when Read
// has begun, and then blocks until its stream context is cancelled.
type blockingContextReadCloser struct {
	ctx         context.Context
	readStarted chan struct{}
	closed      bool
	closeCalls  int
	mu          sync.Mutex
}

func (b *blockingContextReadCloser) Read(p []byte) (int, error) {
	b.mu.Lock()
	if b.readStarted != nil {
		select {
		case <-b.readStarted:
		default:
			close(b.readStarted)
		}
	}
	b.mu.Unlock()

	<-b.ctx.Done()
	return 0, b.ctx.Err()
}

func (b *blockingContextReadCloser) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closeCalls++
	b.closed = true
	return nil
}

// TestFollowLogStreams_LaterStreamCreationError verifies that when a later stream
// fails to open, followLogStreams cancels earlier active follow streams and waits
// for their reader goroutines to complete and release stream resources before returning.
func TestFollowLogStreams_LaterStreamCreationError(t *testing.T) {
	laterStreamErr := errors.New("simulated stream creation failure for pod-2")

	readStarted := make(chan struct{})
	stream1 := &blockingContextReadCloser{
		readStarted: readStarted,
	}

	sources := []streamLogSource{
		{
			podName: "pod-1",
			open: func(ctx context.Context) (io.ReadCloser, error) {
				stream1.ctx = ctx
				return stream1, nil
			},
		},
		{
			podName: "pod-2",
			open: func(ctx context.Context) (io.ReadCloser, error) {
				// Ensure stream 1's reader goroutine is actively blocked on Read before failing
				select {
				case <-readStarted:
				case <-time.After(2 * time.Second):
					return nil, errors.New("timed out waiting for stream 1 to start reading")
				}
				return nil, laterStreamErr
			},
		},
	}

	// Use context.Background() without a deadline so that stream cancellation
	// is solely driven by followLogStreams internal cancellation logic.
	parentCtx := context.Background()

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- followLogStreams(parentCtx, sources, func(string, string) {})
	}()

	// Watchdog timeout to prevent the test from hanging indefinitely if cancellation is broken
	var err error
	select {
	case err = <-resultCh:
		// Succeeded in completing within the watchdog window
	case <-time.After(3 * time.Second):
		t.Fatal("followLogStreams hung: failed to cancel earlier streams or wait for goroutines")
	}

	// Assert the original second-stream error is returned
	if !errors.Is(err, laterStreamErr) {
		t.Fatalf("expected error %v, got %v", laterStreamErr, err)
	}

	// Assert that stream 1 was cancelled and closed exactly once
	stream1.mu.Lock()
	closeCalls := stream1.closeCalls
	closed := stream1.closed
	stream1.mu.Unlock()

	if closeCalls != 1 {
		t.Errorf("stream1.closeCalls = %d, want 1", closeCalls)
	}
	if !closed {
		t.Errorf("stream1.closed = false, want true")
	}
}

// syncLogStream is a mock ReadCloser that provides channel-based barrier synchronization.
// It signals when Read begins, blocks until released by the test coordinator,
// returns its payload, and then on the subsequent read returns EOF.
type syncLogStream struct {
	data        []byte
	readEntered chan struct{}
	proceedRead chan struct{}
	mu          sync.Mutex
	readCalls   int
	closeCalls  int
	closed      bool
}

func (s *syncLogStream) Read(p []byte) (int, error) {
	s.mu.Lock()
	call := s.readCalls
	s.readCalls++
	s.mu.Unlock()

	switch call {
	case 0:
		if s.readEntered != nil {
			close(s.readEntered)
		}
		if s.proceedRead != nil {
			<-s.proceedRead
		}
		n := copy(p, s.data)
		return n, nil
	case 1:
		return 0, io.EOF
	default:
		return 0, errUnexpectedRead
	}
}

func (s *syncLogStream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closeCalls++
	s.closed = true
	return nil
}

// TestFollowLogStreams_NormalCompletion verifies that followLogStreams runs multiple
// streams concurrently, prints their log messages with the correct pod association,
// and terminates cleanly when all streams reach EOF, closing each stream exactly once.
func TestFollowLogStreams_NormalCompletion(t *testing.T) {
	stream1Started := make(chan struct{})
	stream2Started := make(chan struct{})
	proceedStreams := make(chan struct{})

	stream1 := &syncLogStream{
		data:        []byte("log message from pod-1\n"),
		readEntered: stream1Started,
		proceedRead: proceedStreams,
	}
	stream2 := &syncLogStream{
		data:        []byte("log message from pod-2\n"),
		readEntered: stream2Started,
		proceedRead: proceedStreams,
	}

	sources := []streamLogSource{
		{
			podName: "pod-1",
			open: func(ctx context.Context) (io.ReadCloser, error) {
				return stream1, nil
			},
		},
		{
			podName: "pod-2",
			open: func(ctx context.Context) (io.ReadCloser, error) {
				return stream2, nil
			},
		},
	}

	type recordedEntry struct {
		logMsg  string
		podName string
	}
	var mu sync.Mutex
	var recordedLogs []recordedEntry
	captureLog := func(logMsg, podName string) {
		mu.Lock()
		defer mu.Unlock()
		recordedLogs = append(recordedLogs, recordedEntry{logMsg: logMsg, podName: podName})
	}

	// Coordinator goroutine waits until BOTH streams have entered their Read loop,
	// proving that both streams are being consumed concurrently.
	coordinatorDone := make(chan struct{})
	go func() {
		defer close(coordinatorDone)
		select {
		case <-stream1Started:
		case <-time.After(3 * time.Second):
			t.Errorf("timed out waiting for stream1 to start reading")
			return
		}
		select {
		case <-stream2Started:
		case <-time.After(3 * time.Second):
			t.Errorf("timed out waiting for stream2 to start reading")
			return
		}
		// Both streams are actively in Read() concurrently. Release both.
		close(proceedStreams)
	}()

	resultCh := make(chan error, 1)
	go func() {
		resultCh <- followLogStreams(context.Background(), sources, captureLog)
	}()

	var err error
	select {
	case err = <-resultCh:
	case <-time.After(5 * time.Second):
		t.Fatal("followLogStreams timed out in normal completion")
	}

	<-coordinatorDone

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Verify both streams were closed exactly once upon EOF
	stream1.mu.Lock()
	s1Closed, s1CloseCalls, s1ReadCalls := stream1.closed, stream1.closeCalls, stream1.readCalls
	stream1.mu.Unlock()

	if !s1Closed || s1CloseCalls != 1 {
		t.Errorf("stream1 not closed properly: closed=%v, closeCalls=%d (want closed=true, closeCalls=1)", s1Closed, s1CloseCalls)
	}
	if s1ReadCalls != 2 {
		t.Errorf("stream1 readCalls = %d, want 2 (data + EOF)", s1ReadCalls)
	}

	stream2.mu.Lock()
	s2Closed, s2CloseCalls, s2ReadCalls := stream2.closed, stream2.closeCalls, stream2.readCalls
	stream2.mu.Unlock()

	if !s2Closed || s2CloseCalls != 1 {
		t.Errorf("stream2 not closed properly: closed=%v, closeCalls=%d (want closed=true, closeCalls=1)", s2Closed, s2CloseCalls)
	}
	if s2ReadCalls != 2 {
		t.Errorf("stream2 readCalls = %d, want 2 (data + EOF)", s2ReadCalls)
	}

	// Verify each stream's output is associated with the correct pod
	mu.Lock()
	entries := make([]recordedEntry, len(recordedLogs))
	copy(entries, recordedLogs)
	mu.Unlock()

	if len(entries) != 2 {
		t.Fatalf("expected 2 log entries, got %d", len(entries))
	}

	logsByPod := make(map[string]string)
	for _, entry := range entries {
		logsByPod[entry.podName] = entry.logMsg
	}

	if msg, ok := logsByPod["pod-1"]; !ok || msg != "log message from pod-1\n" {
		t.Errorf("pod-1 output mismatch: got %q, want %q", msg, "log message from pod-1\n")
	}
	if msg, ok := logsByPod["pod-2"]; !ok || msg != "log message from pod-2\n" {
		t.Errorf("pod-2 output mismatch: got %q, want %q", msg, "log message from pod-2\n")
	}
}
