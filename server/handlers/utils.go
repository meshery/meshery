package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/machines/kubernetes"
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
	event machines.EventType) error {
	inst, ok := smInstanceTracker.ConnectToInstanceMap[connectionID]
	if !ok {
		var err error

		inst, err = kubernetes.NewK8SMachine(machineCtx, log)
		if err != nil {
			log.Error(err)
			return err
		}
		smInstanceTracker.ConnectToInstanceMap[connectionID] = inst
	}
	err := inst.SendEvent(ctx, event, nil)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}
