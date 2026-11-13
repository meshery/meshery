//go:build windows

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

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
	longPath, err := tokenFileLongPath(path)
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
func writeTokenFile(path string, data []byte) error {
	var random [16]byte
	if _, err := rand.Read(random[:]); err != nil {
		return fmt.Errorf("generate temporary token filename: %w", err)
	}
	tempPath := filepath.Join(filepath.Dir(path), ".meshery-token-"+hex.EncodeToString(random[:]))
	f, err := createPrivateTokenFile(tempPath)
	if err != nil {
		return fmt.Errorf("create private token file: %w", err)
	}
	moved := false
	defer func() {
		_ = f.Close()
		if !moved {
			_ = os.Remove(tempPath)
		}
	}()
	if err := restrictTokenFileToCurrentUser(f); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	toPath, err := tokenFileLongPath(path)
	if err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := replaceTokenFile(f, toPath); err != nil {
		return fmt.Errorf("replace token file: %w", err)
	}
	moved = true
	return f.Close()
}

// replaceTokenFile uses POSIX rename semantics so a destination with a
// delete-sharing metadata handle can be replaced without reusing its object.
func replaceTokenFile(file *os.File, destination string) error {
	name, err := windows.UTF16FromString(destination)
	if err != nil {
		return err
	}
	name = name[:len(name)-1]
	type renameInfo struct {
		Flags          uint32
		RootDirectory  windows.Handle
		FileNameLength uint32
		FileName       [1]uint16
	}
	const replaceIfExists = 0x1
	const posixSemantics = 0x2
	var info renameInfo
	offset := int(unsafe.Offsetof(info.FileName))
	buffer := make([]byte, offset+len(name)*2)
	entry := (*renameInfo)(unsafe.Pointer(&buffer[0]))
	entry.Flags = replaceIfExists | posixSemantics
	entry.FileNameLength = uint32(len(name) * 2)
	copy(unsafe.Slice((*uint16)(unsafe.Pointer(&buffer[offset])), len(name)), name)
	return windows.SetFileInformationByHandle(windows.Handle(file.Fd()), windows.FileRenameInfoEx, &buffer[0], uint32(len(buffer)))
}

func tokenFileLongPath(path string) (string, error) {
	if strings.HasPrefix(path, `\\?\`) || strings.HasPrefix(path, `\??\`) {
		return path, nil
	}
	path, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if len(path) < 248 {
		return path, nil
	}
	if strings.HasPrefix(path, `\\`) {
		return `\\?\UNC\` + strings.TrimPrefix(path, `\\`), nil
	}
	return `\\?\` + path, nil
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
