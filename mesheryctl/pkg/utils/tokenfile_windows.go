//go:build windows

package utils

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

// tokenFileAccess is the access mask granted to the current user on the token
// file: full control of the file object (FILE_ALL_ACCESS), which is the
// Windows counterpart of "owner can read and write" for a file the owner also
// has to be able to replace or delete.
const tokenFileAccess = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0x1FF

// restrictTokenFile limits an open auth token file to the current user.
//
// Chmod(0600) is not enough on Windows: it only toggles the read-only
// attribute and leaves every ACE the file already has, including ACEs it
// inherited from its directory, in place. So the file's DACL is replaced with
// a single ACE that grants the current user access, and the DACL is marked
// protected so that ACEs from the parent directory are no longer inherited.
// The security descriptor is then read back and checked, so a call that
// returns nil means the restriction is in force rather than merely requested.
func restrictTokenFile(f *os.File) error {
	return restrictPathToCurrentUser(f.Name())
}

func restrictPathToCurrentUser(path string) error {
	sid, err := currentUserSID()
	if err != nil {
		return fmt.Errorf("cannot restrict token file %q: %w", path, err)
	}

	acl, err := windows.ACLFromEntries([]windows.EXPLICIT_ACCESS{{
		AccessPermissions: tokenFileAccess,
		AccessMode:        windows.GRANT_ACCESS,
		Inheritance:       windows.NO_INHERITANCE,
		Trustee: windows.TRUSTEE{
			TrusteeForm:  windows.TRUSTEE_IS_SID,
			TrusteeType:  windows.TRUSTEE_IS_USER,
			TrusteeValue: windows.TrusteeValueFromSID(sid),
		},
	}}, nil)
	if err != nil {
		return fmt.Errorf("cannot build access list for token file %q: %w", path, err)
	}

	if err := windows.SetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION|windows.PROTECTED_DACL_SECURITY_INFORMATION,
		nil, nil, acl, nil,
	); err != nil {
		return fmt.Errorf("cannot set owner-only access on token file %q: %w", path, err)
	}

	return verifyOwnerOnly(path, sid)
}

// currentUserSID returns the SID of the user the process runs as.
func currentUserSID() (*windows.SID, error) {
	// The process pseudo-token does not need to be closed.
	tokenUser, err := windows.GetCurrentProcessToken().GetTokenUser()
	if err != nil {
		return nil, fmt.Errorf("cannot look up current user: %w", err)
	}
	return tokenUser.User.Sid, nil
}

// sidOfACE returns the SID an access-allowed ACE applies to. The SID is stored
// inline in the ACE, starting at its SidStart field.
func sidOfACE(ace *windows.ACCESS_ALLOWED_ACE) *windows.SID {
	return (*windows.SID)(unsafe.Pointer(&ace.SidStart))
}

// verifyOwnerOnly reports an error unless the file's DACL is protected from
// inheritance and consists only of ACEs that allow access to sid, with at
// least one such ACE. A missing (NULL) DACL means everyone has full access, so
// it fails the check as well.
func verifyOwnerOnly(path string, sid *windows.SID) error {
	sd, err := windows.GetNamedSecurityInfo(path, windows.SE_FILE_OBJECT, windows.DACL_SECURITY_INFORMATION)
	if err != nil {
		return fmt.Errorf("cannot read access list of token file %q: %w", path, err)
	}

	control, _, err := sd.Control()
	if err != nil {
		return fmt.Errorf("cannot read security descriptor of token file %q: %w", path, err)
	}
	if control&windows.SE_DACL_PROTECTED == 0 {
		return fmt.Errorf("token file %q still inherits access from its directory", path)
	}

	dacl, _, err := sd.DACL()
	if err != nil {
		return fmt.Errorf("cannot read access list of token file %q: %w", path, err)
	}
	if dacl == nil {
		return fmt.Errorf("token file %q has no access list, so it is open to everyone", path)
	}

	granted := false
	for i := uint32(0); i < uint32(dacl.AceCount); i++ {
		var ace *windows.ACCESS_ALLOWED_ACE
		if err := windows.GetAce(dacl, i, &ace); err != nil {
			return fmt.Errorf("cannot read access entry %d of token file %q: %w", i, path, err)
		}
		if ace.Header.AceType != windows.ACCESS_ALLOWED_ACE_TYPE {
			return fmt.Errorf("token file %q has an unexpected access entry type %d", path, ace.Header.AceType)
		}
		aceSID := sidOfACE(ace)
		if !aceSID.Equals(sid) {
			return fmt.Errorf("token file %q grants access to %s besides the current user", path, aceSID.String())
		}
		granted = true
	}
	if !granted {
		return fmt.Errorf("token file %q does not grant access to the current user", path)
	}
	return nil
}
