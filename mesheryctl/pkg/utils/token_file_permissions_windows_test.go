//go:build windows

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

func TestWriteTokenFileRestrictsWindowsACL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "token")
	if err := os.WriteFile(path, []byte("old token"), 0600); err != nil {
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
	f, err := openTokenFile(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		t.Fatal(err)
	}
	if err := verifyTokenFileACL(f, user.User.Sid); err != nil {
		t.Fatalf("verify token ACL: %v", err)
	}
}
