package models

import (
	"encoding/json"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/database"
)

// KeyPersister is the persister for persisting keys in the database
type KeyPersister struct {
	DB *database.Handler
}

// GetUsersKeys returns all keys based on search, order, and updatedAfter without pagination
func (kp *KeyPersister) GetUsersKeys(search, order, updatedAfter string) ([]byte, error) {
	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "function", "category", "subcategory"})

	if order == "" {
		order = "updated_at desc"
	}

	count := int64(0)
	keys := []*Key{}

	query := kp.DB.Where("updated_at > ?", updatedAfter).Order(order)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(keys.function) like ? OR lower(keys.category) like ? OR lower(keys.subcategory) like ?)", like, like, like)
	}

	query.Table("keys").Count(&count).Find(&keys)

	keysPage := &KeysPage{
		Page:       1,          // Assuming default page to be 1
		PageSize:   int(count), // Set the total count as the page size for all keys
		TotalCount: int(count),
		Keys:       keys,
	}

	return marshalKeysPage(keysPage), nil
}

// SaveUsersKey saves a key to the database
func (kp *KeyPersister) SaveUsersKey(key *Key) (*Key, error) {
	if err := kp.DB.Create(key).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	return key, nil
}

// SaveUsersKeys saves a key to the database
func (kp *KeyPersister) SaveUsersKeys(keys []*Key) ([]*Key, error) {
	if err := kp.DB.Save(keys).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	return keys, nil
}

// GetUsersKey retrieves a key by its ID
func (kp *KeyPersister) GetUsersKey(id uuid.UUID) ([]byte, error) {
	var key Key
	err := kp.DB.First(&key, id).Error
	return marshalKey(&key), err
}

func marshalKey(k *Key) []byte {
	res, _ := json.Marshal(k)
	return res
}

func marshalKeysPage(ksPage *KeysPage) []byte {
	res, _ := json.Marshal(ksPage)
	return res
}
