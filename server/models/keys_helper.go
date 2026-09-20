package models

import (
	"strings"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils/csv"
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

// GetIndexForRegisterCol returns the spreadsheet column index that captures whether the
// key should be registered, or -1 if the column is absent.
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

func (kh *KeysRegistrationHelper) SeedKeys(seedLog *SeedLog, filePath string) {
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
		seedLog.Errorf("Failed to read keys CSV %s: %v", filePath, err)
		return
	}

	seeded := 0
	failures := 0
	go func() {
		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			seedLog.Errorf("Failed to parse keys CSV %s: %v", filePath, err)
		}
	}()
	for {
		select {

		case data := <-ch:
			_, err := kh.keyPersister.SaveUsersKey(&data)
			if err != nil {
				failures++
				seedLog.Detailf("Failed to save key %s: %v", data.ID.String(), err)
				kh.log.Error(err)
			} else {
				seeded++
				seedLog.Detailf("Seeded key %s.", data.ID.String())
			}
		case err := <-errorChan:
			failures++
			seedLog.Errorf("Error while parsing keys CSV %s: %v", filePath, err)

		case <-csvReader.Context.Done():
			seedLog.Reportf("Seeded %d keys from %s (%d failures).", seeded, filePath, failures)
			return
		}
	}

}
