package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
)

// The four shapes below are the ones read-only production inspection found (see
// credential_secret_test.go for the counts). Encryption converges them onto one
// envelope, so every one of them has to survive a round trip unchanged - not
// "equivalently", unchanged - or a credential that worked before this shipped
// resolves differently after it.
func credentialShapeFixtures() []struct {
	name   string
	secret map[string]interface{}
} {
	return []struct {
		name   string
		secret map[string]interface{}
	}{
		{
			name:   "canonical prometheus",
			secret: map[string]interface{}{"prometheusURL": "https://prom.example"},
		},
		{
			name:   "canonical grafana",
			secret: map[string]interface{}{"grafanaAPIKey": "gr@f4n4-key", "grafanaURL": "https://grafana.example"},
		},
		{
			name: "kubernetes",
			secret: map[string]interface{}{
				"auth":    map[string]interface{}{"clusterToken": "k8s-service-account-token"},
				"cluster": map[string]interface{}{"server": "https://k8s.example"},
			},
		},
		{
			name: "legacy nested",
			secret: map[string]interface{}{
				"credentialName": "kube-cred",
				"secret": map[string]interface{}{
					"auth":    map[string]interface{}{"clusterToken": "nested-token"},
					"cluster": map[string]interface{}{"server": "https://k8s.example"},
				},
			},
		},
		{
			name:   "legacy string",
			secret: map[string]interface{}{"secret": "bare-token"},
		},
		{
			name:   "legacy string under the registration path's wrapper",
			secret: map[string]interface{}{"name": "prom-cred", "secret": "bare-token"},
		},
		{
			name:   "wrapper whose payload itself carries a string secret",
			secret: map[string]interface{}{"secret": map[string]interface{}{"secret": "inner-token"}},
		},
	}
}

// TestEncryptDecryptCredentialSecretRoundTrip is the compatibility guarantee:
// every production shape decrypts back to exactly what was encrypted, and both
// resolution helpers answer identically on either side of the trip.
func TestEncryptDecryptCredentialSecretRoundTrip(t *testing.T) {
	for _, tt := range credentialShapeFixtures() {
		t.Run(tt.name, func(t *testing.T) {
			wantPayload := CredentialPayload(tt.secret)
			wantAuthSecret := CredentialAuthSecret(tt.secret)

			encrypted, err := EncryptCredentialSecret(tt.secret)
			if err != nil {
				t.Fatalf("EncryptCredentialSecret: %v", err)
			}

			// The datastore round trip is JSON (core.Map marshals on write and
			// unmarshals on read), so run the fixture through it rather than
			// handing the in-memory map straight back to the decrypter.
			decrypted, err := DecryptCredentialSecret(throughJSON(t, encrypted))
			if err != nil {
				t.Fatalf("DecryptCredentialSecret: %v", err)
			}

			if !reflect.DeepEqual(decrypted, tt.secret) {
				t.Errorf("round trip changed the secret:\n got %#v\nwant %#v", decrypted, tt.secret)
			}
			if got := CredentialPayload(decrypted); !reflect.DeepEqual(got, wantPayload) {
				t.Errorf("CredentialPayload after round trip = %#v, want %#v", got, wantPayload)
			}
			if got := CredentialAuthSecret(decrypted); got != wantAuthSecret {
				t.Errorf("CredentialAuthSecret after round trip = %q, want %q", got, wantAuthSecret)
			}
		})
	}
}

// TestEncryptCredentialSecretHidesSecretMaterial is the point of the exercise:
// nothing recognisable as the credential's auth material survives into the
// bytes that reach the datastore.
func TestEncryptCredentialSecretHidesSecretMaterial(t *testing.T) {
	secrets := map[string][]string{
		"canonical grafana": {"gr@f4n4-key"},
		"kubernetes":        {"k8s-service-account-token"},
		"legacy nested":     {"nested-token"},
		"legacy string":     {"bare-token"},
	}

	for _, tt := range credentialShapeFixtures() {
		needles, ok := secrets[tt.name]
		if !ok {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := EncryptCredentialSecret(tt.secret)
			if err != nil {
				t.Fatalf("EncryptCredentialSecret: %v", err)
			}
			stored, err := json.Marshal(encrypted)
			if err != nil {
				t.Fatalf("marshal encrypted secret: %v", err)
			}
			for _, needle := range needles {
				if strings.Contains(string(stored), needle) {
					t.Errorf("stored credential still contains %q in the clear: %s", needle, stored)
				}
			}
		})
	}
}

// TestEncryptCredentialSecretKeepsWrapperNameInTheClear pins the one thing that
// deliberately is not encrypted. The legacy wrapper's `credentialName`/`name` is
// the credential's display name, already stored in the clear in the credentials
// table's own `name` column; hiding it here would buy nothing and would make
// the wrapper unreconstructable.
func TestEncryptCredentialSecretKeepsWrapperNameInTheClear(t *testing.T) {
	encrypted, err := EncryptCredentialSecret(map[string]interface{}{
		"credentialName": "kube-cred",
		"secret":         map[string]interface{}{"auth": map[string]interface{}{"clusterToken": "tok"}},
	})
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}
	if got := encrypted["credentialName"]; got != "kube-cred" {
		t.Errorf("credentialName = %v, want kube-cred", got)
	}
	if _, ok := encrypted["secret"]; ok {
		t.Error("the wrapper's `secret` field is still present; the payload was not moved into the envelope")
	}
}

// TestCredentialEncryptionDetectorIsUnambiguous proves the read path can always
// tell an encrypted secret from a plaintext one without attempting decryption.
// This is what lets a genuine decryption failure be reported as a failure rather
// than mistaken for a plaintext row.
func TestCredentialEncryptionDetectorIsUnambiguous(t *testing.T) {
	t.Run("no production shape is mistaken for ciphertext", func(t *testing.T) {
		for _, tt := range credentialShapeFixtures() {
			if IsEncryptedCredentialSecret(tt.secret) {
				t.Errorf("%s: plaintext secret detected as encrypted", tt.name)
			}
		}
	})

	t.Run("hostile look-alikes are not mistaken for ciphertext", func(t *testing.T) {
		lookAlikes := []struct {
			name   string
			secret map[string]interface{}
		}{
			{"nil", nil},
			{"empty", map[string]interface{}{}},
			{"a field merely named like the marker", map[string]interface{}{"mesheryEncryptedSecret": "meshery.enc.v1:deadbeef:AAAA"}},
			{"the envelope prefix stored under an ordinary key", map[string]interface{}{"secret": "meshery.enc.v1:deadbeef:AAAA"}},
			{"the reserved key holding a non-envelope string", map[string]interface{}{credentialCiphertextKey: "not-an-envelope"}},
			{"the reserved key holding a non-string", map[string]interface{}{credentialCiphertextKey: map[string]interface{}{}}},
		}
		for _, tt := range lookAlikes {
			if IsEncryptedCredentialSecret(tt.secret) {
				t.Errorf("%s: detected as encrypted", tt.name)
			}
		}
	})

	t.Run("every encrypted shape is detected as ciphertext", func(t *testing.T) {
		for _, tt := range credentialShapeFixtures() {
			encrypted, err := EncryptCredentialSecret(tt.secret)
			if err != nil {
				t.Fatalf("%s: EncryptCredentialSecret: %v", tt.name, err)
			}
			if !IsEncryptedCredentialSecret(encrypted) {
				t.Errorf("%s: encrypted secret not detected as encrypted", tt.name)
			}
			if !IsEncryptedCredentialSecret(throughJSON(t, encrypted)) {
				t.Errorf("%s: encrypted secret not detected after a datastore round trip", tt.name)
			}
		}
	})
}

// TestDecryptCredentialSecretPassesPlaintextThrough is the upgrade path: rows
// written before credential encryption shipped carry no envelope and must read
// back byte-for-byte, with no migration and no user action.
func TestDecryptCredentialSecretPassesPlaintextThrough(t *testing.T) {
	for _, tt := range credentialShapeFixtures() {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecryptCredentialSecret(tt.secret)
			if err != nil {
				t.Fatalf("DecryptCredentialSecret on a plaintext secret: %v", err)
			}
			if !reflect.DeepEqual(got, tt.secret) {
				t.Errorf("plaintext secret changed on read:\n got %#v\nwant %#v", got, tt.secret)
			}
		})
	}

	t.Run("nil", func(t *testing.T) {
		got, err := DecryptCredentialSecret(nil)
		if err != nil {
			t.Fatalf("DecryptCredentialSecret(nil): %v", err)
		}
		if got != nil {
			t.Errorf("DecryptCredentialSecret(nil) = %#v, want nil", got)
		}
	})
}

// TestEncryptCredentialSecretRejectsTheReservedProperty pins the write path's
// only defensible answer to a secret that already carries Meshery's ciphertext
// marker. Every read path returns plaintext, so no internal caller can produce
// that shape; the one thing that can is a request body. Passing it through -
// treating the marker as proof the map is already sealed - would store whatever
// else the map carries in the clear and leave the row unreadable forever.
func TestEncryptCredentialSecretRejectsTheReservedProperty(t *testing.T) {
	sealed, err := EncryptCredentialSecret(map[string]interface{}{"grafanaAPIKey": "key"})
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}

	for name, secret := range map[string]map[string]interface{}{
		"a real envelope Meshery produced": sealed,
		"a forged envelope beside a live secret": {
			credentialCiphertextKey: credentialEnvelopePrefix + "aa:bb",
			"apiKey":                "real-token",
		},
		"the marker holding a non-envelope": {
			credentialCiphertextKey: 42,
			"apiKey":                "real-token",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got, err := EncryptCredentialSecret(secret)
			if err == nil {
				t.Fatalf("a secret carrying the reserved property was accepted: %#v", got)
			}
			if got != nil {
				t.Errorf("rejection returned a secret to persist: %#v", got)
			}
			if code := meshkiterrors.GetCode(err); code != ErrReservedCredentialPropertyCode {
				t.Errorf("error code = %s, want %s", code, ErrReservedCredentialPropertyCode)
			}
		})
	}
}

// TestSaveUserCredentialRejectsAReservedPropertyFromTheRequestBody is the same
// guarantee at the layer that is actually reachable: the handler unmarshals a
// request body straight into Credential.Secret, so a client controls that map
// in full. Nothing may be persisted for such a request.
func TestSaveUserCredentialRejectsAReservedPropertyFromTheRequestBody(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	// Byte-for-byte what handlers.SaveUserCredential does with a request body,
	// including the empty-map initialisation it performs before unmarshalling.
	body := []byte(`{"name":"hostile","type":"token","secret":{"__mesheryEncryptedSecret":"meshery.enc.v1:aa:bb","apiKey":"real-token"}}`)
	credential := Credential{Secret: core.Map{}}
	if err := json.Unmarshal(body, &credential); err != nil {
		t.Fatalf("unmarshalling the request body: %v", err)
	}
	credential.UserId = userID

	if _, err := provider.SaveUserCredential("tok", &credential); err == nil {
		t.Fatal("a credential carrying the reserved property was saved")
	} else if code := meshkiterrors.GetCode(err); code != ErrReservedCredentialPropertyCode {
		t.Errorf("error code = %s, want %s", code, ErrReservedCredentialPropertyCode)
	}

	var rows int64
	if err := provider.GetGenericPersister().Table("credentials").Where("user_id = ?", userID).Count(&rows).Error; err != nil {
		t.Fatalf("counting credential rows: %v", err)
	}
	if rows != 0 {
		t.Errorf("%d credential rows were persisted, want 0 - the plaintext apiKey reached the datastore", rows)
	}
}

// TestUpdateUserCredentialRejectsAReservedPropertyFromTheRequestBody covers the
// other client-controlled write path, and proves the rejection does not damage
// the credential that is already stored.
func TestUpdateUserCredentialRejectsAReservedPropertyFromTheRequestBody(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	saved, err := provider.SaveUserCredential("tok", &Credential{
		Name:   "original",
		Type:   "token",
		UserId: userID,
		Secret: core.Map{"grafanaAPIKey": "keep-me"},
	})
	if err != nil {
		t.Fatalf("SaveUserCredential: %v", err)
	}

	_, err = provider.UpdateUserCredential(nil, &Credential{
		ID:     saved.ID,
		UserId: userID,
		Name:   "renamed",
		Type:   "token",
		Secret: core.Map{credentialCiphertextKey: credentialEnvelopePrefix + "aa:bb", "apiKey": "real-token"},
	})
	if err == nil {
		t.Fatal("an update carrying the reserved property was accepted")
	}
	if code := meshkiterrors.GetCode(err); code != ErrReservedCredentialPropertyCode {
		t.Errorf("error code = %s, want %s", code, ErrReservedCredentialPropertyCode)
	}

	fetched, _, err := provider.GetCredentialByID("tok", saved.ID)
	if err != nil {
		t.Fatalf("GetCredentialByID: %v", err)
	}
	if got := fetched.Secret["grafanaAPIKey"]; got != "keep-me" {
		t.Errorf("stored secret after a rejected update = %#v, want the original", fetched.Secret)
	}
}

// TestEncryptCredentialSecretDoesNotMutateInput pins the contract
// server/machines/actions.go depends on: it hands the connection payload's
// credential secret to SaveUserCredential and keeps using that same map
// afterwards.
func TestEncryptCredentialSecretDoesNotMutateInput(t *testing.T) {
	input := map[string]interface{}{"credentialName": "c", "secret": map[string]interface{}{"clusterToken": "tok"}}
	before := throughJSON(t, input)

	if _, err := EncryptCredentialSecret(input); err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}
	if !reflect.DeepEqual(input, before) {
		t.Errorf("input map was mutated:\n got %#v\nwant %#v", input, before)
	}
}

// TestEncryptCredentialSecretUsesAFreshNonce guards the one way AES-GCM fails
// catastrophically: reusing a nonce under the same key.
func TestEncryptCredentialSecretUsesAFreshNonce(t *testing.T) {
	secret := map[string]interface{}{"grafanaAPIKey": "key"}

	first, err := EncryptCredentialSecret(secret)
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}
	second, err := EncryptCredentialSecret(secret)
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}

	if first[credentialCiphertextKey] == second[credentialCiphertextKey] {
		t.Fatal("encrypting the same secret twice produced the same envelope; the nonce is not fresh per call")
	}
}

// TestDecryptCredentialSecretRejectsTampering proves the stored envelope is
// authenticated, not merely obscured.
func TestDecryptCredentialSecretRejectsTampering(t *testing.T) {
	encrypted, err := EncryptCredentialSecret(map[string]interface{}{"grafanaAPIKey": "key"})
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}
	envelope := encrypted[credentialCiphertextKey].(string)

	tampered := map[string]interface{}{
		"flipped payload byte": envelope[:len(envelope)-2] + flipBase64Char(envelope[len(envelope)-2:len(envelope)-1]) + envelope[len(envelope)-1:],
		"truncated payload":    envelope[:len(envelope)-8],
		"missing key id":       credentialEnvelopePrefix + strings.SplitN(strings.TrimPrefix(envelope, credentialEnvelopePrefix), ":", 2)[1],
		"non-base64 payload":   credentialEnvelopePrefix + activeCredentialKeyring().keyID + ":not base64!!",
	}

	for name, bad := range tampered {
		t.Run(name, func(t *testing.T) {
			if _, err := DecryptCredentialSecret(map[string]interface{}{credentialCiphertextKey: bad}); err == nil {
				t.Error("tampered envelope decrypted without error")
			}
		})
	}
}

// TestDecryptCredentialSecretReportsKeyMismatch covers the upgrade hazard the
// key identifier exists for: a datastore written by a build carrying a different
// $TOKEN (a locally built server, or a rotated GLOBAL_TOKEN). The failure has to
// name that cause, because the remedy - run the build that wrote them, or
// re-enter the credentials - is not deducible from "decryption failed".
func TestDecryptCredentialSecretReportsKeyMismatch(t *testing.T) {
	encrypted, err := EncryptCredentialSecret(map[string]interface{}{"grafanaAPIKey": "key"})
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}
	envelope := encrypted[credentialCiphertextKey].(string)
	payload := strings.SplitN(strings.TrimPrefix(envelope, credentialEnvelopePrefix), ":", 2)[1]

	foreign := credentialEnvelopePrefix + "0123456789ab:" + payload
	_, err = DecryptCredentialSecret(map[string]interface{}{credentialCiphertextKey: foreign})
	if err == nil {
		t.Fatal("an envelope naming a foreign key decrypted without error")
	}
	if got := meshkiterrors.GetCode(err); got != ErrCredentialKeyMismatchCode {
		t.Errorf("error code = %s, want %s", got, ErrCredentialKeyMismatchCode)
	}
	if !strings.Contains(err.Error(), "0123456789ab") {
		t.Errorf("error %q does not name the key the ciphertext was written under", err)
	}
	// The error reaches the API caller through writeMeshkitError, and the
	// identifier this build derives is a function of the build-time token.
	// Only the stored key id - which the caller already holds, in the row -
	// belongs in it.
	if active := activeCredentialKeyID(); active == "" {
		t.Fatal("the active keyring has no key identifier; the assertion below would be vacuous")
	} else if strings.Contains(err.Error(), active) {
		t.Errorf("error %q leaks the key identifier this build derives (%s)", err, active)
	}
}

// TestDecryptCredentialSecretRejectsMarkerWithoutEnvelope makes sure the
// reserved property is never handed to a caller as if it were credential data.
func TestDecryptCredentialSecretRejectsMarkerWithoutEnvelope(t *testing.T) {
	for name, value := range map[string]interface{}{
		"a plain string": "not-an-envelope",
		"an object":      map[string]interface{}{"nope": true},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := DecryptCredentialSecret(map[string]interface{}{credentialCiphertextKey: value}); err == nil {
				t.Error("reserved marker holding a non-envelope decrypted without error")
			}
		})
	}
}

// TestActiveCredentialKeyringIsDerivedNotRaw pins the KDF: the build-time token
// is never used as key material directly, and the key identifier is not the key.
func TestActiveCredentialKeyringIsDerivedNotRaw(t *testing.T) {
	ring := deriveCredentialKeyring("dev_token")
	if ring.err != nil {
		t.Fatalf("deriveCredentialKeyring: %v", ring.err)
	}
	if ring.keyID == "" {
		t.Fatal("derived keyring has no key identifier")
	}
	if len(ring.keyID) != credentialKeyIDLen*2 {
		t.Errorf("key identifier %q is %d chars, want %d", ring.keyID, len(ring.keyID), credentialKeyIDLen*2)
	}
	if strings.Contains(ring.keyID, "dev_token") {
		t.Error("key identifier leaks the build token")
	}
	if other := deriveCredentialKeyring("a-different-token"); other.keyID == ring.keyID {
		t.Error("two different build tokens derive the same key identifier")
	}
	if same := deriveCredentialKeyring("dev_token"); same.keyID != ring.keyID {
		t.Error("the same build token derived two different key identifiers")
	}
}

// TestSaveUserCredentialPersistsCiphertext is the end-to-end guarantee against a
// real datastore: the row on disk is ciphertext, and every read path hands back
// plaintext.
func TestSaveUserCredentialPersistsCiphertext(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	for _, tt := range credentialShapeFixtures() {
		t.Run(tt.name, func(t *testing.T) {
			saved, err := provider.SaveUserCredential("tok", &Credential{
				Name:   tt.name,
				Type:   "token",
				UserId: userID,
				Secret: tt.secret,
			})
			if err != nil {
				t.Fatalf("SaveUserCredential: %v", err)
			}

			// The caller keeps its plaintext; only the row is ciphertext.
			if !reflect.DeepEqual(map[string]interface{}(saved.Secret), tt.secret) {
				t.Errorf("SaveUserCredential returned a changed secret:\n got %#v\nwant %#v", saved.Secret, tt.secret)
			}

			var rawSecret string
			if err := provider.GetGenericPersister().Raw("SELECT secret FROM credentials WHERE id = ?", saved.ID).Scan(&rawSecret).Error; err != nil {
				t.Fatalf("reading the stored secret column: %v", err)
			}
			if !strings.Contains(rawSecret, credentialCiphertextKey) {
				t.Fatalf("stored secret column is not an envelope: %s", rawSecret)
			}

			fetched, _, err := provider.GetCredentialByID("tok", saved.ID)
			if err != nil {
				t.Fatalf("GetCredentialByID: %v", err)
			}
			if !reflect.DeepEqual(map[string]interface{}(fetched.Secret), tt.secret) {
				t.Errorf("GetCredentialByID returned:\n got %#v\nwant %#v", fetched.Secret, tt.secret)
			}
		})
	}

	page, err := provider.GetUserCredentials(nil, userID.String(), 0, 25, "", "created_at desc")
	if err != nil {
		t.Fatalf("GetUserCredentials: %v", err)
	}
	if len(page.Credentials) != len(credentialShapeFixtures()) {
		t.Fatalf("listed %d credentials, want %d", len(page.Credentials), len(credentialShapeFixtures()))
	}
	byName := map[string]map[string]interface{}{}
	for _, cred := range page.Credentials {
		byName[cred.Name] = cred.Secret
	}
	for _, tt := range credentialShapeFixtures() {
		if !reflect.DeepEqual(byName[tt.name], tt.secret) {
			t.Errorf("GetUserCredentials returned %s as:\n got %#v\nwant %#v", tt.name, byName[tt.name], tt.secret)
		}
	}
}

// TestPlaintextCredentialRowWrittenBeforeEncryptionStillReads is the upgrade
// case against a real datastore: a row that predates this change is plaintext,
// and every read path must return it untouched with no migration.
func TestPlaintextCredentialRowWrittenBeforeEncryptionStillReads(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	for _, tt := range credentialShapeFixtures() {
		// Insert straight through GORM, bypassing the provider, so the row is
		// exactly what a pre-encryption Meshery Server would have written.
		legacyRow := Credential{
			ID:     uuid.Must(uuid.NewV4()),
			Name:   tt.name,
			Type:   "token",
			UserId: userID,
			Secret: core.Map(tt.secret),
		}
		if err := provider.GetGenericPersister().Table("credentials").Create(&legacyRow).Error; err != nil {
			t.Fatalf("seeding pre-encryption row %s: %v", tt.name, err)
		}

		fetched, _, err := provider.GetCredentialByID("tok", legacyRow.ID)
		if err != nil {
			t.Fatalf("GetCredentialByID on a pre-encryption row: %v", err)
		}
		if !reflect.DeepEqual(map[string]interface{}(fetched.Secret), tt.secret) {
			t.Errorf("pre-encryption row %s read back as:\n got %#v\nwant %#v", tt.name, fetched.Secret, tt.secret)
		}
	}

	page, err := provider.GetUserCredentials(nil, userID.String(), 0, 25, "", "created_at desc")
	if err != nil {
		t.Fatalf("GetUserCredentials over pre-encryption rows: %v", err)
	}
	if len(page.Credentials) != len(credentialShapeFixtures()) {
		t.Fatalf("listed %d pre-encryption credentials, want %d", len(page.Credentials), len(credentialShapeFixtures()))
	}
	for _, cred := range page.Credentials {
		if cred.Secret == nil {
			t.Errorf("pre-encryption credential %s was listed without its secret", cred.Name)
		}
	}
}

// TestUpdateUserCredentialRewritesPlaintextRowAsCiphertext pins the migration
// story: there is no bulk migration, rows convert as they are rewritten.
func TestUpdateUserCredentialRewritesPlaintextRowAsCiphertext(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	legacyRow := Credential{
		ID:     uuid.Must(uuid.NewV4()),
		Name:   "legacy",
		Type:   "token",
		UserId: userID,
		Secret: core.Map{"grafanaAPIKey": "old-key"},
	}
	if err := provider.GetGenericPersister().Table("credentials").Create(&legacyRow).Error; err != nil {
		t.Fatalf("seeding pre-encryption row: %v", err)
	}

	updated, err := provider.UpdateUserCredential(nil, &Credential{
		ID:     legacyRow.ID,
		UserId: userID,
		Name:   "legacy",
		Type:   "token",
		Secret: core.Map{"grafanaAPIKey": "new-key"},
	})
	if err != nil {
		t.Fatalf("UpdateUserCredential: %v", err)
	}
	if got := updated.Secret["grafanaAPIKey"]; got != "new-key" {
		t.Errorf("UpdateUserCredential returned grafanaAPIKey = %v, want new-key", got)
	}

	var rawSecret string
	if err := provider.GetGenericPersister().Raw("SELECT secret FROM credentials WHERE id = ?", legacyRow.ID).Scan(&rawSecret).Error; err != nil {
		t.Fatalf("reading the stored secret column: %v", err)
	}
	if strings.Contains(rawSecret, "new-key") {
		t.Errorf("rewritten row still holds the secret in the clear: %s", rawSecret)
	}
	if !strings.Contains(rawSecret, credentialCiphertextKey) {
		t.Errorf("rewritten row is not an envelope: %s", rawSecret)
	}
}

// TestGetUserCredentialsDropsUndecryptableSecretsWithoutFailingThePage is the
// degradation contract. One row written under a different build's token must not
// make the whole credentials page unreadable; it is listed without its secret
// and the reason is logged.
func TestGetUserCredentialsDropsUndecryptableSecretsWithoutFailingThePage(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	if _, err := provider.SaveUserCredential("tok", &Credential{
		Name:   "readable",
		Type:   "token",
		UserId: userID,
		Secret: core.Map{"grafanaAPIKey": "key"},
	}); err != nil {
		t.Fatalf("SaveUserCredential: %v", err)
	}

	// Re-label a well-formed envelope with a key identifier this build cannot
	// derive: exactly what a row written by a differently built server looks
	// like.
	encrypted, err := EncryptCredentialSecret(map[string]interface{}{"grafanaAPIKey": "unreadable"})
	if err != nil {
		t.Fatalf("EncryptCredentialSecret: %v", err)
	}
	payload := strings.SplitN(strings.TrimPrefix(encrypted[credentialCiphertextKey].(string), credentialEnvelopePrefix), ":", 2)[1]
	foreignRow := Credential{
		ID:     uuid.Must(uuid.NewV4()),
		Name:   "written-by-another-build",
		Type:   "token",
		UserId: userID,
		Secret: core.Map{credentialCiphertextKey: credentialEnvelopePrefix + "0123456789ab:" + payload},
	}
	if err := provider.GetGenericPersister().Table("credentials").Create(&foreignRow).Error; err != nil {
		t.Fatalf("seeding foreign-key row: %v", err)
	}

	page, err := provider.GetUserCredentials(nil, userID.String(), 0, 25, "", "name asc")
	if err != nil {
		t.Fatalf("GetUserCredentials must not fail the whole page for one undecryptable row: %v", err)
	}
	if len(page.Credentials) != 2 {
		t.Fatalf("listed %d credentials, want 2", len(page.Credentials))
	}
	for _, cred := range page.Credentials {
		switch cred.Name {
		case "readable":
			if got := cred.Secret["grafanaAPIKey"]; got != "key" {
				t.Errorf("readable credential returned grafanaAPIKey = %v, want key", got)
			}
		case "written-by-another-build":
			if cred.Secret != nil {
				t.Errorf("undecryptable credential was listed with a secret: %#v", cred.Secret)
			}
		}
	}
}

// throughJSON puts a value through the same encode/decode the datastore does, so
// tests exercise what actually comes back off disk rather than the in-memory map
// they built.
func throughJSON(t *testing.T, in map[string]interface{}) map[string]interface{} {
	t.Helper()
	encoded, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(encoded, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

// flipBase64Char returns a different base64 character, so a tampering test
// changes the payload rather than accidentally reproducing it.
func flipBase64Char(s string) string {
	if s == "A" {
		return "B"
	}
	return "A"
}

// TestUpdateUserCredentialWithoutSecretKeepsTheStoredOne is a regression test
// for silent credential loss. UpdateUserCredential is a partial update, and the
// HTTP handler initialises Secret to an empty map before unmarshalling, so a
// request body carrying only a rename reaches the provider with Secret == {}.
// GORM does not treat an empty non-nil map as a zero value, so that used to be
// written over the credential's real secret - the credential silently stopped
// authenticating and nothing reported why.
func TestUpdateUserCredentialWithoutSecretKeepsTheStoredOne(t *testing.T) {
	provider := newTestProviderWithCredentialDB(t)
	userID := uuid.Must(uuid.NewV4())

	saved, err := provider.SaveUserCredential("tok", &Credential{
		Name:   "original",
		Type:   "token",
		UserId: userID,
		Secret: core.Map{"grafanaAPIKey": "keep-me"},
	})
	if err != nil {
		t.Fatalf("SaveUserCredential: %v", err)
	}

	updated, err := provider.UpdateUserCredential(nil, &Credential{
		ID:     saved.ID,
		UserId: userID,
		Name:   "renamed",
		Type:   "token",
		Secret: core.Map{}, // exactly what the handler produces for a body with no `secret`
	})
	if err != nil {
		t.Fatalf("UpdateUserCredential: %v", err)
	}
	if updated.Name != "renamed" {
		t.Errorf("name = %q, want renamed; the rename itself must still apply", updated.Name)
	}
	if got := updated.Secret["grafanaAPIKey"]; got != "keep-me" {
		t.Errorf("UpdateUserCredential returned grafanaAPIKey = %v, want keep-me", got)
	}

	fetched, _, err := provider.GetCredentialByID("tok", saved.ID)
	if err != nil {
		t.Fatalf("GetCredentialByID: %v", err)
	}
	if got := fetched.Secret["grafanaAPIKey"]; got != "keep-me" {
		t.Errorf("stored secret after a secretless update = %#v, want the original", fetched.Secret)
	}
}

// TestDecryptCredentialEnvelopeReadsAHandRolledV1Envelope pins the on-disk
// format across the constructor change underneath it. `meshery.enc.v1` is a
// persisted contract: envelopes were first written by sealing with
// cipher.NewGCM and prepending a hand-generated 96-bit nonce, and are now
// written by cipher.NewGCMWithRandomNonce, which is the constructor
// `GODEBUG=fips140=only` permits. The two layouts are byte-identical, and this
// builds the old one explicitly so that a future change that breaks the
// equivalence fails here rather than in an upgraded user's datastore.
func TestDecryptCredentialEnvelopeReadsAHandRolledV1Envelope(t *testing.T) {
	ring := activeCredentialKeyring()
	if ring.err != nil {
		t.Fatalf("activeCredentialKeyring: %v", ring.err)
	}

	ikm := sha256.Sum256([]byte(GlobalTokenForAnonymousResults))
	key, err := hkdf.Key(sha256.New, ikm[:], []byte(credentialKDFSalt), credentialKeyInfo, credentialKeyLen)
	if err != nil {
		t.Fatalf("hkdf.Key: %v", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		t.Fatalf("aes.NewCipher: %v", err)
	}
	legacyAEAD, err := cipher.NewGCM(block)
	if err != nil {
		t.Skipf("cipher.NewGCM is unavailable on this runtime (%v); the pre-change writer could not have run here either", err)
	}

	plaintext := []byte(`{"wrapped":false,"payload":{"grafanaAPIKey":"key"}}`)
	nonce := make([]byte, legacyAEAD.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		t.Fatalf("rand.Read: %v", err)
	}
	sealed := legacyAEAD.Seal(nonce, nonce, plaintext, nil)
	envelope := credentialEnvelopePrefix + ring.keyID + ":" + base64.StdEncoding.EncodeToString(sealed)

	decrypted, err := DecryptCredentialSecret(map[string]interface{}{credentialCiphertextKey: envelope})
	if err != nil {
		t.Fatalf("an envelope written by the pre-change sealer no longer decrypts: %v", err)
	}
	if got := decrypted["grafanaAPIKey"]; got != "key" {
		t.Errorf("grafanaAPIKey = %v, want key", got)
	}
}
