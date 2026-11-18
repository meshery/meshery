package models

import (
	"runtime/debug"
	"strings"

	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils/csv"
)

var (
	rowIndex = 1
	// The column in the spreadsheet which tracks whether the key should be registerd with Local Provider or not.
	shouldRegister = "Local Provider"
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
	for index, col := range cols {
		if col == shouldRegister {
			return index
		}
	}
	return -1
}

func (kh *KeysRegistrationHelper) SeedKeys(filePath string) {
	ch := make(chan Key, 1)
	errorChan := make(chan error, 1)
	// The header row is identical for every row the parser visits, so a
	// missing register column is reported once, on the first row it blocks,
	// rather than once per row. Without this, a renamed or dropped header
	// selects no row at all and SeedKeys persists zero keys silently.
	registerColMissingReported := false
	csvReader, err := csv.NewCSVParser[Key](filePath, rowIndex, map[string]string{
		"Key ID": "id",
	}, func(columns []string, currentRow []string) bool {
		index := kh.GetIndexForRegisterCol(columns)
		if index != -1 && index < len(currentRow) {
			shouldRegister := currentRow[index]
			return strings.ToLower(shouldRegister) == "true"
		}
		if !registerColMissingReported {
			registerColMissingReported = true
			kh.log.Error(ErrKeysRegisterColumnMissing(shouldRegister))
		}
		return false
	})

	if err != nil {
		kh.log.Error(err)
		return
	}

	// This goroutine is outside the reach of RunSeedStage's recover, which
	// only covers the stage's own goroutine, so it recovers on its own. The
	// failure is reported through ErrSeedingStagePanic with the same stage
	// name the call sites register ("user keys"), making it indistinguishable
	// in the log from a panic RunSeedStage caught itself. Parse's deferred
	// cancel closes its context while unwinding, so the select loop below
	// still observes Done and returns instead of hanging.
	go func() {
		defer func() {
			if r := recover(); r != nil {
				kh.log.Error(ErrSeedingStagePanic("user keys", r, debug.Stack()))
			}
		}()
		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			kh.log.Error(err)
		}
	}()
	for {
		select {

		case data := <-ch:
			kh.saveSeededKey(data)
		case err := <-errorChan:
			kh.log.Error(err)

		case <-csvReader.Context.Done():
			// The parser calls its deferred cancel only as Parse returns,
			// after every send has completed, so once Done is observed no
			// further send can occur. Either channel may still hold a
			// buffered row or error alongside Done, so drain both
			// non-blockingly before returning instead of dropping
			// whatever the select did not pick.
			for {
				select {
				case data := <-ch:
					kh.saveSeededKey(data)
				case err := <-errorChan:
					kh.log.Error(err)
				default:
					return
				}
			}
		}
	}

}

// saveSeededKey persists one parsed key. A panic from the persistence call is
// contained here so it cannot unwind the seed loop: without a consumer the
// parser would block forever on its next send and never run the deferred
// cleanup that ends parsing, leaking the goroutine and the file it holds.
func (kh *KeysRegistrationHelper) saveSeededKey(data Key) {
	defer func() {
		if r := recover(); r != nil {
			kh.log.Error(ErrSeedingStagePanic("user keys", r, debug.Stack()))
		}
	}()
	if _, err := kh.keyPersister.SaveUsersKey(&data); err != nil {
		kh.log.Error(err)
	}
}
