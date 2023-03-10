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

func (r *Resolver) subscribeConfiguration(ctx context.Context, provider models.Provider, applicationSelector model.PageFilter, patternSelector model.PageFilter, filterSelector model.PageFilter) (<-chan *model.ConfigurationPage, error) {
	chp := make(chan struct{}, 1)
	cha := make(chan struct{}, 1)
	chf := make(chan struct{}, 1)
	chp <- struct{}{}
	cha <- struct{}{}
	chf <- struct{}{}
	key, err := uuid.NewV4()
	if err != nil {
		r.Log.Error(models.ErrGenerateUUID(err))
		return nil, err
	}
	
	r.Config.ConfigurationChannel.SubscribePatterns(chp, key)
	r.Config.ConfigurationChannel.SubscribeApplications(cha, key)
	r.Config.ConfigurationChannel.SubscribeFilters(chf, key)

	configuration := make(chan *model.ConfigurationPage)
	go func() {
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
					Patterns:     patterns,
					Applications: applications,
					Filters:      filters,
				}
				configuration <- conf

			case <-cha:
				applications, err = r.fetchApplications(ctx, provider, applicationSelector)
				if err != nil {
					r.Log.Error(ErrApplicationsSubscription(err))
					break
				}

				conf = &model.ConfigurationPage{
					Patterns:     patterns,
					Applications: applications,
					Filters:      filters,
				}
				configuration <- conf

			case <-chf:
				filters, err = r.fetchFilters(ctx, provider, filterSelector)
				if err != nil {
					r.Log.Error(ErrFiltersSubscription(err))
					break
				}

				conf = &model.ConfigurationPage{
					Patterns:     patterns,
					Applications: applications,
					Filters:      filters,
				}
				configuration <- conf

			case <-ctx.Done():
				close(configuration)
				r.Config.ConfigurationChannel.UnSubscribeApplications(key)
				r.Config.ConfigurationChannel.UnSubscribeFilters(key)
				r.Config.ConfigurationChannel.UnSubscribePatterns(key)
				r.Log.Info("Configuration subscription stopped")
				return
			}
		}
	}()
	return configuration, nil
}
