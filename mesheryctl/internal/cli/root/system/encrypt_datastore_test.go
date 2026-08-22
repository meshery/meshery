package system

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/server/pkg/encryption"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestEncryptionService(t *testing.T) *encryption.Service {
	t.Helper()
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i + 1)
	}
	svc, err := encryption.New(key)
	if err != nil {
		t.Fatalf("encryption.New: %v", err)
	}
	return svc
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dbFile := filepath.Join(t.TempDir(), "test_mesherydb.sql")
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	if err := db.AutoMigrate(&rawCredential{}, &rawK8sContext{}); err != nil {
		t.Fatalf("AutoMigrate: %v", err)
	}

	return db
}

func TestMigrateCredentialsTable(t *testing.T) {
	db := setupTestDB(t)
	svc := newTestEncryptionService(t)

	// Insert test credentials
	originalSecret := map[string]interface{}{
		"apiKey": "super-secret-token-12345",
	}
	secretBytes, _ := json.Marshal(originalSecret)
	cred := rawCredential{
		ID:     "cred-1",
		Secret: string(secretBytes),
	}
	if err := db.Create(&cred).Error; err != nil {
		t.Fatalf("db.Create: %v", err)
	}

	// 1. Dry run encryption
	count, err := migrateCredentialsTable(db, svc, false, true)
	if err != nil {
		t.Fatalf("migrateCredentialsTable (dry-run encrypt): %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 credential to be affected, got %d", count)
	}
	var afterDryRun rawCredential
	if err := db.First(&afterDryRun, "id = ?", "cred-1").Error; err != nil {
		t.Fatalf("db.First: %v", err)
	}
	if afterDryRun.Secret != string(secretBytes) {
		t.Error("dry-run modified the database")
	}

	// 2. Perform encryption
	count, err = migrateCredentialsTable(db, svc, false, false)
	if err != nil {
		t.Fatalf("migrateCredentialsTable (encrypt): %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 credential encrypted, got %d", count)
	}

	var encryptedRow rawCredential
	if err := db.First(&encryptedRow, "id = ?", "cred-1").Error; err != nil {
		t.Fatalf("db.First: %v", err)
	}
	var encMap map[string]interface{}
	if err := json.Unmarshal([]byte(encryptedRow.Secret), &encMap); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if !svc.IsEncrypted(encMap) {
		t.Fatal("expected credential row to be encrypted")
	}

	// Re-encrypting should be a no-op (count = 0)
	count, err = migrateCredentialsTable(db, svc, false, false)
	if err != nil {
		t.Fatalf("migrateCredentialsTable (second encrypt): %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 credentials to be re-encrypted, got %d", count)
	}

	// 3. Perform decryption
	count, err = migrateCredentialsTable(db, svc, true, false)
	if err != nil {
		t.Fatalf("migrateCredentialsTable (decrypt): %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 credential decrypted, got %d", count)
	}

	var decryptedRow rawCredential
	if err := db.First(&decryptedRow, "id = ?", "cred-1").Error; err != nil {
		t.Fatalf("db.First: %v", err)
	}
	var decMap map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedRow.Secret), &decMap); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if svc.IsEncrypted(decMap) {
		t.Fatal("expected credential row to be decrypted")
	}
	if decMap["apiKey"] != "super-secret-token-12345" {
		t.Errorf("decrypted apiKey = %v, want super-secret-token-12345", decMap["apiKey"])
	}
}

func TestMigrateK8sContextsTable(t *testing.T) {
	db := setupTestDB(t)
	svc := newTestEncryptionService(t)

	authMap := map[string]interface{}{"token": "bearer-token-abc"}
	clusterMap := map[string]interface{}{"server": "https://10.0.0.1:6443"}
	authBytes, err := json.Marshal(authMap)
	if err != nil {
		t.Fatalf("json.Marshal auth: %v", err)
	}
	clusterBytes, err := json.Marshal(clusterMap)
	if err != nil {
		t.Fatalf("json.Marshal cluster: %v", err)
	}

	k8sCtx := rawK8sContext{
		ID:      "ctx-1",
		Auth:    string(authBytes),
		Cluster: string(clusterBytes),
	}
	if err := db.Create(&k8sCtx).Error; err != nil {
		t.Fatalf("db.Create: %v", err)
	}

	// 1. Perform encryption
	count, err := migrateK8sContextsTable(db, svc, false, false)
	if err != nil {
		t.Fatalf("migrateK8sContextsTable (encrypt): %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 context encrypted, got %d", count)
	}

	var encryptedRow rawK8sContext
	if err := db.First(&encryptedRow, "id = ?", "ctx-1").Error; err != nil {
		t.Fatalf("db.First: %v", err)
	}
	var encAuth, encCluster map[string]interface{}
	if err := json.Unmarshal([]byte(encryptedRow.Auth), &encAuth); err != nil {
		t.Fatalf("json.Unmarshal auth: %v", err)
	}
	if err := json.Unmarshal([]byte(encryptedRow.Cluster), &encCluster); err != nil {
		t.Fatalf("json.Unmarshal cluster: %v", err)
	}

	if !svc.IsEncrypted(encAuth) {
		t.Fatal("expected Auth to be encrypted")
	}
	if !svc.IsEncrypted(encCluster) {
		t.Fatal("expected Cluster to be encrypted")
	}

	// 2. Perform decryption
	count, err = migrateK8sContextsTable(db, svc, true, false)
	if err != nil {
		t.Fatalf("migrateK8sContextsTable (decrypt): %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 context decrypted, got %d", count)
	}

	var decryptedRow rawK8sContext
	if err := db.First(&decryptedRow, "id = ?", "ctx-1").Error; err != nil {
		t.Fatalf("db.First: %v", err)
	}
	var decAuth, decCluster map[string]interface{}
	if err := json.Unmarshal([]byte(decryptedRow.Auth), &decAuth); err != nil {
		t.Fatalf("json.Unmarshal auth: %v", err)
	}
	if err := json.Unmarshal([]byte(decryptedRow.Cluster), &decCluster); err != nil {
		t.Fatalf("json.Unmarshal cluster: %v", err)
	}

	if svc.IsEncrypted(decAuth) {
		t.Fatal("expected Auth to be decrypted")
	}
	if svc.IsEncrypted(decCluster) {
		t.Fatal("expected Cluster to be decrypted")
	}
	if decAuth["token"] != "bearer-token-abc" {
		t.Errorf("decrypted Auth token = %v, want bearer-token-abc", decAuth["token"])
	}
}

func TestEncryptDatastoreCmd_FlagDefaults(t *testing.T) {
	cmd := encryptDatastoreCmd
	if cmd.Use != "encrypt-datastore" {
		t.Errorf("unexpected command use: %s", cmd.Use)
	}

	// Verify flags exist
	if cmd.Flags().Lookup("key") == nil {
		t.Error("flag --key missing")
	}
	if cmd.Flags().Lookup("key-file") == nil {
		t.Error("flag --key-file missing")
	}
	if cmd.Flags().Lookup("db-path") == nil {
		t.Error("flag --db-path missing")
	}
	if cmd.Flags().Lookup("dry-run") == nil {
		t.Error("flag --dry-run missing")
	}
	if cmd.Flags().Lookup("decrypt") == nil {
		t.Error("flag --decrypt missing")
	}
}
