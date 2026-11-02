package models

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// Event metadata keys that identify the design an event is about. Events
// carry the identity of a design, never its contents (#22764).
const (
	EventMetadataDesignID      = "designId"
	EventMetadataDesignName    = "designName"
	EventMetadataDesignVersion = "designVersion"
	EventMetadataDesignLink    = "designLink"
	EventMetadataHistoryTitle  = "historyTitle"
)

// MaxEventMetadataBytes is the largest JSON-serialized event metadata that is
// persisted as-is. Anything larger is reduced to the design identity by
// GuardEventMetadataSize.
const MaxEventMetadataBytes = 16 * 1024

// designViewLinkFormat is the meshmap view URL for a design.
//
// NOTE: design_id (snake_case) is the meshmap extension's URL query-param
// contract; flipping it requires a coordinated meshmap-side change.
const designViewLinkFormat = "%s/extension/meshmap?mode=operator&type=view&design_id=%s"

// DesignLink returns the URL that opens the design in the meshmap extension.
func DesignLink(serverURL string, designID core.Uuid) string {
	return fmt.Sprintf(designViewLinkFormat, serverURL, designID)
}

// ServerURLFromRequest returns the Meshery server URL stored on the request
// context, or "" when it is absent.
func ServerURLFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	serverURL, _ := r.Context().Value(MesheryServerURL).(string)
	return serverURL
}

// DesignEventIdentity returns the identity fields every design-related event
// carries: designId, designName, designLink and, when non-empty,
// designVersion. Callers add event-specific fields to the returned map.
func DesignEventIdentity(serverURL string, designID core.Uuid, name, version string) map[string]interface{} {
	metadata := map[string]interface{}{
		EventMetadataDesignID:   designID,
		EventMetadataDesignName: name,
		EventMetadataDesignLink: DesignLink(serverURL, designID),
	}
	if version != "" {
		metadata[EventMetadataDesignVersion] = version
	}
	return metadata
}

// eventMetadataIdentityKeys are the keys an oversized metadata map is reduced to.
var eventMetadataIdentityKeys = []string{
	EventMetadataDesignID,
	EventMetadataDesignName,
	EventMetadataDesignVersion,
	EventMetadataDesignLink,
	EventMetadataHistoryTitle,
}

// GuardEventMetadataSize bounds the metadata of an event about to be
// persisted. When its JSON serialization exceeds MaxEventMetadataBytes, the
// returned event's metadata holds only the design identity fields (and
// historyTitle) that were present, and a warning naming the event's category,
// action and original size - never its contents - is logged.
//
// The event's own metadata map is never mutated: the same map is often shared
// with the copy published to the EventBroadcaster. A nil log is allowed.
func GuardEventMetadataSize(event events.Event, log logger.Handler) events.Event {
	if len(event.Metadata) == 0 {
		return event
	}

	size := 0
	data, err := json.Marshal(event.Metadata)
	if err == nil {
		size = len(data)
	}
	// Metadata that cannot be serialized cannot be persisted either, so it is
	// reduced the same way as oversized metadata.
	if err == nil && size <= MaxEventMetadataBytes {
		return event
	}

	slim := make(map[string]interface{}, len(eventMetadataIdentityKeys))
	for _, key := range eventMetadataIdentityKeys {
		if value, ok := event.Metadata[key]; ok {
			slim[key] = value
		}
	}
	event.Metadata = slim

	if log != nil {
		if err != nil {
			log.Warn(fmt.Errorf("event metadata for category %q action %q could not be serialized and was reduced to design identity fields: %w", event.Category, event.Action, err))
		} else {
			log.Warn(fmt.Errorf("event metadata for category %q action %q was %d bytes, exceeding the %d byte limit, and was reduced to design identity fields", event.Category, event.Action, size, MaxEventMetadataBytes))
		}
	}
	return event
}
