package models

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// EventsPersister assists with persisting events in local SQLite DB
type EventsPersister struct {
	DB *database.Handler
	// Log receives the warning raised when an event's metadata is reduced
	// by GuardEventMetadataSize. Optional.
	Log logger.Handler
}

type EventsResponse struct {
	Events               []*events.Event         `json:"events"`
	Page                 int                     `json:"page"`
	PageSize             int                     `json:"pageSize"`
	CountBySeverityLevel []*CountBySeverityLevel `json:"countBySeverityLevel"`
	TotalCount           int64                   `json:"totalCount"`
	ReadCount            int64                   `json:"readCount"`
}

type CountBySeverityLevel struct {
	Severity string `json:"severity"`
	Count    int    `json:"count"`
}

func (e *EventsPersister) getCountBySeverity(userID core.Uuid, eventStatus events.EventStatus) ([]*CountBySeverityLevel, error) {
	if eventStatus == "" {
		eventStatus = events.Unread
	}
	// Get the system ID from the config for the current instance. This is used to filter events that are not associated with the user but are associated with the system
	systemID := viper.GetString("INSTANCE_ID")
	sysID := uuid.FromStringOrNil(systemID)

	eventsBySeverity := []*CountBySeverityLevel{}
	err := e.DB.Model(&events.Event{}).
		Select("severity, count(severity) as count").
		Where("status = ? AND (owner = ? OR owner = ?)", eventStatus, userID, sysID).
		Group("severity").
		Find(&eventsBySeverity).Error

	if err != nil {
		return nil, err
	}

	return eventsBySeverity, nil
}

// PersistEvent records event in the local database, honouring the same
// createEvent contract the remote provider implements
// (docs/event-correlation-contract.md in meshery/schemas), so the local and
// remote providers behave alike:
//
//   - An event without a correlationKey, or one that identifies no user (no
//     owner), is recorded as a new row, with its correlation fields as sent.
//   - An event with a correlationKey and no resolvedAt reports an occurrence of
//     a condition. When the owner already has an open event with that key (one
//     without resolvedAt), no row is recorded: the open event's
//     occurrenceCount grows by the event's occurrenceCount (1 when unset), its
//     lastSeenAt becomes the later of its own and the event's (the time of
//     receipt when unset), and every other field keeps its first value.
//     Otherwise the event is recorded and opens the condition, with
//     occurrenceCount 1, firstSeenAt the event's firstSeenAt, else its
//     lastSeenAt, else now, and lastSeenAt the event's lastSeenAt, else
//     firstSeenAt.
//   - An event with a correlationKey and a resolvedAt is a resolution. It
//     closes the owner's open event with that key, recording nothing new: the
//     open event takes the resolvedAt, its lastSeenAt becomes the later of its
//     own and the resolution's (its lastSeenAt, else resolvedAt), and its
//     severity and description become the resolution's, which say what
//     cleared. A resolution is a close, never an occurrence: it does not
//     change occurrenceCount. With none open it is recorded as an already
//     resolved event with occurrenceCount 1 (when unset) and firstSeenAt and
//     lastSeenAt from the event, else its resolvedAt.
//
// A create is a replay when it carries an idempotencyKey its caller already
// used, or the ID of an event already stored, whatever that event's key or
// state: it records and changes nothing and is answered with the stored event.
// A retry of an event that has since been resolved is therefore never counted
// on, and never resolves, a newer condition opened with the same key.
//
// An owner has at most one open event per key. The partial unique index
// AutoMigrateSystemTables creates enforces it, and a create that loses the race
// to open a condition is retried once as an occurrence of the winner.
//
// Every locally persisted event - local provider, system events and the remote
// provider's fallback - is written here, so this is where
// PrepareEventForPersistence applies.
//
// It returns the event as stored: the recorded event when one was created, or
// the condition's open event, as updated, when an occurrence was counted on it
// or it was resolved. The event is the server's own: one that names no source
// is recorded with source EventSourceServer.
func (e *EventsPersister) PersistEvent(event events.Event, token string) (PersistedEvent, error) {
	return e.persist(withServerSource(event))
}

// PersistClientEvent records an event a client sent as PersistEvent does, but
// with its source as sent: one that names no source is recorded without one.
func (e *EventsPersister) PersistClientEvent(event events.Event, token string) (PersistedEvent, error) {
	return e.persist(event)
}

// CorrelatesConditions is always true: the local database applies the
// createEvent condition correlation contract.
func (e *EventsPersister) CorrelatesConditions() bool {
	return true
}

func (e *EventsPersister) persist(event events.Event) (PersistedEvent, error) {
	event = PrepareEventForPersistence(event, e.Log)
	// Replay runs first: a create carrying an idempotencyKey its caller
	// already used records nothing and is answered with the bound event. A
	// create without a key issues no query here; a retry by event ID is
	// caught by the insert's primary key instead.
	if bound, ok, err := boundEvent(e.DB.DB, event); err != nil {
		return PersistedEvent{}, ErrPersistEvent(err)
	} else if ok {
		return PersistedEvent{Event: bound, Disposition: EventReplayed}, nil
	}
	if !correlatesEvent(event) {
		// A new event is inserted, never upserted: the producer minted its ID
		// (events.NewEvent does, and BeforeCreate keeps it), so a plain INSERT
		// is the one statement it needs. GORM's Create requires a pointer.
		err := e.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&event).Error; err != nil {
				return err
			}
			return bindIdempotencyKey(tx, event, event.ID)
		})
		if err != nil {
			if isUniqueViolation(err) {
				if replayed, ok, lookupErr := e.replayedEvent(event); lookupErr != nil {
					return PersistedEvent{}, ErrPersistEvent(lookupErr)
				} else if ok {
					return replayed, nil
				}
			}
			return PersistedEvent{}, ErrPersistEvent(err)
		}
		return PersistedEvent{Event: &event, Disposition: EventCreated}, nil
	}

	persisted, err := e.persistCorrelatedEvent(event)
	if err != nil && isUniqueViolation(err) {
		// Either this event was already recorded (a retry of the same event,
		// whose ID is taken), or a concurrent create opened the condition
		// between this create's lookup and its insert. The first is a replay.
		// For the second, re-run the create: the lookup now finds the open
		// event and the occurrence is counted on it.
		if replayed, ok, lookupErr := e.replayedEvent(event); lookupErr != nil {
			return PersistedEvent{}, ErrPersistEvent(lookupErr)
		} else if ok {
			return replayed, nil
		}
		persisted, err = e.persistCorrelatedEvent(event)
	}
	if err != nil {
		return PersistedEvent{}, ErrPersistEvent(err)
	}
	return persisted, nil
}

// replayedEvent reports whether event was already recorded, and returns the
// stored event as the answer to the retry: a create carrying an idempotencyKey
// its caller already used is answered with the event the key is bound to, and
// an event value persisted twice (a producer retrying after a timeout; an
// event built with events.NewEvent carries its own ID) with the stored row of
// that ID. Either way the retry records nothing; a retry by ID binds the
// idempotencyKey it carries, as replayStoredEvent does.
func (e *EventsPersister) replayedEvent(event events.Event) (PersistedEvent, bool, error) {
	if bound, ok, err := boundEvent(e.DB.DB, event); err != nil || ok {
		if err != nil {
			return PersistedEvent{}, false, err
		}
		return PersistedEvent{Event: bound, Disposition: EventReplayed}, true, nil
	}
	replayed, ok, err := replayStoredEvent(e.DB.DB, event)
	if err != nil && isUniqueViolation(err) {
		// A concurrent create bound the key first. The key is checked before
		// the ID, so the retry is answered with the event the key is bound to.
		if bound, ok, boundErr := boundEvent(e.DB.DB, event); boundErr != nil {
			return PersistedEvent{}, false, boundErr
		} else if ok {
			return PersistedEvent{Event: bound, Disposition: EventReplayed}, true, nil
		}
	}
	return replayed, ok, err
}

// replayStoredEvent answers a create carrying the ID of a stored event as a
// replay of that event (see storedEventWithID), and binds the idempotencyKey
// the create carries, when it carries one in a caller's scope, to the stored
// event, as a create that records, counts or resolves an event does. A retry
// that pairs a fresh key with an event already stored is therefore answered
// the same way when it is retried again with only that key. A key bound
// concurrently fails the bind on the binding's primary key; the caller answers
// that as a replay of the bound event.
func replayStoredEvent(db *gorm.DB, event events.Event) (PersistedEvent, bool, error) {
	replayed, ok, err := storedEventWithID(db, event.ID)
	if err != nil || !ok {
		return PersistedEvent{}, false, err
	}
	if err := bindIdempotencyKey(db, event, replayed.Event.ID); err != nil {
		return PersistedEvent{}, false, err
	}
	return replayed, true, nil
}

// storedEventWithID returns the stored event with id, answered as a replay,
// and false when id is nil or no event has it. An event ID names one recorded
// event whatever its correlation key or state, so a create carrying the ID of
// a stored event is a retry of that event: it is answered with the event and
// is never applied to another row - not counted as an occurrence of, nor
// resolving, a newer condition opened with the same key.
func storedEventWithID(db *gorm.DB, id uuid.UUID) (PersistedEvent, bool, error) {
	if id.IsNil() {
		return PersistedEvent{}, false, nil
	}
	var stored events.Event
	found := db.Where("id = ?", id).Limit(1).Find(&stored)
	if found.Error != nil {
		return PersistedEvent{}, false, found.Error
	}
	if found.RowsAffected == 0 {
		return PersistedEvent{}, false, nil
	}
	return PersistedEvent{Event: &stored, Disposition: EventReplayed}, true, nil
}

func (e *EventsPersister) PersistSystemEvent(event events.Event) (PersistedEvent, error) {
	return e.PersistEvent(event, "")
}

// PersistAndPublishSystemEvent persists a system event and broadcasts it to
// userID as stored: with the stored row's ID, and, for an occurrence counted on or the
// resolution of an open condition, as the condition's stored event. When
// persisting fails, the persistence error is returned and broadcast is skipped
// to prevent ghost notifications. A nil persister skips persistence and
// broadcasts the event as sent. A nil broadcaster skips the broadcast.
//
// It is the system-event counterpart of the handlers' persistAndPublish: every
// site that persists a system event and broadcasts it goes through it, never
// publishing its own copy after PersistSystemEvent. A nil broadcaster skips
// the broadcast.
func PersistAndPublishSystemEvent(persister SystemEventPersister, broadcaster *Broadcast, userID core.Uuid, event *events.Event) error {
	var err error
	if persister != nil {
		var persisted PersistedEvent
		persisted, err = persister.PersistSystemEvent(*event)
		if err != nil {
			return err
		}
		if persisted.Event == nil {
			// Nothing new was recorded and the stored event is unknown.
			return nil
		}
		event = persisted.Event
	}
	if broadcaster != nil {
		broadcaster.Publish(userID, event)
	}
	return nil
}

// correlatedEventsMu serializes correlated creates and resolutions within this
// process. The local provider's persister and the system event sink are separate
// EventsPersister values over one SQLite database, and SQLite fails one of two
// transactions that both read and then write ("database is locked") instead of
// waiting. The partial unique index still guards the invariant against any
// other writer.
var correlatedEventsMu sync.Mutex

// ResolveConnectionCondition closes the open rows for a server-generated
// connection resolution. Each row keeps its owner, ID and occurrence history.
// The remote provider uses its owner-scoped event API instead.
func (l *DefaultLocalProvider) ResolveConnectionCondition(event events.Event) ([]events.Event, error) {
	if event.CorrelationKey == nil || *event.CorrelationKey == "" || event.ResolvedAt == nil {
		return nil, nil
	}

	correlatedEventsMu.Lock()
	defer correlatedEventsMu.Unlock()

	var resolved []events.Event
	now := time.Now().UTC()
	err := l.EventsPersister.DB.Transaction(func(tx *gorm.DB) error {
		var open []events.Event
		if err := tx.Where("correlation_key = ? AND resolved_at IS NULL AND owner IS NOT NULL AND owner <> ?",
			*event.CorrelationKey, uuid.Nil).Find(&open).Error; err != nil {
			return err
		}
		for i := range open {
			var persisted PersistedEvent
			updates := resolutionUpdates(&open[i], event, now)
			if err := updateOpenEvent(tx, &open[i], updates, EventResolved, &persisted); err != nil {
				return err
			}
			resolved = append(resolved, *persisted.Event)
		}
		return nil
	})
	if err != nil {
		return nil, ErrPersistEvent(err)
	}
	return resolved, nil
}

// persistCorrelatedEvent applies the createEvent correlation steps to an event
// for which correlatesEvent holds, in one transaction, and returns the event as
// stored.
func (e *EventsPersister) persistCorrelatedEvent(event events.Event) (PersistedEvent, error) {
	correlatedEventsMu.Lock()
	defer correlatedEventsMu.Unlock()

	now := time.Now().UTC()
	var persisted PersistedEvent
	err := e.DB.Transaction(func(tx *gorm.DB) error {
		// An event already recorded, sent again, is a retry of that event,
		// not another occurrence or the resolution of the condition. The ID is
		// looked up among every stored event, not only the condition's open
		// one: once the event is resolved and a newer event opens the same
		// key, its retry must still be answered with it rather than counted
		// on, or resolve, the newer condition. The key the retry carries is
		// bound to the event; a key bound concurrently fails the transaction
		// on the binding's primary key, and persist answers it as a replay of
		// the bound event.
		if replayed, ok, err := replayStoredEvent(tx, event); err != nil {
			return err
		} else if ok {
			persisted = replayed
			return nil
		}

		var open events.Event
		found := tx.
			Where("owner = ? AND correlation_key = ? AND resolved_at IS NULL", *event.Owner, *event.CorrelationKey).
			Limit(1).
			Find(&open)
		if found.Error != nil {
			return found.Error
		}

		if found.RowsAffected == 0 {
			if event.ResolvedAt == nil {
				applyConditionOpeningDefaults(&event, now)
			} else {
				applyResolvedConditionDefaults(&event)
			}
			if err := tx.Create(&event).Error; err != nil {
				return err
			}
			persisted = PersistedEvent{Event: &event, Disposition: EventCreated}
			return bindIdempotencyKey(tx, event, event.ID)
		}

		// The update runs on the stored row, whose status is valid, so the
		// event's BeforeUpdate status check passes.
		if event.ResolvedAt != nil {
			updates := resolutionUpdates(&open, event, now)
			if err := updateOpenEvent(tx, &open, updates, EventResolved, &persisted); err != nil {
				return err
			}
			return bindIdempotencyKey(tx, event, open.ID)
		}

		increment := occurrenceCountOf(event)
		lastSeen := now
		if event.LastSeenAt != nil {
			lastSeen = event.LastSeenAt.UTC()
		}
		if open.LastSeenAt != nil && open.LastSeenAt.After(lastSeen) {
			lastSeen = open.LastSeenAt.UTC()
		}
		if err := updateOpenEvent(tx, &open, map[string]interface{}{
			// An open row recorded by an older server, or as sent by a producer
			// that set no count, counts as one occurrence.
			"occurrence_count": gorm.Expr("COALESCE(occurrence_count, 1) + ?", increment),
			"last_seen_at":     lastSeen,
			"updated_at":       now,
		}, EventIncremented, &persisted); err != nil {
			return err
		}
		return bindIdempotencyKey(tx, event, open.ID)
	})
	if err != nil {
		return PersistedEvent{}, err
	}
	return persisted, nil
}

// resolutionUpdates records what cleared and retains the latest observation.
func resolutionUpdates(open *events.Event, event events.Event, now time.Time) map[string]interface{} {
	resolvedAt := event.ResolvedAt.UTC()
	lastSeen := resolvedAt
	if event.LastSeenAt != nil {
		lastSeen = event.LastSeenAt.UTC()
	}
	if open.LastSeenAt != nil && open.LastSeenAt.After(lastSeen) {
		lastSeen = open.LastSeenAt.UTC()
	}
	updates := map[string]interface{}{
		"resolved_at":  resolvedAt,
		"last_seen_at": lastSeen,
		"updated_at":   now,
	}
	if event.Severity != "" {
		updates["severity"] = event.Severity
	}
	if event.Description != "" {
		updates["description"] = event.Description
	}
	return updates
}

// updateOpenEvent applies updates to a condition's open event and reads it
// back, so the caller gets the event as stored rather than as it was before.
func updateOpenEvent(tx *gorm.DB, open *events.Event, updates map[string]interface{}, disposition EventDisposition, persisted *PersistedEvent) error {
	if err := tx.Model(open).Updates(updates).Error; err != nil {
		return err
	}
	var stored events.Event
	if err := tx.Where("id = ?", open.ID).Take(&stored).Error; err != nil {
		return err
	}
	*persisted = PersistedEvent{Event: &stored, Disposition: disposition}
	return nil
}

// occurrenceCountOf is the number of occurrences event reports: its
// occurrenceCount, and 1 when that is unset or below 1. Only ClientEventHandler
// validates a payload, so a server producer's zero or negative count must
// neither leave the stored count unchanged nor decrement it.
func occurrenceCountOf(event events.Event) int {
	if event.OccurrenceCount == nil || *event.OccurrenceCount < 1 {
		return 1
	}
	return *event.OccurrenceCount
}

// applyConditionOpeningDefaults fills the occurrence fields of an event that
// opens a condition, as createEvent does.
func applyConditionOpeningDefaults(event *events.Event, now time.Time) {
	count := occurrenceCountOf(*event)
	event.OccurrenceCount = &count
	if event.FirstSeenAt == nil {
		firstSeen := now
		if event.LastSeenAt != nil {
			firstSeen = *event.LastSeenAt
		}
		event.FirstSeenAt = &firstSeen
	}
	if event.LastSeenAt == nil {
		lastSeen := *event.FirstSeenAt
		event.LastSeenAt = &lastSeen
	}
}

// applyResolvedConditionDefaults fills the occurrence fields of a resolution
// recorded with no open condition, as createEvent does: one occurrence, seen
// at the payload's firstSeenAt and lastSeenAt, else at resolvedAt. A
// firstSeenAt that would fall after lastSeenAt is moved back to it.
func applyResolvedConditionDefaults(event *events.Event) {
	count := occurrenceCountOf(*event)
	event.OccurrenceCount = &count
	if event.FirstSeenAt == nil {
		firstSeen := *event.ResolvedAt
		event.FirstSeenAt = &firstSeen
	}
	if event.LastSeenAt == nil {
		lastSeen := *event.ResolvedAt
		event.LastSeenAt = &lastSeen
	}
	if event.FirstSeenAt.After(*event.LastSeenAt) {
		firstSeen := *event.LastSeenAt
		event.FirstSeenAt = &firstSeen
	}
}

// isUniqueViolation reports whether err is a unique constraint failure from
// either database engine Meshery Server runs on.
func isUniqueViolation(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") || strings.Contains(msg, "duplicate key value")
}
