package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestKeyPersister(t *testing.T) (*KeyPersister, *database.Handler) {
	t.Helper()
	db, err := database.New(database.Options{Engine: database.SQLITE, Filename: ":memory:"})
	require.NoError(t, err)

	err = db.AutoMigrate(&Key{})
	require.NoError(t, err)

	return &KeyPersister{DB: &db}, &db
}

func TestKeyPersister_NilReceiverSafety(t *testing.T) {
	testID := core.Uuid(uuid.Must(uuid.NewV4()))
	testKey := &Key{ID: testID, Function: "test-func"}

	// 1. Nil persister pointer
	var nilKP *KeyPersister

	assert.NotPanics(t, func() {
		res, err := nilKP.GetUsersKeys("", "", "")
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := nilKP.SaveUsersKey(testKey)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := nilKP.SaveUsersKeys([]Key{*testKey})
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := nilKP.GetUsersKey(testID)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	// 2. Persister with nil DB handler
	kpNilDB := &KeyPersister{DB: nil}

	assert.NotPanics(t, func() {
		res, err := kpNilDB.GetUsersKeys("", "", "")
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := kpNilDB.SaveUsersKey(testKey)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := kpNilDB.SaveUsersKeys([]Key{*testKey})
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := kpNilDB.GetUsersKey(testID)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})
}

func TestKeyPersister_SaveUsersKeys_EmptySlice(t *testing.T) {
	kp, _ := setupTestKeyPersister(t)

	res, err := kp.SaveUsersKeys([]Key{})
	require.NoError(t, err)
	assert.Empty(t, res)
}

func TestKeyPersister_ZeroTimestampDeletedAtIsActive(t *testing.T) {
	kp, db := setupTestKeyPersister(t)

	keyID := core.Uuid(uuid.Must(uuid.NewV4()))
	k := &Key{
		ID:          keyID,
		Function:    "legacy_active_key",
		Category:    "core",
		Subcategory: "auth",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_, err := kp.SaveUsersKey(k)
	require.NoError(t, err)

	// Simulate legacy row where deleted_at holds a zero timestamp rather than SQL NULL
	zeroTime := core.NewTime(time.Time{}) // 0001-01-01 00:00:00
	err = db.Model(&Key{}).Where("id = ?", keyID).Update("deleted_at", zeroTime).Error
	require.NoError(t, err)

	// Key should still be retrievable by ID
	bytes, err := kp.GetUsersKey(keyID)
	require.NoError(t, err)
	assert.NotNil(t, bytes)

	// Key should also still appear in GetUsersKeys listing
	listBytes, err := kp.GetUsersKeys("legacy_active_key", "", "")
	require.NoError(t, err)
	var page KeysPage
	err = json.Unmarshal(listBytes, &page)
	require.NoError(t, err)
	assert.Equal(t, 1, page.TotalCount)
	require.Len(t, page.Keys, 1)
	assert.Equal(t, keyID, page.Keys[0].ID)
}

func TestKeyPersister_SaveUsersKey_NilKey(t *testing.T) {
	kp, _ := setupTestKeyPersister(t)

	res, err := kp.SaveUsersKey(nil)
	assert.Nil(t, res)
	assert.Equal(t, ErrResultData(), err)
}

func TestKeyPersister_SaveAndGetUsersKey(t *testing.T) {
	kp, _ := setupTestKeyPersister(t)

	keyID := core.Uuid(uuid.Must(uuid.NewV4()))
	ownerID := core.Uuid(uuid.Must(uuid.NewV4()))
	k := &Key{
		ID:          keyID,
		Owner:       ownerID,
		Function:    "manage_clusters",
		Category:    "infrastructure",
		Subcategory: "kubernetes",
		Description: "Key to manage cluster resources",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	saved, err := kp.SaveUsersKey(k)
	require.NoError(t, err)
	require.NotNil(t, saved)
	assert.Equal(t, keyID, saved.ID)

	// Retrieve by ID
	bytes, err := kp.GetUsersKey(keyID)
	require.NoError(t, err)
	require.NotNil(t, bytes)

	var retrieved Key
	err = json.Unmarshal(bytes, &retrieved)
	require.NoError(t, err)
	assert.Equal(t, keyID, retrieved.ID)
	assert.Equal(t, "manage_clusters", retrieved.Function)
	assert.Equal(t, "infrastructure", retrieved.Category)
	assert.Equal(t, "kubernetes", retrieved.Subcategory)
}

func TestKeyPersister_GetUsersKey_NotFound(t *testing.T) {
	kp, _ := setupTestKeyPersister(t)

	nonExistentID := core.Uuid(uuid.Must(uuid.NewV4()))
	bytes, err := kp.GetUsersKey(nonExistentID)

	assert.Nil(t, bytes)
	require.Error(t, err)
	assert.Equal(t, ErrResultNotFoundCode, meshkiterrors.GetCode(err))
}

func TestKeyPersister_GetUsersKey_SoftDeleted(t *testing.T) {
	kp, db := setupTestKeyPersister(t)

	keyID := core.Uuid(uuid.Must(uuid.NewV4()))
	k := &Key{
		ID:          keyID,
		Function:    "read_logs",
		Category:    "observability",
		Subcategory: "logging",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_, err := kp.SaveUsersKey(k)
	require.NoError(t, err)

	// Mark soft-deleted
	deletedAt := core.NewTime(time.Now().UTC())
	err = db.Model(&Key{}).Where("id = ?", keyID).Update("deleted_at", deletedAt).Error
	require.NoError(t, err)

	// Fetching soft-deleted key must return ErrResultNotFound
	bytes, err := kp.GetUsersKey(keyID)
	assert.Nil(t, bytes)
	require.Error(t, err)
	assert.Equal(t, ErrResultNotFoundCode, meshkiterrors.GetCode(err))
}

func TestKeyPersister_GetUsersKeys_FilterAndOrder(t *testing.T) {
	kp, db := setupTestKeyPersister(t)

	now := time.Now().UTC()
	key1 := Key{
		ID:          core.Uuid(uuid.Must(uuid.NewV4())),
		Function:    "deploy_workload",
		Category:    "workload",
		Subcategory: "deployment",
		CreatedAt:   now.Add(-2 * time.Hour),
		UpdatedAt:   now.Add(-2 * time.Hour),
	}
	key2 := Key{
		ID:          core.Uuid(uuid.Must(uuid.NewV4())),
		Function:    "view_metrics",
		Category:    "observability",
		Subcategory: "metrics",
		CreatedAt:   now.Add(-1 * time.Hour),
		UpdatedAt:   now.Add(-1 * time.Hour),
	}
	key3 := Key{
		ID:          core.Uuid(uuid.Must(uuid.NewV4())),
		Function:    "delete_cluster",
		Category:    "infrastructure",
		Subcategory: "kubernetes",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	savedKeys, err := kp.SaveUsersKeys([]Key{key1, key2, key3})
	require.NoError(t, err)
	assert.Len(t, savedKeys, 3)

	// 1. List all keys with explicit order by function asc
	bytes, err := kp.GetUsersKeys("", "function asc", "")
	require.NoError(t, err)
	var page KeysPage
	err = json.Unmarshal(bytes, &page)
	require.NoError(t, err)
	assert.Equal(t, 3, page.TotalCount)
	require.Len(t, page.Keys, 3)
	assert.Equal(t, "delete_cluster", page.Keys[0].Function)
	assert.Equal(t, "deploy_workload", page.Keys[1].Function)
	assert.Equal(t, "view_metrics", page.Keys[2].Function)

	// 2. Search by category
	bytes, err = kp.GetUsersKeys("observability", "", "")
	require.NoError(t, err)
	var searchPage KeysPage
	err = json.Unmarshal(bytes, &searchPage)
	require.NoError(t, err)
	assert.Equal(t, 1, searchPage.TotalCount)
	require.Len(t, searchPage.Keys, 1)
	assert.Equal(t, "view_metrics", searchPage.Keys[0].Function)

	// 3. Search by subcategory
	bytes, err = kp.GetUsersKeys("deployment", "", "")
	require.NoError(t, err)
	var subcatPage KeysPage
	err = json.Unmarshal(bytes, &subcatPage)
	require.NoError(t, err)
	assert.Equal(t, 1, subcatPage.TotalCount)
	require.Len(t, subcatPage.Keys, 1)
	assert.Equal(t, "deploy_workload", subcatPage.Keys[0].Function)

	// 4. Soft-delete key2 and verify it is excluded
	deletedAt := core.NewTime(now)
	err = db.Model(&Key{}).Where("id = ?", key2.ID).Update("deleted_at", deletedAt).Error
	require.NoError(t, err)

	bytes, err = kp.GetUsersKeys("", "", "")
	require.NoError(t, err)
	var afterDeletePage KeysPage
	err = json.Unmarshal(bytes, &afterDeletePage)
	require.NoError(t, err)
	assert.Equal(t, 2, afterDeletePage.TotalCount)
	assert.Len(t, afterDeletePage.Keys, 2)
	for _, k := range afterDeletePage.Keys {
		assert.NotEqual(t, key2.ID, k.ID)
	}

	// 5. Test updatedAfter filter
	futureTime := now.Add(24 * time.Hour).Format(time.RFC3339)
	bytes, err = kp.GetUsersKeys("", "", futureTime)
	require.NoError(t, err)
	var futurePage KeysPage
	err = json.Unmarshal(bytes, &futurePage)
	require.NoError(t, err)
	assert.Equal(t, 0, futurePage.TotalCount)
	assert.Empty(t, futurePage.Keys)
}

func TestKeyPersister_ClosedDB_ReturnsErrDBReadAndCreate(t *testing.T) {
	kp, db := setupTestKeyPersister(t)

	// Close underlying SQL database
	sqlDB, err := db.DB.DB()
	require.NoError(t, err)
	err = sqlDB.Close()
	require.NoError(t, err)

	// GetUsersKeys on closed DB should return ErrDBRead
	res, err := kp.GetUsersKeys("", "", "")
	assert.Nil(t, res)
	require.Error(t, err)
	assert.Equal(t, ErrDBReadCode, meshkiterrors.GetCode(err))

	// GetUsersKey on closed DB should return ErrDBRead
	testID := core.Uuid(uuid.Must(uuid.NewV4()))
	res, err = kp.GetUsersKey(testID)
	assert.Nil(t, res)
	require.Error(t, err)
	assert.Equal(t, ErrDBReadCode, meshkiterrors.GetCode(err))

	// SaveUsersKey on closed DB should return ErrDBCreate
	k := &Key{ID: testID, Function: "f"}
	saved, err := kp.SaveUsersKey(k)
	assert.Nil(t, saved)
	require.Error(t, err)
	assert.Equal(t, ErrDBCreateCode, meshkiterrors.GetCode(err))

	// SaveUsersKeys on closed DB should return ErrDBCreate
	savedList, err := kp.SaveUsersKeys([]Key{*k})
	assert.Nil(t, savedList)
	require.Error(t, err)
	assert.Equal(t, ErrDBCreateCode, meshkiterrors.GetCode(err))
}

func TestDefaultLocalProvider_Keys_NilPersisterSafety(t *testing.T) {
	provider := &DefaultLocalProvider{KeyPersister: nil}
	testID := uuid.Must(uuid.NewV4()).String()

	assert.NotPanics(t, func() {
		res, err := provider.GetUsersKeys("", "", "", "", "", "", "")
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := provider.GetUsersKey(nil, testID)
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})

	assert.NotPanics(t, func() {
		res, err := provider.SaveUsersKey("", []Key{})
		assert.Nil(t, res)
		assert.Equal(t, ErrDBConnection, err)
	})
}
