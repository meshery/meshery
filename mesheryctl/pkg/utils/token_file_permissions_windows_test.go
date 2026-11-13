//go:build windows

package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestTokenFileOwnerRejectsDifferentSID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := WriteTokenFile(path, []byte("old token")); err != nil {
		t.Fatal(err)
	}
	file, err := openTokenFile(path)
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
