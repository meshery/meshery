package machines

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// ExtractActorIDs safely extracts user and system IDs from context.Context
// without unchecked type assertions or nil-pointer dereferences.
//
// If models.UserCtxKey contains a valid *models.User, its ID is returned.
// If models.SystemIDKey contains a valid *core.Uuid or core.Uuid, it is returned.
// Any missing, nil, or mismatched value returns the zero value (uuid.Nil / nil)
// without panicking.
func ExtractActorIDs(ctx context.Context) (core.Uuid, *core.Uuid) {
	if ctx == nil {
		return core.Uuid(uuid.Nil), nil
	}

	var userUUID core.Uuid
	if user, ok := ctx.Value(models.UserCtxKey).(*models.User); ok && user != nil {
		userUUID = user.ID
	}

	var sysID *core.Uuid
	if sID, ok := ctx.Value(models.SystemIDKey).(*core.Uuid); ok && sID != nil {
		sysID = sID
	} else if sIDVal, ok := ctx.Value(models.SystemIDKey).(core.Uuid); ok && sIDVal != core.Uuid(uuid.Nil) {
		sysID = &sIDVal
	}

	return userUUID, sysID
}

// ExtractProvider safely retrieves models.Provider from context.Context.
// It guards against missing context keys, non-provider types, and typed-nil pointers.
func ExtractProvider(ctx context.Context) (models.Provider, bool) {
	if ctx == nil {
		return nil, false
	}
	p, ok := ctx.Value(models.ProviderCtxKey).(models.Provider)
	if !ok || p == nil {
		return nil, false
	}
	return p, true
}

// HasActorAttribution checks whether at least one valid actor identity
// (user UUID or system UUID) is present to attribute events.
func HasActorAttribution(userUUID core.Uuid, sysID *core.Uuid) bool {
	return userUUID != core.Uuid(uuid.Nil) || (sysID != nil && *sysID != core.Uuid(uuid.Nil))
}

// ApplyEventAttribution safely sets FromOwner and FromSystem
// on an events.EventBuilder only when valid, non-nil IDs are available.
// It avoids nil pointer dereferences and avoids inventing synthetic fallback identities.
func ApplyEventAttribution(eb *events.EventBuilder, userUUID core.Uuid, sysID *core.Uuid) *events.EventBuilder {
	if eb == nil {
		return nil
	}
	if userUUID != core.Uuid(uuid.Nil) {
		eb.FromOwner(userUUID)
	}
	if sysID != nil && *sysID != core.Uuid(uuid.Nil) {
		eb.FromSystem(*sysID)
	}
	return eb
}
