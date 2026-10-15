package machines

import (
	"context"
	"reflect"

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

// isNilSafe checks if an interface contains nil or an underlying nil pointer/chan/func/interface/map/slice.
// It avoids calling reflect.Value.IsNil() on non-nil-capable kinds.
func isNilSafe(i interface{}) bool {
	if i == nil {
		return true
	}
	val := reflect.ValueOf(i)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice, reflect.UnsafePointer:
		return val.IsNil()
	default:
		return false
	}
}

// ExtractProvider safely retrieves models.Provider from context.Context.
// It guards against missing context keys, non-provider types, ordinary nil interfaces,
// and typed-nil pointers that pass interface type assertions.
func ExtractProvider(ctx context.Context) (models.Provider, bool) {
	if ctx == nil {
		return nil, false
	}
	raw := ctx.Value(models.ProviderCtxKey)
	if isNilSafe(raw) {
		return nil, false
	}
	p, ok := raw.(models.Provider)
	if !ok || isNilSafe(p) {
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
