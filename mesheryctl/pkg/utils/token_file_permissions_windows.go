//go:build windows

package utils

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

type destLockEntry struct {
	mu  sync.Mutex
	ref int
}

var (
	destLocksMu              sync.Mutex
	destLocks                = make(map[string]*destLockEntry)
	acquireTokenFileLockFunc = acquireTokenFileLock
)

func canonicalTokenPath(path string) (string, error) {
	normalized := path
	switch {
	case strings.HasPrefix(normalized, `\??\UNC\`):
		normalized = `\\` + strings.TrimPrefix(normalized, `\??\UNC\`)
	case strings.HasPrefix(normalized, `\\?\UNC\`):
		normalized = `\\` + strings.TrimPrefix(normalized, `\\?\UNC\`)
	case strings.HasPrefix(normalized, `\??\`):
		normalized = strings.TrimPrefix(normalized, `\??\`)
	case strings.HasPrefix(normalized, `\\?\`):
		normalized = strings.TrimPrefix(normalized, `\\?\`)
	}

	abs, err := filepath.Abs(normalized)
	if err != nil {
		return "", err
	}
	return strings.ToLower(filepath.Clean(abs)), nil
}

func acquireDestLock(canonicalPath string) func() {
	destLocksMu.Lock()
	entry, ok := destLocks[canonicalPath]
	if !ok {
		entry = &destLockEntry{}
		destLocks[canonicalPath] = entry
	}
	entry.ref++
	destLocksMu.Unlock()

	entry.mu.Lock()
	return func() {
		entry.mu.Unlock()
		destLocksMu.Lock()
		entry.ref--
		if entry.ref == 0 {
			delete(destLocks, canonicalPath)
		}
		destLocksMu.Unlock()
	}
}

// tokenFileMutexName returns a session-local named mutex identifier derived from
// the canonical lowercase destination path.
func tokenFileMutexName(path string) (string, error) {
	norm, err := canonicalTokenPath(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte(norm))
	// "Local\" namespace ensures session-isolation for the current logon session.
	return fmt.Sprintf(`Local\meshery-token-%x`, sum[:16]), nil
}

// acquireTokenFileLock acquires both a per-destination process-local lock and a cross-process Windows
// named mutex scoped to the destination token file. The lock is held through the
// entire critical section (creation, write, permission hardening, installation,
// and failure cleanup).
func acquireTokenFileLock(path string) (func(), error) {
	norm, err := canonicalTokenPath(path)
	if err != nil {
		return nil, fmt.Errorf("canonicalize token path: %w", err)
	}

	releaseDestLock := acquireDestLock(norm)
	acquiredDestLock := true
	defer func() {
		if acquiredDestLock {
			releaseDestLock()
		}
	}()

	mutexName, err := tokenFileMutexName(path)
	if err != nil {
		return nil, fmt.Errorf("generate token mutex name: %w", err)
	}
	namePtr, err := windows.UTF16PtrFromString(mutexName)
	if err != nil {
		return nil, fmt.Errorf("convert token mutex name: %w", err)
	}

	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, fmt.Errorf("get current user SID: %w", err)
	}

	var pinner runtime.Pinner
	pinner.Pin(user.User.Sid)
	defer pinner.Unpin()

	sid := user.User.Sid.String()
	sd, err := windows.SecurityDescriptorFromString("O:" + sid + "D:P(A;;GA;;;" + sid + ")")
	if err != nil {
		return nil, fmt.Errorf("create token mutex security descriptor: %w", err)
	}

	sa := windows.SecurityAttributes{
		Length:             uint32(unsafe.Sizeof(windows.SecurityAttributes{})),
		SecurityDescriptor: sd,
	}

	// Create or open the session-local named mutex with restricted security descriptor.
	handle, err := windows.CreateMutex(&sa, false, namePtr)
	if err != nil && !errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		return nil, fmt.Errorf("create token mutex: %w", err)
	}

	// Lock the calling OS thread because Windows named mutexes have thread affinity.
	runtime.LockOSThread()
	acquiredOSThread := true
	defer func() {
		if acquiredOSThread {
			runtime.UnlockOSThread()
			_ = windows.CloseHandle(handle)
		}
	}()

	const timeoutMs = 15000 // 15 seconds timeout
	event, err := windows.WaitForSingleObject(handle, timeoutMs)
	if err != nil {
		return nil, fmt.Errorf("wait for token mutex: %w", err)
	}

	const waitTimeout = 0x00000102
	switch event {
	case windows.WAIT_OBJECT_0, windows.WAIT_ABANDONED:
		// Mutex ownership acquired. WAIT_ABANDONED handles a previous crashed holder safely.
	case waitTimeout:
		return nil, fmt.Errorf("timeout waiting for token file lock on %s", path)
	default:
		return nil, fmt.Errorf("unexpected event %#x waiting for token mutex on %s", event, path)
	}

	acquiredDestLock = false
	acquiredOSThread = false

	unlock := func() {
		_ = windows.ReleaseMutex(handle)
		_ = windows.CloseHandle(handle)
		runtime.UnlockOSThread()
		releaseDestLock()
	}
	return unlock, nil
}

// createPrivateTokenFile creates a new file with a private DACL before any token is written.
// Permission helpers return wrapped system errors for callers to add command context.
func createPrivateTokenFile(path string) (*os.File, error) {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, fmt.Errorf("get current user SID: %w", err)
	}
	sid := user.User.Sid.String()
	sd, err := windows.SecurityDescriptorFromString("O:" + sid + "D:P(A;;FA;;;" + sid + ")")
	if err != nil {
		return nil, fmt.Errorf("create token file security descriptor: %w", err)
	}
	longPath, err := tokenFileWin32Path(path)
	if err != nil {
		return nil, err
	}
	name, err := windows.UTF16PtrFromString(longPath)
	if err != nil {
		return nil, err
	}
	sa := windows.SecurityAttributes{Length: uint32(unsafe.Sizeof(windows.SecurityAttributes{})), SecurityDescriptor: sd}
	handle, err := windows.CreateFile(name, windows.GENERIC_WRITE|windows.WRITE_DAC|windows.READ_CONTROL|windows.DELETE,
		0, &sa, windows.CREATE_NEW, windows.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(handle), path), nil
}

// writeTokenFile writes a newly created private file and atomically replaces
// the destination. Existing handles retain access only to the old file.
//
// Source identity is strictly bound to the verified open handle f.Fd()
// via handle-based atomic rename operations, eliminating path-replacement races.
// Concurrent writers are synchronized across the entire critical section,
// and failed writers clean up only their own temporary file, never deleting
// the destination.
func writeTokenFile(path string, data []byte) error {
	unlock, err := acquireTokenFileLockFunc(path)
	if err != nil {
		return ErrAcquireTokenFileLock(err)
	}
	defer unlock()

	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return fmt.Errorf("generate temporary token filename: %w", err)
	}
	tempPath := filepath.Join(filepath.Dir(path), ".meshery-token-"+hex.EncodeToString(random[:]))
	f, err := createPrivateTokenFile(tempPath)
	if err != nil {
		return fmt.Errorf("create private token file: %w", err)
	}
	installed := false
	defer func() {
		_ = f.Close()
		if !installed {
			// Cleanup is strictly tied to the current operation's temporary file.
			// Destination is never deleted on failure.
			_ = os.Remove(tempPath)
		}
	}()

	if err := restrictTokenFileToCurrentUser(f); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}

	toPath, err := tokenFileLongPath(path)
	if err != nil {
		return err
	}
	if err := replaceTokenFileFunc(f, toPath); err != nil {
		// Some filesystems reject the extended rename operation. The legacy
		// rename still replaces the file object and fails closed if an open
		// handle prevents replacement; it never writes into the old file.
		if !errors.Is(err, windows.ERROR_INVALID_NAME) &&
			!errors.Is(err, windows.ERROR_INVALID_PARAMETER) &&
			!errors.Is(err, windows.ERROR_NOT_SUPPORTED) &&
			!errors.Is(err, windows.ERROR_INVALID_FUNCTION) &&
			!errors.Is(err, windows.ERROR_CALL_NOT_IMPLEMENTED) {
			return ErrReplaceTokenFile(err)
		}
		var beforeInfo windows.ByHandleFileInformation
		if infoErr := windows.GetFileInformationByHandle(windows.Handle(f.Fd()), &beforeInfo); infoErr != nil {
			return replacementFallbackError(err, "record file info", infoErr)
		}
		if beforeInfo.VolumeSerialNumber == 0 || (beforeInfo.FileIndexHigh == 0 && beforeInfo.FileIndexLow == 0) {
			return replacementFallbackError(err, "file identity cannot be verified on this filesystem", nil)
		}
		if closeErr := f.Close(); closeErr != nil {
			return replacementFallbackError(err, "close temporary file", closeErr)
		}
		fromPath, pathErr := tokenFileWin32Path(tempPath)
		if pathErr != nil {
			return replacementFallbackError(err, "resolve source path", pathErr)
		}
		from, pathErr := windows.UTF16PtrFromString(fromPath)
		if pathErr != nil {
			return replacementFallbackError(err, "convert source path", pathErr)
		}
		win32ToPath, pathErr := tokenFileWin32Path(path)
		if pathErr != nil {
			return replacementFallbackError(err, "resolve destination path", pathErr)
		}
		to, pathErr := windows.UTF16PtrFromString(win32ToPath)
		if pathErr != nil {
			return replacementFallbackError(err, "convert destination path", pathErr)
		}
		if moveErr := windows.MoveFileEx(from, to, windows.MOVEFILE_REPLACE_EXISTING|windows.MOVEFILE_WRITE_THROUGH); moveErr != nil {
			return replacementFallbackError(err, "fallback", moveErr)
		}
		destHandle, openErr := windows.CreateFile(
			to,
			windows.GENERIC_READ|windows.READ_CONTROL,
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
			nil,
			windows.OPEN_EXISTING,
			windows.FILE_ATTRIBUTE_NORMAL,
			0,
		)
		if openErr != nil {
			_ = windows.DeleteFile(to)
			return replacementFallbackError(err, "verify destination file", openErr)
		}
		destFile := os.NewFile(uintptr(destHandle), toPath)
		if verifyErr := verifyDestinationTokenFile(destFile, beforeInfo, data); verifyErr != nil {
			_ = destFile.Close()
			_ = windows.DeleteFile(to)
			return replacementFallbackError(err, "verify destination file", verifyErr)
		}
		if closeErr := destFile.Close(); closeErr != nil {
			_ = windows.DeleteFile(to)
			return replacementFallbackError(err, "close destination file", closeErr)
		}
		installed = true
		return nil
	}

	installed = true
	return f.Close()
}

var (
	replaceTokenFileExFunc  = replaceTokenFileEx
	replaceTokenFileStdFunc = replaceTokenFileStandard
)
var replaceTokenFileFunc = replaceTokenFile

func replacementFallbackError(extended error, operation string, operationErr error) error {
	cause := errors.Join(extended, operationErr)
	return ErrReplaceTokenFile(fmt.Errorf(
		"replace token file: extended rename: %v; %s: %w",
		extended,
		operation,
		cause,
	))
}

// verifyDestinationTokenFile verifies that the destination file after a MoveFileEx fallback
// is the exact file that was created, written, and restricted.
func verifyDestinationTokenFile(destFile *os.File, expectedInfo windows.ByHandleFileInformation, expectedData []byte) error {
	var actualInfo windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(windows.Handle(destFile.Fd()), &actualInfo); err != nil {
		return fmt.Errorf("read destination file info: %w", err)
	}
	if !sameWindowsFileIdentity(expectedInfo, actualInfo) {
		return fmt.Errorf("token file identity mismatch after move: expected (%x:%x:%x), got (%x:%x:%x)",
			expectedInfo.VolumeSerialNumber, expectedInfo.FileIndexHigh, expectedInfo.FileIndexLow,
			actualInfo.VolumeSerialNumber, actualInfo.FileIndexHigh, actualInfo.FileIndexLow)
	}
	if actualInfo.FileSizeHigh != expectedInfo.FileSizeHigh || actualInfo.FileSizeLow != expectedInfo.FileSizeLow {
		return fmt.Errorf("token file size mismatch after move: expected %d bytes, got %d bytes",
			uint64(expectedInfo.FileSizeHigh)<<32|uint64(expectedInfo.FileSizeLow),
			uint64(actualInfo.FileSizeHigh)<<32|uint64(actualInfo.FileSizeLow))
	}
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return fmt.Errorf("get current user SID: %w", err)
	}
	if err := verifyTokenFileOwner(destFile, user.User.Sid); err != nil {
		return fmt.Errorf("verify destination owner: %w", err)
	}
	if err := verifyTokenFileACL(destFile, user.User.Sid); err != nil {
		return fmt.Errorf("verify destination ACL: %w", err)
	}
	if len(expectedData) > 0 {
		content := make([]byte, len(expectedData))
		var bytesRead uint32
		if err := windows.ReadFile(windows.Handle(destFile.Fd()), content, &bytesRead, nil); err != nil {
			return fmt.Errorf("read destination content: %w", err)
		}
		if int(bytesRead) != len(expectedData) || !bytes.Equal(content, expectedData) {
			return fmt.Errorf("destination content does not match expected token")
		}
	}
	return nil
}

// sameWindowsFileIdentity compares the unique volume and file index identifiers.
// If the underlying filesystem does not support unique file indexing (all zero),
// it returns false to fail closed.
func sameWindowsFileIdentity(a, b windows.ByHandleFileInformation) bool {
	if a.VolumeSerialNumber == 0 || (a.FileIndexHigh == 0 && a.FileIndexLow == 0) {
		return false
	}
	return a.VolumeSerialNumber == b.VolumeSerialNumber &&
		a.FileIndexHigh == b.FileIndexHigh &&
		a.FileIndexLow == b.FileIndexLow
}

// renameFallbackError captures both the extended rename error and the standard
// rename error when both fail, ensuring neither native error is lost from the error chain.
type renameFallbackError struct {
	extended error
	standard error
}

func (e *renameFallbackError) Error() string {
	return fmt.Sprintf("handle-based rename: extended: %v; standard: %v", e.extended, e.standard)
}

func (e *renameFallbackError) Unwrap() []error {
	return []error{e.extended, e.standard}
}

// replaceTokenFile uses handle-based atomic rename operations to replace the
// destination. The source file is never closed before replacement and its identity
// is bound strictly to file.Fd().
func replaceTokenFile(file *os.File, destination string) error {
	name, err := windows.UTF16FromString(destination)
	if err != nil {
		return err
	}
	name = name[:len(name)-1]

	// 1. Try FileRenameInfoEx (class 22) with POSIX rename semantics first.
	// This replaces an existing destination even if a metadata handle is open with delete-sharing.
	errEx := replaceTokenFileExFunc(file, name)
	if errEx == nil {
		return nil
	}

	// 2. If class 22 is unsupported by the filesystem or Windows version,
	// fall back to standard handle-based rename (FileRenameInfo, class 3).
	// This also operates on the open handle file.Fd(), preserving source identity.
	if isUnsupportedRenameClassError(errEx) {
		errStd := replaceTokenFileStdFunc(file, name)
		if errStd == nil {
			return nil
		}
		if errors.Is(errStd, windows.ERROR_ACCESS_DENIED) {
			return errStd
		}
		return &renameFallbackError{extended: errEx, standard: errStd}
	}

	// For all other errors (such as ERROR_ACCESS_DENIED or ERROR_SHARING_VIOLATION),
	// fail closed immediately without weakening security.
	return errEx
}

func isUnsupportedRenameClassError(err error) bool {
	return errors.Is(err, windows.ERROR_INVALID_PARAMETER) ||
		errors.Is(err, windows.ERROR_INVALID_FUNCTION) ||
		errors.Is(err, windows.ERROR_NOT_SUPPORTED) ||
		errors.Is(err, windows.ERROR_CALL_NOT_IMPLEMENTED) ||
		// Some filesystems (e.g. FAT32, SMB shares) return ERROR_INVALID_NAME
		// when FileRenameInfoEx is passed the POSIX-semantics flag they do not
		// support. The standard rename (FileRenameInfo, class 3) is a safe
		// fallback: if the path itself is malformed it will also fail there,
		// so this cannot mask a genuine path error.
		errors.Is(err, windows.ERROR_INVALID_NAME)
}

func replaceTokenFileEx(file *os.File, name []uint16) error {
	type renameInfoEx struct {
		Flags          uint32
		RootDirectory  windows.Handle
		FileNameLength uint32
		FileName       [1]uint16
	}
	const replaceIfExists = 0x1
	const posixSemantics = 0x2
	var info renameInfoEx
	offset := int(unsafe.Offsetof(info.FileName))
	buffer := make([]byte, offset+len(name)*2)
	entry := (*renameInfoEx)(unsafe.Pointer(&buffer[0]))
	entry.Flags = replaceIfExists | posixSemantics
	entry.FileNameLength = uint32(len(name) * 2)
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(&buffer[offset])), len(name)), name)
	return windows.SetFileInformationByHandle(windows.Handle(file.Fd()), windows.FileRenameInfoEx, &buffer[0], uint32(len(buffer)))
}

func replaceTokenFileStandard(file *os.File, name []uint16) error {
	type renameInfoStd struct {
		ReplaceIfExists uint32
		RootDirectory   windows.Handle
		FileNameLength  uint32
		FileName        [1]uint16
	}
	var info renameInfoStd
	offset := int(unsafe.Offsetof(info.FileName))
	buffer := make([]byte, offset+len(name)*2)
	entry := (*renameInfoStd)(unsafe.Pointer(&buffer[0]))
	entry.ReplaceIfExists = 1
	entry.FileNameLength = uint32(len(name) * 2)
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(&buffer[offset])), len(name)), name)
	return windows.SetFileInformationByHandle(windows.Handle(file.Fd()), windows.FileRenameInfo, &buffer[0], uint32(len(buffer)))
}

func tokenFileLongPath(path string) (string, error) {
	// SetFileInformationByHandle (FileRenameInfo / FileRenameInfoEx) requires the
	// destination in the NT object-namespace format (\??\...) rather than the
	// Win32 format used by most user-mode APIs. Using a plain Win32 path such as
	// C:\path\to\file causes ERROR_INVALID_NAME on Windows.
	//
	// Mapping rules:
	//   \??\ prefix    → already NT namespace, return as-is
	//   \\?\UNC\...    → Win32 extended UNC, convert to NT \??\UNC\
	//   \\?\...        → Win32 extended local, convert to NT \??\
	//   \\server\share → UNC path, convert to \??\UNC\server\share
	//   C:\...         → local drive path, convert to \??\C:\...
	if strings.HasPrefix(path, `\??\`) {
		return path, nil
	}
	if strings.HasPrefix(path, `\\?\UNC\`) {
		return `\??\UNC\` + strings.TrimPrefix(path, `\\?\UNC\`), nil
	}
	if strings.HasPrefix(path, `\\?\`) {
		return `\??\` + path[4:], nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(abs, `\\`) {
		// UNC path: \\server\share\... → \??\UNC\server\share\...
		return `\??\UNC\` + strings.TrimPrefix(abs, `\\`), nil
	}
	// Local drive path: C:\... → \??\C:\...
	return `\??\` + abs, nil
}

// tokenFileWin32Path returns the Win32-extended path suitable for CreateFile.
// Short paths (< 248 chars) are returned unchanged; longer paths are prefixed
// with \\?\ to bypass the MAX_PATH limit.
func tokenFileWin32Path(path string) (string, error) {
	if strings.HasPrefix(path, `\\?\`) {
		return path, nil
	}
	if strings.HasPrefix(path, `\??\UNC\`) {
		return `\\?\UNC\` + strings.TrimPrefix(path, `\??\UNC\`), nil
	}
	if strings.HasPrefix(path, `\??\`) {
		return `\\?\` + strings.TrimPrefix(path, `\??\`), nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if len(abs) < 248 {
		return abs, nil
	}
	if strings.HasPrefix(abs, `\\`) {
		return `\\?\UNC\` + strings.TrimPrefix(abs, `\\`), nil
	}
	return `\\?\` + abs, nil
}

// restrictTokenFileToCurrentUser replaces inherited ACL entries with an ACL
// granting full access only to the account running mesheryctl.
func restrictTokenFileToCurrentUser(file *os.File) error {
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return fmt.Errorf("get current user SID: %w", err)
	}
	if err := verifyTokenFileOwner(file, user.User.Sid); err != nil {
		return err
	}

	var pinner runtime.Pinner
	pinner.Pin(user.User.Sid)
	defer pinner.Unpin()

	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: windows.GENERIC_ALL,
		AccessMode:        windows.GRANT_ACCESS,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(user.User.Sid),
		},
	}}, nil)
	if err != nil {
		return fmt.Errorf("create token file ACL: %w", err)
	}

	if err := windows.SetSecurityInfo(
		windows.Handle(file.Fd()),
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil,
	); err != nil {
		return fmt.Errorf("restrict token file ACL: %w", err)
	}
	return verifyTokenFileACL(file, user.User.Sid)
}

func verifyTokenFileOwner(file *os.File, sid *windows.SID) error {
	descriptor, err := windows.GetSecurityInfo(windows.Handle(file.Fd()), windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION)
	if err != nil {
		return fmt.Errorf("read token file owner: %w", err)
	}
	owner, _, err := descriptor.Owner()
	if err != nil {
		return fmt.Errorf("read token file owner: %w", err)
	}
	if owner == nil {
		return fmt.Errorf("token file has no owner")
	}
	if windows.EqualSid(owner, sid) {
		return nil
	}

	// Elevated Windows tokens may use the built-in Administrators group as
	// their default owner. Do not accept any other non-user owner.
	var size uint32
	token := windows.GetCurrentProcessToken()
	if err := windows.GetTokenInformation(token, windows.TokenOwner, nil, 0, &size); err != windows.ERROR_INSUFFICIENT_BUFFER {
		return fmt.Errorf("query token default owner: %w", err)
	}
	if size == 0 {
		return fmt.Errorf("token default owner is empty")
	}
	buf := make([]byte, size)
	if err := windows.GetTokenInformation(token, windows.TokenOwner, &buf[0], size, &size); err != nil {
		return fmt.Errorf("read token default owner: %w", err)
	}
	defaultOwner := *(**windows.SID)(unsafe.Pointer(&buf[0]))
	administrators, err := windows.CreateWellKnownSid(windows.WinBuiltinAdministratorsSid)
	if err != nil {
		return fmt.Errorf("get Administrators SID: %w", err)
	}
	trusted := isTrustedTokenOwner(owner, sid, defaultOwner, administrators)
	runtime.KeepAlive(buf)
	if !trusted {
		return fmt.Errorf("token file is not owned by the current user")
	}
	return nil
}

func isTrustedTokenOwner(owner, user, defaultOwner, administrators *windows.SID) bool {
	return windows.EqualSid(owner, user) || (defaultOwner != nil &&
		windows.EqualSid(owner, defaultOwner) && windows.EqualSid(owner, administrators))
}

func verifyTokenFileACL(file *os.File, sid *windows.SID) error {
	actual, err := windows.GetSecurityInfo(windows.Handle(file.Fd()), windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return fmt.Errorf("read token file ACL: %w", err)
	}
	control, _, err := actual.Control()
	if err != nil {
		return fmt.Errorf("read token file ACL control: %w", err)
	}
	dacl, _, err := actual.DACL()
	if err != nil || dacl == nil || dacl.AceCount != 1 || control&windows.SE_DACL_PROTECTED == 0 {
		return fmt.Errorf("token file ACL is not restricted to the current user")
	}
	var ace *windows.ACCESS_ALLOWED_ACE
	if err := windows.GetAce(dacl, 0, &ace); err != nil {
		return fmt.Errorf("read token file ACL entry: %w", err)
	}
	if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE || !windows.EqualSid((*windows.SID)(unsafe.Pointer(&ace.SidStart)), sid) {
		return fmt.Errorf("token file ACL grants access to an unexpected account")
	}
	const fileAllAccess = 0x001f01ff
	if ace.Mask != windows.GENERIC_ALL && ace.Mask != fileAllAccess {
		return fmt.Errorf("token file ACL does not grant the current user full access")
	}
	return nil
}
