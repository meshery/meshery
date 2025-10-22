package models

import (
	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
)

// KeyPersister is the persister for persisting keys in the database
type KeyPersister struct {
	BasePersister
}

// NewKeyPersister creates a new key persister with optimized base functionality
func NewKeyPersister(db *database.Handler) *KeyPersister {
	return &KeyPersister{
		BasePersister: BasePersister{DB: db},
	}
}

// GetUsersKeys returns all keys with optimized query building
func (kp *KeyPersister) GetUsersKeys(search, order, updatedAfter string) ([]byte, error) {
	keys := []*Key{}

	qb := kp.NewQueryBuilder(&Key{}).
		WithSearch(search, []string{"keys.function", "keys.category", "keys.subcategory"}).
		WithOrder(order, []string{"created_at", "updated_at", "function", "category", "subcategory"}).
		WithUpdatedAfter(updatedAfter)

	count := qb.Count("keys")
	if err := qb.Find(&keys); err != nil {
		return nil, ErrDBRead(err)
	}

	keysPage := &KeysPage{
		Page:       1,
		PageSize:   int(count),
		TotalCount: int(count),
		Keys:       keys,
	}

	return MarshalJSON(keysPage), nil
}

// SaveUsersKey saves a key using base persister utilities
func (kp *KeyPersister) SaveUsersKey(key *Key) (*Key, error) {
	if err := kp.SaveEntity(key, &key.ID); err != nil {
		return nil, ErrDBCreate(err)
	}
	return key, nil
}

// SaveUsersKeys saves multiple keys in a batch operation for better performance
func (kp *KeyPersister) SaveUsersKeys(keys []*Key) ([]*Key, error) {
	// Generate UUIDs for keys that don't have them
	for _, key := range keys {
		if key.ID == uuid.Nil {
			id, err := kp.GenerateUUID()
			if err != nil {
				return nil, ErrGenerateUUID(err)
			}
			key.ID = id
		}
	}

	// Batch insert for better performance
	if err := kp.DB.CreateInBatches(keys, 100).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	return keys, nil
}

// GetUsersKey retrieves a key by its ID
func (kp *KeyPersister) GetUsersKey(id uuid.UUID) ([]byte, error) {
	var key Key
	if err := kp.FindByID(&key, id); err != nil {
		return nil, ErrDBRead(err)
	}
	return MarshalJSON(&key), nil
}
