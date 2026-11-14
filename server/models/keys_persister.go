package models

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"gorm.io/gorm"
)

// KeyPersister is the persister for persisting keys in the database
type KeyPersister struct {
	DB *database.Handler
}

// GetUsersKeys returns all keys based on search, order, and updatedAfter without pagination
func (kp *KeyPersister) GetUsersKeys(search, order, updatedAfter string) ([]byte, error) {
	if kp == nil || kp.DB == nil {
		return nil, ErrDBConnection
	}

	order = SanitizeOrderInput(order, []string{"created_at", "updated_at", "function", "category", "subcategory"})

	if order == "" {
		order = defaultOrderUpdatedAtDesc
	}

	count := int64(0)
	keys := []Key{}

	query := kp.DB.Table("keys").Where("deleted_at IS NULL OR deleted_at <= '1970-01-01 00:00:00'")
	if updatedAfter != "" {
		query = query.Where("updated_at > ?", updatedAfter)
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("(lower(keys.function) like ? OR lower(keys.category) like ? OR lower(keys.subcategory) like ?)", like, like, like)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, ErrDBRead(err)
	}

	if err := query.Order(order).Find(&keys).Error; err != nil {
		return nil, ErrDBRead(err)
	}

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
	if kp == nil || kp.DB == nil {
		return nil, ErrDBConnection
	}

	if key == nil {
		return nil, ErrResultData()
	}

	if err := kp.DB.Save(key).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	return key, nil
}

// SaveUsersKeys saves multiple keys to the database
func (kp *KeyPersister) SaveUsersKeys(keys []Key) ([]Key, error) {
	if kp == nil || kp.DB == nil {
		return nil, ErrDBConnection
	}
	if len(keys) == 0 {
		return keys, nil
	}

	if err := kp.DB.Save(&keys).Error; err != nil {
		return nil, ErrDBCreate(err)
	}

	return keys, nil
}

// GetUsersKey retrieves a key by its ID
func (kp *KeyPersister) GetUsersKey(id core.Uuid) ([]byte, error) {
	if kp == nil || kp.DB == nil {
		return nil, ErrDBConnection
	}

	var key Key
	err := kp.DB.Where("id = ? AND (deleted_at IS NULL OR deleted_at <= '1970-01-01 00:00:00')", id).First(&key).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrResultNotFound(err)
		}
		return nil, ErrDBRead(err)
	}

	return marshalKey(&key), nil
}

func marshalKey(k *Key) []byte {
	res, _ := json.Marshal(k)
	return res
}

func marshalKeysPage(ksPage *KeysPage) []byte {
	res, _ := json.Marshal(ksPage)
	return res
}
