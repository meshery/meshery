package resolver

import (
	"context"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
)

func (r *Resolver) subscribeConfiguration(ctx context.Context, provider models.Provider, user models.User, patternSelector model.PageFilter, filterSelector model.PageFilter) (<-chan *model.ConfigurationPage, error) {
	UserID := user.ID

	chp, unsubscribePatterns := r.Config.PatternChannel.Subscribe(UserID)
	// Filters are not widely used, better to keep the subscription for filters disable.
	// chf, unsubscribeFilters := r.Config.FilterChannel.Subscribe(UserID)

	r.Config.PatternChannel.Publish(UserID, struct{}{})
	// r.Config.FilterChannel.Publish(UserID, struct{}{})
	configuration := make(chan *model.ConfigurationPage, 1)

	go func(UserID core.Uuid) {
		defer close(configuration)

		r.Log.Info("Configuration subscription started")
		for {
			select {
			case <-ctx.Done():
				unsubscribePatterns()
				r.Log.Info("Configuration subscription stopped")
				return

			case <-chp:
				patterns, err := r.fetchPatterns(ctx, provider, patternSelector)
				if err != nil {
					r.Log.Error(ErrPatternsSubscription(err))
					break
				}

				conf := &model.ConfigurationPage{
					Patterns: patterns,
				}
				configuration <- conf

			}
		}
	}(UserID)
	return configuration, nil
}
