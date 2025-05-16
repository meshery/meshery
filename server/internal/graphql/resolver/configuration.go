package resolver

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) subscribeConfiguration(ctx context.Context, provider models.Provider, user models.User, patternSelector model.PageFilter, filterSelector model.PageFilter) (<-chan *model.ConfigurationPage, error) {
	userID, _ := uuid.FromString(user.ID)

	chp, unsubscribePatterns := r.Config.PatternChannel.Subscribe(userID)
	// Filters are not widely used, better to keep the subscription for filters disable.
	// chf, unsubscribeFilters := r.Config.FilterChannel.Subscribe(userID)

	r.Config.PatternChannel.Publish(userID, struct{}{})
	// r.Config.FilterChannel.Publish(userID, struct{}{})
	configuration := make(chan *model.ConfigurationPage, 1)

	go func(userID uuid.UUID) {
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
	}(userID)
	return configuration, nil
}
