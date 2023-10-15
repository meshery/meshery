package resolver

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

var (
	patterns     *model.PatternPageResult
	applications *model.ApplicationPage
	filters      *model.FilterPage
	conf         *model.ConfigurationPage
	err          error
)

func (r *Resolver) subscribeConfiguration(ctx context.Context, provider models.Provider, user models.User, applicationSelector model.PageFilter, patternSelector model.PageFilter, filterSelector model.PageFilter) (<-chan *model.ConfigurationPage, error) {
	userID, _ := uuid.FromString(user.ID)

	cha, unsubscribeApps := r.Config.ApplicationChannel.Subscribe(userID)
	chp, unsubscribePatterns := r.Config.PatternChannel.Subscribe(userID)
	chf, unsubscribeFilters := r.Config.FilterChannel.Subscribe(userID)

	r.Config.ApplicationChannel.Publish(userID, struct{}{})
	r.Config.PatternChannel.Publish(userID, struct{}{})
	r.Config.FilterChannel.Publish(userID, struct{}{})
	configuration := make(chan *model.ConfigurationPage)
	go func(userID uuid.UUID) {
		r.Log.Info("Configuration subscription started")
		for {
			select {
			case <-chp:
				patterns, err = r.fetchPatterns(ctx, provider, patternSelector)
				if err != nil {
					r.Log.Error(ErrPatternsSubscription(err))
					break
				}

				conf = &model.ConfigurationPage{
					Patterns: patterns,
				}
				configuration <- conf

			case <-cha:
				applications, err = r.fetchApplications(ctx, provider, applicationSelector)
				if err != nil {
					r.Log.Error(ErrApplicationsSubscription(err))
					break
				}

				conf = &model.ConfigurationPage{
					Applications: applications,
				}
				configuration <- conf

			case <-chf:
				filters, err = r.fetchFilters(ctx, provider, filterSelector)
				if err != nil {
					r.Log.Error(ErrFiltersSubscription(err))
					break
				}

				conf = &model.ConfigurationPage{
					Filters: filters,
				}
				configuration <- conf

			case <-ctx.Done():
				unsubscribeApps()
				unsubscribePatterns()
				unsubscribeFilters()
				close(configuration)
				r.Log.Info("Configuration subscription stopped")
				return
			}
		}
	}(userID)
	return configuration, nil
}
