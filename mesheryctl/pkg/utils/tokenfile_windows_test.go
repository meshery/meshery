//go:build windows

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/sys/windows"
)

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
		sids = append(sids, sidOfACE(ace))
	}
	return control&windows.SE_DACL_PROTECTED != 0, sids
}

func mustCurrentUserSID(t *testing.T) *windows.SID {
	t.Helper()
	sid, err := currentUserSID()
	if err != nil {
		t.Fatalf("currentUserSID: %v", err)
	}
	return sid
}

func everyoneSID(t *testing.T) *windows.SID {
	t.Helper()
	sid, err := windows.CreateWellKnownSid(windows.WinWorldSid)
	if err != nil {
		t.Fatalf("CreateWellKnownSid: %v", err)
	}
	return sid
}

// grantEveryoneRead replaces the file's DACL with an unprotected one that
// gives the current user full access and Everyone read access, which is what a
// token file in a shared directory can look like.
func grantEveryoneRead(t *testing.T, path string) {
	t.Helper()

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
	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{
		entry(mustCurrentUserSID(t), tokenFileAccess),
		entry(everyoneSID(t), windows.GENERIC_READ),
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
	if !sids[0].Equals(mustCurrentUserSID(t)) {
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

	// Guard the test: the setup must really have left Everyone with access.
	if err := verifyOwnerOnly(path, mustCurrentUserSID(t)); err == nil {
		t.Fatal("setup failed: file is already owner-only before WriteTokenFile")
	}

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

// verifyOwnerOnly is what lets WriteTokenFile report success, so it must
// reject every state that is not owner-only.
func TestVerifyOwnerOnlyRejectsOpenFiles(t *testing.T) {
	sid := mustCurrentUserSID(t)

	t.Run("inherits access from directory", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")
		if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := verifyOwnerOnly(path, sid); err == nil {
			t.Fatal("verifyOwnerOnly accepted a file with an inherited DACL")
		}
	})

	t.Run("grants another user access", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")
		if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		grantEveryoneRead(t, path)
		if err := verifyOwnerOnly(path, sid); err == nil {
			t.Fatal("verifyOwnerOnly accepted a file that Everyone can read")
		}
	})

	t.Run("owner-only file passes", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "auth.json")
		if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := restrictPathToCurrentUser(path); err != nil {
			t.Fatalf("restrictPathToCurrentUser: %v", err)
		}
		if err := verifyOwnerOnly(path, sid); err != nil {
			t.Fatalf("verifyOwnerOnly rejected an owner-only file: %v", err)
		}
	})
}

func TestRestrictPathToCurrentUserMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	if err := restrictPathToCurrentUser(path); err == nil {
		t.Fatal("restrictPathToCurrentUser succeeded for a missing file")
	}
}
