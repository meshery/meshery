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
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAuthTokenFileCreatesOwnerOnly(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.json")

	if err := writeAuthTokenFile(path, []byte(`{"token":"x"}`)); err != nil {
		t.Fatalf("writeAuthTokenFile: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("new auth file mode = %04o, want 0600", perm)
	}
}

func TestWriteAuthTokenFileHardensExistingPermissiveFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.json")

	if err := os.WriteFile(path, []byte(`{"token":"old"}`), 0o600); err != nil {
		t.Fatalf("seed WriteFile: %v", err)
	}
	// Make the seed mode deterministic regardless of process umask.
	if err := os.Chmod(path, 0o644); err != nil {
		t.Fatalf("chmod seed: %v", err)
	}
	if info, err := os.Stat(path); err != nil {
		t.Fatalf("stat seed: %v", err)
	} else if perm := info.Mode().Perm(); perm != 0o644 {
		t.Fatalf("seed mode = %04o, want 0644", perm)
	}

	if err := writeAuthTokenFile(path, []byte(`{"token":"new"}`)); err != nil {
		t.Fatalf("writeAuthTokenFile: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("hardened auth file mode = %04o, want 0600", perm)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(data) != `{"token":"new"}` {
		t.Fatalf("content = %q, want updated token payload", data)
	}
}

func TestWriteAuthTokenFileClearsContents(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "auth.json")

	if err := writeAuthTokenFile(path, []byte(`{"token":"x"}`)); err != nil {
		t.Fatalf("seed writeAuthTokenFile: %v", err)
	}
	if err := writeAuthTokenFile(path, []byte{}); err != nil {
		t.Fatalf("clear writeAuthTokenFile: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(data) != 0 {
		t.Fatalf("cleared file len = %d, want 0", len(data))
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		t.Fatalf("cleared auth file mode = %04o, want 0600", perm)
	}
}
