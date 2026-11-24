package models

import (
	"encoding/json"
	"testing"

	userv1beta2 "github.com/meshery/schemas/models/v1beta2/user"
)

// TestAnonymousFlowResponseDecodesSchemasContract pins the anonymous user flow
// reply to the canonical schemas v1beta2 user.AnonymousFlowResponse construct.
//
// The reply is a CONSUMED contract: the remote provider produces it and Meshery
// Server only decodes it. A local re-declaration of the struct silently renamed
// the id field to `owner` while meshery-cloud kept emitting `userId`, so
// flowResponse's id decoded to the nil UUID on every anonymous sign-in and the
// session's capabilities were written under the zero UUID instead of the real
// user. This is that regression's guard - the payload below is byte-for-byte
// what meshery-cloud's anonymous handler encodes.
func TestAnonymousFlowResponseDecodesSchemasContract(t *testing.T) {
	const providerReply = `{
		"accessToken": "eyJhbGciOiJIUzI1NiJ9.e30.signature",
		"capability": {"capabilities": []},
		"userId": "0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c40"
	}`

	var got AnonymousFlowResponse
	if err := json.Unmarshal([]byte(providerReply), &got); err != nil {
		t.Fatalf("decode anonymous flow response: %v", err)
	}

	if got.AccessToken == "" {
		t.Error("accessToken did not decode")
	}
	if got.UserID.IsNil() {
		t.Fatal(`userId decoded to the nil UUID: the struct is not reading the canonical "userId" key`)
	}
	if want := "0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c40"; got.UserID.String() != want {
		t.Errorf("userId = %q, want %q", got.UserID, want)
	}
	if got.Capabilities == nil {
		t.Error(`capability did not decode`)
	}
}

// TestAnonymousFlowResponseRejectsLegacyOwnerKey pins the negative half of the
// contract: an `owner`-keyed reply is NOT the schemas construct, so it must
// leave the id nil rather than quietly resolving. Without this, reintroducing a
// local struct keyed on `owner` would pass the positive test above by accident.
func TestAnonymousFlowResponseRejectsLegacyOwnerKey(t *testing.T) {
	const legacyReply = `{
		"accessToken": "eyJhbGciOiJIUzI1NiJ9.e30.signature",
		"owner": "0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c40"
	}`

	var got AnonymousFlowResponse
	if err := json.Unmarshal([]byte(legacyReply), &got); err != nil {
		t.Fatalf("decode anonymous flow response: %v", err)
	}

	if !got.UserID.IsNil() {
		t.Errorf(`the legacy "owner" key resolved a user id (%v); the wire key is "userId"`, got.UserID)
	}
}

// TestAnonymousFlowResponseIsTheSchemasType is the compile-time half of the
// guard: AnonymousFlowResponse must be an alias for the schemas construct, not
// a structurally-similar local copy that can drift from it again.
func TestAnonymousFlowResponseIsTheSchemasType(t *testing.T) {
	var fromSchemas userv1beta2.AnonymousFlowResponse
	// The explicit type is the assertion: only an alias assigns here. Letting
	// staticcheck's ST1023 infer it from the right-hand side would leave
	// nothing for this test to check.
	var local AnonymousFlowResponse = fromSchemas //nolint:staticcheck
	_ = local
}
