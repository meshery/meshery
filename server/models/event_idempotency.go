package models

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/models/events"
	"gorm.io/gorm"
)

// EventIdempotencyBinding binds a caller's idempotencyKey to the event the
// create that carried it recorded, incremented or resolved (step 5 of the
// createEvent contract, docs/event-correlation-contract.md in meshery/schemas).
// A later create with the same key from the same caller records nothing and is
// answered with the bound event (step 1).
//
// The binding is its own table because the bound event is not always the row
// the create inserted: an occurrence counted on an open condition, or its
// resolution, binds the key to the condition's event, which keeps the key of
// the create that opened it. The primary key makes binding atomic: of two
// concurrent creates with one key, one binds and the other is a replay.
type EventIdempotencyBinding struct {
	Owner          uuid.UUID `gorm:"column:owner;type:uuid;primaryKey"`
	IdempotencyKey string    `gorm:"column:idempotency_key;primaryKey"`
	EventID        uuid.UUID `gorm:"column:event_id;type:uuid;not null"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

// TableName is the binding table's name.
func (EventIdempotencyBinding) TableName() string { return "event_idempotency_bindings" }

// idempotencyScope returns the caller and key a create's idempotencyKey is
// scoped to, and false when the create carries no key or identifies no user:
// such a create takes no part in replay.
func idempotencyScope(event events.Event) (uuid.UUID, string, bool) {
	if event.IdempotencyKey == nil || *event.IdempotencyKey == "" || event.Owner == nil || event.Owner.IsNil() {
		return uuid.Nil, "", false
	}
	return *event.Owner, *event.IdempotencyKey, true
}

// bindIdempotencyKey binds event's idempotencyKey, when it carries one in a
// caller's scope, to the stored event eventID.
func bindIdempotencyKey(tx *gorm.DB, event events.Event, eventID uuid.UUID) error {
	owner, key, ok := idempotencyScope(event)
	if !ok {
		return nil
	}
	return tx.Create(&EventIdempotencyBinding{Owner: owner, IdempotencyKey: key, EventID: eventID, CreatedAt: time.Now().UTC()}).Error
}

// boundEvent returns the event event's idempotencyKey is bound to, and false
// when it carries no key in a caller's scope or the key is unused. A binding
// whose event no longer exists (it was deleted) frees the key: it is removed
// and the create proceeds.
//
// The stale binding is removed by its event ID as well as its owner and key,
// so only the binding that was read as stale is removed. Two creates that both
// read the stale binding may both proceed, but once one of them has bound the
// key to its own event the other's delete matches nothing, its bind then
// fails on the binding's primary key, and it is answered as a replay of the
// winner, so one key never records two events.
func boundEvent(db *gorm.DB, event events.Event) (*events.Event, bool, error) {
	owner, key, ok := idempotencyScope(event)
	if !ok {
		return nil, false, nil
	}
	var binding EventIdempotencyBinding
	found := db.Where("owner = ? AND idempotency_key = ?", owner, key).Limit(1).Find(&binding)
	if found.Error != nil {
		return nil, false, found.Error
	}
	if found.RowsAffected == 0 {
		return nil, false, nil
	}
	var stored events.Event
	read := db.Where("id = ?", binding.EventID).Limit(1).Find(&stored)
	if read.Error != nil {
		return nil, false, read.Error
	}
	if read.RowsAffected == 0 {
		if err := db.Where("owner = ? AND idempotency_key = ? AND event_id = ?", owner, key, binding.EventID).Delete(&EventIdempotencyBinding{}).Error; err != nil {
			return nil, false, err
		}
		return nil, false, nil
	}
	return &stored, true, nil
}
