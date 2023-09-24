package resolver

import (
	"context"
	"database/sql"

	"github.com/layer5io/meshery/server/helpers"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
)

func (r *Resolver) getTelemetryComps(ctx context.Context, provider models.Provider, k8sContextIDs []string) ([]*model.TelemetryComp, error) {
	query := `
	SELECT rom.name, rs.attribute, rst.attribute FROM objects o LEFT JOIN resource_object_meta rom on o.id = rom.id INNER JOIN resource_specs rs on o.id = rs.id INNER JOIN resource_statuses rst on o.id = rst.id WHERE o.kind = 'Service' AND o.cluster_id IN (?);
	`
	var rows *sql.Rows
	var err error
	var ctxIDs []string

	if len(k8sContextIDs) == 0 {
		return nil, ErrEmptyCurrentK8sContext
	}

	ctxIDs = k8sContextIDs

	rows, err = provider.GetGenericPersister().Raw(query, ctxIDs).Rows()

	if err != nil {
		r.Log.Error(ErrGettingTelemetryComponents(err))
	}

	components := make([]*model.TelemetryComp, 0)
	defer rows.Close()
	for rows.Next() {
		var component model.TelemetryComp
		err := rows.Scan(&component.Name, &component.Spec, &component.Status)
		if err != nil {
			r.Log.Error(ErrGettingTelemetryComponents(err))
			return nil, err
		}

		if utils.SliceContains(helpers.TelemetryComps, component.Name) {
			components = append(components, &component)
		}
	}
	return components, nil
}
