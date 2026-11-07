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

func (e *EventsPersister) persist(event events.Event) (PersistedEvent, error) {
	event = PrepareEventForPersistence(event, e.Log)
	if !correlatesEvent(event) {
		// GORM's Save requires a pointer (it reflects on the struct to read and
		// update the primary key and timestamps). Passing the value directly
		// produced "invalid value, should be pointer to struct or slice" and
		// dropped every system event silently.
		if err := e.DB.Save(&event).Error; err != nil {
			return PersistedEvent{}, ErrPersistEvent(err)
		}
		return PersistedEvent{Event: &event, Disposition: EventCreated}, nil
	}

	persisted, err := e.persistCorrelatedEvent(event)
	if err != nil && isUniqueViolation(err) {
		// A concurrent create opened the condition between this create's lookup
		// and its insert. Re-run it: the lookup now finds the open event and
		// the occurrence is counted on it.
		persisted, err = e.persistCorrelatedEvent(event)
	}
	if err != nil {
		return PersistedEvent{}, ErrPersistEvent(err)
	}
	return persisted, nil
}

func (e *EventsPersister) PersistSystemEvent(event events.Event) (PersistedEvent, error) {
	return e.PersistEvent(event, "")
}

// PersistAndPublishSystemEvent persists a system event and broadcasts it to
// userID as stored: with the stored row's ID, which an event built with a bare
// events.NewEvent() does not carry, and, for an occurrence counted on or the
// resolution of an open condition, as the condition's stored event. When
// persisting fails, or there is no persister, the event is still broadcast, as
// sent, and the persistence error is returned for the caller to report.
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
		if err == nil && persisted.Event != nil {
			event = persisted.Event
		}
	}
	if broadcaster != nil {
		broadcaster.Publish(userID, event)
	}
	return err
}

// correlatedEventsMu serializes correlated creates within this process. The
// local provider's persister and the system event sink are separate
// EventsPersister values over one SQLite database, and SQLite fails one of two
// transactions that both read and then write ("database is locked") instead of
// waiting. The partial unique index still guards the invariant against any
// other writer.
var correlatedEventsMu sync.Mutex

// persistCorrelatedEvent applies the createEvent correlation steps to an event
// for which correlatesEvent holds, in one transaction, and returns the event as
// stored.
func (e *EventsPersister) persistCorrelatedEvent(event events.Event) (PersistedEvent, error) {
	correlatedEventsMu.Lock()
	defer correlatedEventsMu.Unlock()

	now := time.Now().UTC()
	var persisted PersistedEvent
	err := e.DB.Transaction(func(tx *gorm.DB) error {
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
			return nil
		}

		// The update runs on the stored row, whose status is valid, so the
		// event's BeforeUpdate status check passes.
		if event.ResolvedAt != nil {
			resolvedAt := event.ResolvedAt.UTC()
			// The resolution is the condition's latest observation, and it
			// says what cleared: the closed event takes its time, severity and
			// description (for example success, "Connection restored").
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
			return updateOpenEvent(tx, &open, updates, EventResolved, &persisted)
		}

		increment := 1
		if event.OccurrenceCount != nil {
			increment = *event.OccurrenceCount
		}
		lastSeen := now
		if event.LastSeenAt != nil {
			lastSeen = event.LastSeenAt.UTC()
		}
		if open.LastSeenAt != nil && open.LastSeenAt.After(lastSeen) {
			lastSeen = open.LastSeenAt.UTC()
		}
		return updateOpenEvent(tx, &open, map[string]interface{}{
			// An open row recorded by an older server, or as sent by a producer
			// that set no count, counts as one occurrence.
			"occurrence_count": gorm.Expr("COALESCE(occurrence_count, 1) + ?", increment),
			"last_seen_at":     lastSeen,
			"updated_at":       now,
		}, EventIncremented, &persisted)
	})
	if err != nil {
		return PersistedEvent{}, err
	}
	return persisted, nil
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

// applyConditionOpeningDefaults fills the occurrence fields of an event that
// opens a condition, as createEvent does.
func applyConditionOpeningDefaults(event *events.Event, now time.Time) {
	if event.OccurrenceCount == nil {
		one := 1
		event.OccurrenceCount = &one
	}
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
	if event.OccurrenceCount == nil {
		one := 1
		event.OccurrenceCount = &one
	}
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
