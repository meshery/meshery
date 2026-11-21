//go:build windows

package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

func TestWriteTokenFileRestrictsWindowsACL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("old token")); err != nil {
		t.Fatal(err)
	}
	// Model an existing token file that grants access to everyone.
	sd, err := windows.SecurityDescriptorFromString("D:(A;;FA;;;WD)")
	if err != nil {
		t.Fatal(err)
	}
	acl, _, err := sd.DACL()
	if err != nil {
		t.Fatal(err)
	}
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION, nil, nil, acl, nil); err != nil {
		t.Fatal(err)
	}
	if err := WriteTokenFile(path, []byte("new token")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "new token" {
		t.Fatalf("read token: %q, %v", data, err)
	}
	descriptor, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION|windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatal(err)
	}
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		t.Fatal(err)
	}
	owner, _, err := descriptor.Owner()
	if err != nil || owner == nil || !windows.EqualSid(owner, user.User.Sid) {
		t.Fatalf("unexpected token owner: %v", err)
	}
	control, _, err := descriptor.Control()
	if err != nil || control&windows.SE_DACL_PROTECTED == 0 {
		t.Fatalf("token DACL is not protected: %v", err)
	}
	dacl, _, err := descriptor.DACL()
	if err != nil || dacl == nil || dacl.AceCount != 1 {
		t.Fatalf("unexpected token DACL: %v", err)
	}
	var ace *windows.ACCESS_ALLOWED_ACE
	if err := windows.GetAce(dacl, 0, &ace); err != nil {
		t.Fatal(err)
	}
	if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE || !windows.EqualSid((*windows.SID)(unsafe.Pointer(&ace.SidStart)), user.User.Sid) {
		t.Fatal("token DACL permits an unexpected account")
	}
	const fileAllAccess = 0x001f01ff
	if ace.Mask != windows.GENERIC_ALL && ace.Mask != fileAllAccess {
		t.Fatalf("unexpected token ACE rights: %#x", ace.Mask)
	}
}

func TestWriteTokenFileRejectsOpenReader(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("old token")); err != nil {
		t.Fatal(err)
	}
	reader, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteTokenFile(path, []byte("new token")); err == nil {
		reader.Close()
		t.Fatal("expected an open reader to block the token update")
	}
	reader.Close()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "old token" {
		t.Fatalf("old token was changed: %q, %v", data, err)
	}
}

func TestWriteTokenFileReplacesFileWithRetainedWriteDACHandle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("old token")); err != nil {
		t.Fatal(err)
	}
	name, err := windows.UTF16PtrFromString(path)
	if err != nil {
		t.Fatal(err)
	}
	old, err := windows.CreateFile(name, windows.WRITE_DAC|windows.READ_CONTROL,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil, windows.OPEN_EXISTING, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		t.Fatal(err)
	}
	defer windows.CloseHandle(old)
	if err := WriteTokenFile(path, []byte("new token")); err != nil {
		t.Fatalf("replace token with old metadata handle open: %v", err)
	}
	// An attacker can still change the old file's DACL with the retained handle.
	permissive, err := windows.SecurityDescriptorFromString("D:(A;;FA;;;WD)")
	if err != nil {
		t.Fatal(err)
	}
	oldACL, _, err := permissive.DACL()
	if err != nil {
		t.Fatal(err)
	}
	if err := windows.SetSecurityInfo(old, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION, nil, nil, oldACL, nil); err != nil {
		t.Fatalf("change old DACL: %v", err)
	}
	var read uint32
	if err := windows.ReadFile(old, make([]byte, 32), &read, nil); err == nil {
		t.Fatal("metadata-only handle unexpectedly read the replacement token")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "new token" {
		t.Fatalf("read replacement token: %q, %v", data, err)
	}
	current, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatal(err)
	}
	acl, _, err := current.DACL()
	if err != nil || acl == nil || acl.AceCount != 1 {
		t.Fatalf("unexpected replacement DACL: %v", err)
	}
	var ace *windows.ACCESS_ALLOWED_ACE
	if err := windows.GetAce(acl, 0, &ace); err != nil {
		t.Fatal(err)
	}
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		t.Fatal(err)
	}
	if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE || !windows.EqualSid((*windows.SID)(unsafe.Pointer(&ace.SidStart)), user.User.Sid) {
		t.Fatal("old handle changed the replacement file's DACL")
	}
}

func TestTokenFileOwnerRejectsDifferentSID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("old token")); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	other, err := windows.StringToSid("S-1-1-0") // Everyone is not a user owner.
	if err != nil {
		t.Fatal(err)
	}
	if err := verifyTokenFileOwner(file, other); err == nil {
		t.Fatal("expected an owner mismatch to be rejected")
	}
}

func TestTokenFileOwnerTrust(t *testing.T) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		t.Fatal(err)
	}
	admins, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		t.Fatal(err)
	}
	other, err := windows.StringToSid("S-1-1-0")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		name                string
		owner, defaultOwner *windows.SID
		want                bool
	}{
		{"current user", user.User.Sid, admins, true},
		{"default administrators", admins, admins, true},
		{"administrators not default", admins, user.User.Sid, false},
		{"other default owner", other, other, false},
		{"other owner", other, admins, false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTrustedTokenOwner(tt.owner, user.User.Sid, tt.defaultOwner, admins); got != tt.want {
				t.Fatalf("isTrustedTokenOwner = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenFileOwnerDefaultAdministrator(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	// Let Windows assign the token's default owner to an ordinary new file.
	if err := os.WriteFile(path, []byte("old token"), 0600); err != nil {
		t.Fatal(err)
	}
	descriptor, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION)
	if err != nil {
		t.Fatal(err)
	}
	owner, _, err := descriptor.Owner()
	if err != nil {
		t.Fatal(err)
	}
	admins, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		t.Fatal(err)
	}
	if owner == nil || !windows.EqualSid(owner, admins) {
		t.Skip("this Windows token does not default to the Administrators owner")
	}
	if err := WriteTokenFile(path, []byte("new token")); err != nil {
		t.Fatalf("write token with default Administrators owner: %v", err)
	}
}

func TestWriteTokenFileLongPath(t *testing.T) {
	dir := t.TempDir()
	for len(dir) < 270 {
		dir = filepath.Join(dir, strings.Repeat("a", 40))
	}
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "token")
	if err := WriteTokenFile(path, []byte("token")); err != nil {
		t.Fatalf("write long token path: %v", err)
	}
}

func TestTokenFileLongPathProducesNTNamespace(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "local drive path",
			input: `C:\Users\test\.meshery\token.json`,
			want:  `\??\C:\Users\test\.meshery\token.json`,
		},
		{
			name:  "Win32 extended already",
			input: `\\?\C:\Users\test\.meshery\token.json`,
			want:  `\??\C:\Users\test\.meshery\token.json`,
		},
		{
			name:  "NT namespace already",
			input: `\??\C:\Users\test\.meshery\token.json`,
			want:  `\??\C:\Users\test\.meshery\token.json`,
		},
		{
			name:  "Win32 extended UNC",
			input: `\\?\UNC\server\share\token.json`,
			want:  `\??\UNC\server\share\token.json`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tokenFileLongPath(tc.input)
			if err != nil {
				t.Fatalf("tokenFileLongPath(%q): %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("tokenFileLongPath(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestTokenFileWin32PathPreservesShortPaths(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "short local path",
			input: `C:\Users\test\token.json`,
			want:  `C:\Users\test\token.json`,
		},
		{
			name:  "Win32 extended unchanged",
			input: `\\?\C:\path\token.json`,
			want:  `\\?\C:\path\token.json`,
		},
		{
			name:  "NT namespace converted",
			input: `\??\C:\path\token.json`,
			want:  `\\?\C:\path\token.json`,
		},
		{
			name:  "NT UNC converted",
			input: `\??\UNC\server\share\token.json`,
			want:  `\\?\UNC\server\share\token.json`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tokenFileWin32Path(tc.input)
			if err != nil {
				t.Fatalf("tokenFileWin32Path(%q): %v", tc.input, err)
			}
			if got != tc.want {
				t.Fatalf("tokenFileWin32Path(%q) = %q; want %q", tc.input, got, tc.want)
			}
		})
	}
}

func TestTokenFileMutexNameCanonicalization(t *testing.T) {
	base := `C:\Users\test\.meshery\token.json`
	aliases := []string{
		base,
		`c:\users\test\.meshery\token.json`,
		`\\?\C:\Users\test\.meshery\token.json`,
		`\??\C:\Users\test\.meshery\token.json`,
	}

	name1, err := tokenFileMutexName(aliases[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, alias := range aliases[1:] {
		name, err := tokenFileMutexName(alias)
		if err != nil {
			t.Fatal(err)
		}
		if name != name1 {
			t.Fatalf("expected identical mutex names for aliases %q and %q: %q != %q", aliases[0], alias, name1, name)
		}
	}

	uncBase := `\\server\share\.meshery\token.json`
	uncAliases := []string{
		uncBase,
		`\\?\UNC\server\share\.meshery\token.json`,
		`\??\UNC\server\share\.meshery\token.json`,
	}
	uncName, err := tokenFileMutexName(uncAliases[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, alias := range uncAliases[1:] {
		name, err := tokenFileMutexName(alias)
		if err != nil {
			t.Fatal(err)
		}
		if name != uncName {
			t.Fatalf("expected identical mutex names for UNC aliases %q and %q: %q != %q", uncAliases[0], alias, uncName, name)
		}
	}

	relative := filepath.Join(".meshery", "token.json")
	absolute, err := filepath.Abs(relative)
	if err != nil {
		t.Fatal(err)
	}
	relativeName, err := tokenFileMutexName(relative)
	if err != nil {
		t.Fatal(err)
	}
	absoluteName, err := tokenFileMutexName(absolute)
	if err != nil {
		t.Fatal(err)
	}
	if relativeName != absoluteName {
		t.Fatalf("expected relative and absolute paths to share a mutex: %q != %q", relativeName, absoluteName)
	}

	diffName, err := tokenFileMutexName(`C:\Users\test\.meshery\other.json`)
	if err != nil {
		t.Fatal(err)
	}
	if name1 == diffName {
		t.Fatal("expected distinct paths to produce different mutex names")
	}

	if !strings.HasPrefix(name1, `Local\meshery-token-`) {
		t.Fatalf("expected Local namespace prefix, got %q", name1)
	}
}

func TestHandleBasedRenameFallbackToStandard(t *testing.T) {
	origEx := replaceTokenFileExFunc
	defer func() { replaceTokenFileExFunc = origEx }()

	// Simulate a filesystem where FileRenameInfoEx (class 22) returns ERROR_NOT_SUPPORTED.
	replaceTokenFileExFunc = func(file *os.File, name []uint16) error {
		return windows.ERROR_NOT_SUPPORTED
	}

	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("fallback-standard-token")); err != nil {
		t.Fatalf("expected fallback to standard handle-based rename to succeed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil || string(data) != "fallback-standard-token" {
		t.Fatalf("read fallback token: %q, %v", data, err)
	}
}

func TestSameWindowsFileIdentity(t *testing.T) {
	info := windows.ByHandleFileInformation{
		VolumeSerialNumber: 12345,
		FileIndexHigh:      1,
		FileIndexLow:       2,
	}

	if !sameWindowsFileIdentity(info, info) {
		t.Fatal("expected identical file info to match")
	}

	diffVol := info
	diffVol.VolumeSerialNumber = 99999
	if sameWindowsFileIdentity(info, diffVol) {
		t.Fatal("expected different volume serial number to mismatch")
	}

	diffIdxHi := info
	diffIdxHi.FileIndexHigh = 9
	if sameWindowsFileIdentity(info, diffIdxHi) {
		t.Fatal("expected different file index high to mismatch")
	}

	diffIdxLo := info
	diffIdxLo.FileIndexLow = 9
	if sameWindowsFileIdentity(info, diffIdxLo) {
		t.Fatal("expected different file index low to mismatch")
	}

	zeroInfo := windows.ByHandleFileInformation{}
	if sameWindowsFileIdentity(zeroInfo, zeroInfo) {
		t.Fatal("expected all-zero file info to be rejected (cannot guarantee identity)")
	}
}

func TestWriteTokenFileFallbackIdentityPreserved(t *testing.T) {
	orig := replaceTokenFileFunc
	replaceTokenFileFunc = func(file *os.File, destination string) error {
		return windows.ERROR_NOT_SUPPORTED
	}
	defer func() { replaceTokenFileFunc = orig }()

	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("fallback token")); err != nil {
		t.Fatalf("WriteTokenFile fallback failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil || string(data) != "fallback token" {
		t.Fatalf("read fallback token: %q, %v", data, err)
	}
}

func TestHandleBasedRenameFailsClosedOnAccessDenied(t *testing.T) {
	origEx := replaceTokenFileExFunc
	origStd := replaceTokenFileStdFunc
	defer func() {
		replaceTokenFileExFunc = origEx
		replaceTokenFileStdFunc = origStd
	}()

	exCalled := false
	replaceTokenFileExFunc = func(file *os.File, name []uint16) error {
		exCalled = true
		return windows.ERROR_ACCESS_DENIED
	}

	stdCalled := false
	replaceTokenFileStdFunc = func(file *os.File, name []uint16) error {
		stdCalled = true
		return nil
	}

	path := filepath.Join(t.TempDir(), "token")
	err := WriteTokenFile(path, []byte("token"))
	if err == nil {
		t.Fatal("expected access denied error to fail closed")
	}
	if !exCalled {
		t.Fatal("expected extended rename stub to be called, but it was not")
	}
	if stdCalled {
		t.Fatal("expected ERROR_ACCESS_DENIED NOT to fall back to standard rename")
	}
	if !errors.Is(err, windows.ERROR_ACCESS_DENIED) {
		t.Fatalf("expected errors.Is(err, ERROR_ACCESS_DENIED) to be true; got: %v", err)
	}
}

func TestHandleBasedStandardRenameAccessDeniedFailsClosed(t *testing.T) {
	origEx := replaceTokenFileExFunc
	origStd := replaceTokenFileStdFunc
	defer func() {
		replaceTokenFileExFunc = origEx
		replaceTokenFileStdFunc = origStd
	}()

	replaceTokenFileExFunc = func(file *os.File, name []uint16) error {
		return windows.ERROR_NOT_SUPPORTED
	}
	replaceTokenFileStdFunc = func(file *os.File, name []uint16) error {
		return windows.ERROR_ACCESS_DENIED
	}

	path := filepath.Join(t.TempDir(), "token")
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	err = replaceTokenFile(file, path)
	if !errors.Is(err, windows.ERROR_ACCESS_DENIED) {
		t.Fatalf("expected standard rename access denied; got: %v", err)
	}
	if errors.Is(err, windows.ERROR_NOT_SUPPORTED) {
		t.Fatalf("access-denied result must not retain the unsupported fallback condition: %v", err)
	}
}

func TestConcurrentWritersProduceConsistentToken(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")

	const numWriters = 10
	var wg sync.WaitGroup
	wg.Add(numWriters)

	validTokens := make(map[string]bool)
	var mu sync.Mutex

	for i := 0; i < numWriters; i++ {
		tok := fmt.Sprintf("token-payload-%d", i)
		mu.Lock()
		validTokens[tok] = true
		mu.Unlock()

		go func(content string) {
			defer wg.Done()
			if err := WriteTokenFile(path, []byte(content)); err != nil {
				t.Errorf("concurrent WriteTokenFile failed: %v", err)
			}
		}(tok)
	}

	wg.Wait()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read final token file: %v", err)
	}

	if !validTokens[string(data)] {
		t.Fatalf("destination file contains corrupted or partial token: %q", data)
	}
}

func TestConcurrentWritersDifferentPaths(t *testing.T) {
	dir := t.TempDir()
	path1 := filepath.Join(dir, "token1")
	path2 := filepath.Join(dir, "token2")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := WriteTokenFile(path1, []byte("token-1")); err != nil {
			t.Errorf("write path1: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := WriteTokenFile(path2, []byte("token-2")); err != nil {
			t.Errorf("write path2: %v", err)
		}
	}()

	wg.Wait()

	data1, err1 := os.ReadFile(path1)
	data2, err2 := os.ReadFile(path2)
	if err1 != nil || string(data1) != "token-1" {
		t.Fatalf("path1 mismatch: %q, %v", data1, err1)
	}
	if err2 != nil || string(data2) != "token-2" {
		t.Fatalf("path2 mismatch: %q, %v", data2, err2)
	}
}

func TestFailedWriterPreservesExistingDestination(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")

	if err := WriteTokenFile(path, []byte("original valid token")); err != nil {
		t.Fatalf("initial write failed: %v", err)
	}

	origEx := replaceTokenFileExFunc
	defer func() { replaceTokenFileExFunc = origEx }()

	// Simulate an error during rename
	replaceTokenFileExFunc = func(file *os.File, name []uint16) error {
		return windows.ERROR_ACCESS_DENIED
	}

	err := WriteTokenFile(path, []byte("new token that should fail"))
	if err == nil {
		t.Fatal("expected write to fail")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read destination after failed write: %v", err)
	}
	if string(data) != "original valid token" {
		t.Fatalf("destination token was clobbered by failed writer: %q", data)
	}

	// Verify no orphaned temporary files remain in dir
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".meshery-token-") {
			t.Fatalf("found leftover temporary file %q after failed write", entry.Name())
		}
	}
}

func TestFailedWriterDoesNotDeleteDestinationWonByAnotherWriter(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")

	// Writer A starts and prepares temporary file.
	// Meanwhile Writer B successfully writes its token to destination.
	if err := WriteTokenFile(path, []byte("token from Writer B")); err != nil {
		t.Fatalf("Writer B failed: %v", err)
	}

	// Now simulate Writer A failing during its write attempt.
	origEx := replaceTokenFileExFunc
	defer func() { replaceTokenFileExFunc = origEx }()
	replaceTokenFileExFunc = func(file *os.File, name []uint16) error {
		return windows.ERROR_ACCESS_DENIED
	}

	err := WriteTokenFile(path, []byte("token from Writer A"))
	if err == nil {
		t.Fatal("expected Writer A to fail")
	}

	// Assert: Writer B's destination file was NOT deleted by Writer A!
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("destination was deleted by failed Writer A: %v", err)
	}
	if string(data) != "token from Writer B" {
		t.Fatalf("destination was corrupted by failed Writer A: %q", data)
	}
}

func TestCrossProcessTokenMutex(t *testing.T) {
	if os.Getenv("MESHERY_TEST_HELPER_MUTEX") == "1" {
		path := os.Getenv("MESHERY_TEST_MUTEX_PATH")
		unlock, err := acquireTokenFileLock(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "CHILD_ACQUIRE_LOCK_FAILED: %v\n", err)
			os.Exit(2)
		}
		// Signal parent by writing to stdout
		fmt.Println("LOCKED")
		// Hold lock until parent signals or stdin closes
		buf := make([]byte, 1)
		_, _ = os.Stdin.Read(buf)
		unlock()
		os.Exit(0)
		return
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "token")

	cmd := exec.Command(os.Args[0], "-test.run=TestCrossProcessTokenMutex")
	cmd.Env = append(os.Environ(),
		"MESHERY_TEST_HELPER_MUTEX=1",
		"MESHERY_TEST_MUTEX_PATH="+path,
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Skipf("cannot start child process: %v", err)
	}

	var stderrBuf strings.Builder
	stderrDone := make(chan struct{})
	go func() {
		defer close(stderrDone)
		_, _ = io.Copy(&stderrBuf, stderr)
	}()

	var waitOnce sync.Once
	var waitErr error
	waitForChild := func() error {
		waitOnce.Do(func() {
			_ = stdin.Close()
			waitErr = cmd.Wait()
			<-stderrDone
		})
		return waitErr
	}
	defer func() {
		_ = waitForChild()
	}()

	// Wait for child to acquire the named mutex; use io.ReadFull so a short
	// Read does not silently accept a partial "LOCKED" prefix. Apply a deadline
	// so a silent child failure does not hang the test suite.
	readyCh := make(chan string, 1)
	go func() {
		buf := make([]byte, 6)
		if _, err := io.ReadFull(stdout, buf); err == nil {
			readyCh <- string(buf)
		} else {
			readyCh <- ""
		}
	}()
	var readyStr string
	select {
	case readyStr = <-readyCh:
	case <-time.After(10 * time.Second):
		_ = cmd.Process.Kill()
		wErr := waitForChild()
		t.Fatalf("timed out waiting for child to signal LOCKED; waitErr=%v, stderr=%s", wErr, stderrBuf.String())
	}
	if readyStr != "LOCKED" {
		_ = cmd.Process.Kill()
		wErr := waitForChild()
		t.Fatalf("child did not signal LOCKED: got %q; exit status/waitErr=%v; stderr=%s", readyStr, wErr, stderrBuf.String())
	}

	// While child holds the cross-process mutex, verify parent fails to acquire it quickly
	mutexName, err := tokenFileMutexName(path)
	if err != nil {
		t.Fatal(err)
	}
	namePtr, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		t.Fatal(err)
	}
	handle, err := windows.OpenMutex(windows.SYNCHRONIZE, false, namePtr)
	if err != nil {
		t.Fatalf("open child mutex: %v", err)
	}
	defer windows.CloseHandle(handle)

	// Wait with 50ms timeout; must return WAIT_TIMEOUT because child holds it
	const waitTimeout = 0x00000102
	event, err := windows.WaitForSingleObject(handle, 50)
	if err != nil || event != waitTimeout {
		t.Fatalf("expected WAIT_TIMEOUT while child holds mutex, got event=%#x, err=%v", event, err)
	}

	// Release child and collect its exit result cleanly
	if err := waitForChild(); err != nil {
		t.Fatalf("child exited with error after release: %v; stderr=%s", err, stderrBuf.String())
	}

	// Now parent can acquire without timeout. Named mutexes have thread
	// affinity on Windows, so pin the goroutine to its OS thread for
	// WaitForSingleObject and ReleaseMutex.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	eventAfter, err := windows.WaitForSingleObject(handle, 2000)
	if err != nil || eventAfter != windows.WAIT_OBJECT_0 {
		t.Fatalf("expected WAIT_OBJECT_0 after child exit, got event=%#x, err=%v", eventAfter, err)
	}
	if err := windows.ReleaseMutex(handle); err != nil {
		t.Fatalf("ReleaseMutex: %v", err)
	}
}

func TestVerifyDestinationTokenFile_DetectsReplacement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")
	if err := WriteTokenFile(path, []byte("token")); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(windows.Handle(file.Fd()), &info); err != nil {
		t.Fatal(err)
	}

	tamperedInfo := info
	tamperedInfo.FileIndexLow++

	if err := verifyDestinationTokenFile(file, tamperedInfo, []byte("token")); err == nil {
		t.Fatal("expected verification to fail on file index mismatch")
	}
}

func TestVerifyDestinationTokenFile_DetectsContentMismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "token")
	if err := WriteTokenFile(path, []byte("actual token")); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(windows.Handle(file.Fd()), &info); err != nil {
		t.Fatal(err)
	}

	if err := verifyDestinationTokenFile(file, info, []byte("forged token")); err == nil {
		t.Fatal("expected verification to fail on content mismatch")
	}
}


// readDACL returns whether the file's DACL is protected from inheritance and
// the SIDs of the entries in it.
func readDACL(t *testing.T, path string) (protected bool, sids []*windows.SID) {
	t.Helper()

	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		t.Fatalf("GetNamedSecurityInfo: %v", err)
	}
	control, _, err := sd.Control()
	if err != nil {
		t.Fatalf("Control: %v", err)
	}
	dacl, _, err := sd.DACL()
	if err != nil {
		t.Fatalf("DACL: %v", err)
	}
	if dacl == nil {
		t.Fatal("file has a NULL DACL")
	}
	for i := uint32(0); i < uint32(dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(dacl, i, &ace); err != nil {
			t.Fatalf("GetAce(%d): %v", i, err)
		}
		sids = append(sids, (*windows.SID)(unsafe.Pointer(&ace.SidStart)))
	}
	return control&windows.SE_DACL_PROTECTED != 0, sids
}

func mustCurrentUserSID(t *testing.T) *windows.SID {
	t.Helper()
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		t.Fatalf("GetCurrentProcessToken().GetTokenUser: %v", err)
	}
	return user.User.Sid
}

func everyoneSID(t *testing.T) *windows.SID {
	t.Helper()
	sid, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	if err != nil {
		t.Fatalf("CreateWellKnownSid: %v", err)
	}
	sidPtr := (*windows.SID)(unsafe.Pointer(sid))
	return sidPtr
}

// grantEveryoneRead replaces the file's DACL with an unprotected one that
// gives the current user full access and Everyone read access, which is what a
// token file in a shared directory can look like.
func grantEveryoneRead(t *testing.T, path string) {
	t.Helper()

	ownerSID := mustCurrentUserSID(t)
	everyone := everyoneSID(t)

	var pinner runtime.Pinner
	pinner.Pin(ownerSID)
	pinner.Pin(everyone)
	defer pinner.Unpin()

	entry := func(sid *windows.SID, access windows.ACCESS_MASK) windows.EXPLICIT_ACCESS {
		return windows.EXPLICIT_ACCESS{
			AccessPermissions: access,
			AccessMode:        windows.GRANT_ACCESS,
			Inheritance:       windows.NO_INHERITANCE,
			Trustee: windows.TRUSTEE{
				TrusteeForm:  windows.TRUSTEE_IS_SID,
				TrusteeType:  windows.TRUSTEE_IS_UNKNOWN,
				TrusteeValue: windows.TrusteeValueFromSID(sid),
			},
		}
	}
	const standardTokenAccess = windows.STANDARD_RIGHTS_ALL | windows.FILE_GENERIC_READ | windows.FILE_GENERIC_WRITE
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{
		entry(ownerSID, standardTokenAccess),
		entry(everyone, windows.GENERIC_READ),
	}, nil)
	if err != nil {
		t.Fatalf("ACLFromEntries: %v", err)
	}
	if err := windows.SetNamedSecurityInfo(path, windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.UNPROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil); err != nil {
		t.Fatalf("SetNamedSecurityInfo: %v", err)
	}
}

func assertOnlyCurrentUser(t *testing.T, path string) {
	t.Helper()

	protected, sids := readDACL(t, path)
	if !protected {
		t.Error("token file DACL is not protected from inheritance")
	}
	if len(sids) != 1 {
		t.Fatalf("token file DACL has %d entries, want 1", len(sids))
	}
	if !windows.EqualSid(sids[0], mustCurrentUserSID(t)) {
		t.Errorf("token file DACL grants %s, want only the current user %s", sids[0], mustCurrentUserSID(t))
	}
}

func TestWriteTokenFileRestrictsNewFileToCurrentUser(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")

	if err := WriteTokenFile(path, []byte(`{"token":"t"}`)); err != nil {
		t.Fatalf("WriteTokenFile: %v", err)
	}

	assertOnlyCurrentUser(t, path)
	if got, err := os.ReadFile(path); err != nil || string(got) != `{"token":"t"}` {
		t.Fatalf("token file content = %q, err = %v", got, err)
	}
}

// The case from the issue: a token file that already grants other users access
// must lose that access when it is rewritten.
func TestWriteTokenFileRemovesExistingBroadAccess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(path, []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}
	grantEveryoneRead(t, path)

	// Guard the test: verify that the insecure file fails verifyTokenFileACL before WriteTokenFile.
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := verifyTokenFileACL(f, mustCurrentUserSID(t)); err == nil {
		f.Close()
		t.Fatal("setup failed: file is already owner-only before WriteTokenFile")
	}
	f.Close()

	if err := WriteTokenFile(path, []byte(`{"token":"new"}`)); err != nil {
		t.Fatalf("WriteTokenFile: %v", err)
	}

	assertOnlyCurrentUser(t, path)
	if got, err := os.ReadFile(path); err != nil || string(got) != `{"token":"new"}` {
		t.Fatalf("token file content = %q, err = %v", got, err)
	}
}

func TestWriteTokenFileRewriteKeepsRestriction(t *testing.T) {
	path := filepath.Join(t.TempDir(), "auth.json")

	for i := 0; i < 2; i++ {
		if err := WriteTokenFile(path, []byte(`{"token":"t"}`)); err != nil {
			t.Fatalf("WriteTokenFile #%d: %v", i+1, err)
		}
		assertOnlyCurrentUser(t, path)
	}
}
