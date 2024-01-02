package models

import (
	"strings"

	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/csv"
)

var rowIndex = 1

type KeysRegistrationHelper struct {
	log          logger.Handler
	keysChan     chan Key
	keyPersister *KeyPersister
}

func NewKeysRegistrationHelper(dbHandler *database.Handler, log logger.Handler) *KeysRegistrationHelper {
	return &KeysRegistrationHelper{
		log:      log,
		keysChan: make(chan Key, 1),
		keyPersister: &KeyPersister{
			DB: dbHandler,
		},
	}
}

func (kh *KeysRegistrationHelper) SeedKeys(filePath string) {
	ch := make(chan Key)

	csvReader, err := csv.NewCSVParser[Key](filePath, rowIndex, map[string]string{
		"Key ID": "id",
	}, func(columns []string, currentRow []string) bool {
		if len(columns) == len(currentRow) {
			shouldRegister := currentRow[len(columns)-1]
			return strings.ToLower(shouldRegister) == "true"
		}
		return false
	})

	if err != nil {
		kh.log.Error(err)
		return
	}

	go func() {
		err := csvReader.Parse(ch)
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
		case <-csvReader.Context.Done():
			return
		}
	}

}
