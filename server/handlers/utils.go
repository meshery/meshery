package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/machines/kubernetes"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
)

const (
	defaultPageSize = 25
)

func getPaginationParams(req *http.Request) (page, offset, limit int, search, order, sortOnCol, status string) {

	urlValues := req.URL.Query()
	page, _ = strconv.Atoi(urlValues.Get("page"))
	limitstr := urlValues.Get("pagesize")
	if limitstr != "all" {
		limit, _ = strconv.Atoi(limitstr)
		if limit == 0 {
			limit = defaultPageSize
		}
	}

	search = urlValues.Get("search")
	order = urlValues.Get("order")
	sortOnCol = urlValues.Get("sort")
	status = urlValues.Get("status")

	if page <= 0 {
		page = 1
	}
	offset = (page - 1) * limit

	if sortOnCol == "" {
		sortOnCol = "updated_at"
	}
	return
}

func InitializeMachineWithContext(
	machineCtx *kubernetes.MachineCtx,
	ctx context.Context,
	connectionID uuid.UUID,
	smInstanceTracker *ConnectionToStateMachineInstanceTracker,
	log logger.Handler,
	provider models.Provider,
) (*machines.StateMachine, error) {

	inst, err := kubernetes.NewK8SMachine(connectionID.String(), log)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	inst.Provider = provider
	_, err = inst.Start(ctx, machineCtx, log, kubernetes.AssignInitialCtx)
	smInstanceTracker.ConnectToInstanceMap[connectionID] = inst
	if err != nil {
		return nil, err
	}

	return inst, nil
}
