package models

import (
	"context"
)

type QueryTrackerInterface interface {
	AddOrFlagQuery(ctx context.Context, uuid, query string, flag bool)
	RemoveUUID(ctx context.Context, uuid string)
	GetQueriesForUUID(ctx context.Context, uuid string) map[string]bool
}
