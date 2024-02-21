package models

import (
	"strings"

	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/csv"
)

var (
	rowIndex = 1
	// The column in the spreadsheet which tracks whether the key should be registerd with Local Provider or not.
	shouldRegister         = "Local Provider"
	shouldRegisterColIndex = -1
)

type KeysRegistrationHelper struct {
	log          logger.Handler
	keysChan     chan Key
	keyPersister *KeyPersister
}

func NewKeysRegistrationHelper(dbHandler *database.Handler, log logger.Handler) (*KeysRegistrationHelper, error) {
	krh := &KeysRegistrationHelper{
		log:      log,
		keysChan: make(chan Key, 1),
		keyPersister: &KeyPersister{
			DB: dbHandler,
		},
	}
	err := krh.keyPersister.DB.AutoMigrate(
		Key{},
	)
	return krh, err
}

// returns the spreadsheet column index that captures whether the key should be registered.
func (kh *KeysRegistrationHelper) GetIndexForRegisterCol(cols []string) int {
	if shouldRegisterColIndex != -1 {
		return shouldRegisterColIndex
	}

	for index, col := range cols {
		if col == shouldRegister {
			return index
		}
	}
	return shouldRegisterColIndex
}

func (kh *KeysRegistrationHelper) SeedKeys(filePath string) {
	ch := make(chan Key, 1)
	errorChan := make(chan error, 1)
	csvReader, err := csv.NewCSVParser[Key](filePath, rowIndex, map[string]string{
		"Key ID": "id",
	}, func(columns []string, currentRow []string) bool {
		index := kh.GetIndexForRegisterCol(columns)
		if index != -1 && index < len(currentRow) {
			shouldRegister := currentRow[index]
			return strings.ToLower(shouldRegister) == "true"
		}
		return false
	})

	if err != nil {
		kh.log.Error(err)
		return
	}

	go func() {
		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			kh.log.Error(err)
		}
	}()
	for {
		select {

		case data := <-ch:
			_, err := kh.keyPersister.SaveUsersKey(&data)
			if err != nil {
				kh.log.Error(err)
			}
		case err := <-errorChan:
			kh.log.Error(err)

		case <-csvReader.Context.Done():
			return
		}
	}

}
