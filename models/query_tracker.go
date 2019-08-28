package models

import (
	"context"
)

// QueryTrackerInterface defines the methods for working with query UUIDs
type QueryTrackerInterface interface {
	AddOrFlagQuery(ctx context.Context, uuid, query string, flag bool)
	RemoveUUID(ctx context.Context, uuid string)
	GetQueriesForUUID(ctx context.Context, uuid string) map[string]bool
}
