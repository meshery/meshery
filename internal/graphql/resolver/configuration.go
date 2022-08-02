package resolver

import (
	"context"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
)

var (
	patterns     *model.PatternPageResult
	applications *model.ApplicationPage
	filters      *model.FilterPage
	conf         *model.ConfigurationPage
	err          error
)

func (r *Resolver) subscribeConfiguration(ctx context.Context, provider models.Provider, selector model.PageFilter) (<-chan *model.ConfigurationPage, error) {
	if r.Config.PatternsChannel == nil {
		r.Config.PatternsChannel = make(chan struct{})
	}
	if r.Config.ApplicationsChannel == nil {
		r.Config.ApplicationsChannel = make(chan struct{})
	}
	if r.Config.FiltersChannel == nil {
		r.Config.FiltersChannel = make(chan struct{})
	}
	configuration := make(chan *model.ConfigurationPage)
	go func() {
		r.Log.Info("Configuration subscription started")
		for {
			select {
			case <-r.Config.PatternsChannel:
				patterns, err = r.fetchPatterns(ctx, provider, selector)
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

			case <-r.Config.ApplicationsChannel:
				applications, err = r.fetchApplications(ctx, provider, selector)
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

			case <-r.Config.FiltersChannel:
				filters, err = r.fetchFilters(ctx, provider, selector)
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
				close(r.Config.ApplicationsChannel)
				close(r.Config.PatternsChannel)
				close(r.Config.FiltersChannel)
				r.Log.Info("Configuration subscription stopped")
				return
			}
		}
	}()
	return configuration, nil
}
