package handlers

import (
	"net/http"
	"strconv"

	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshery/server/models"
	"github.com/gofrs/uuid"
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

func (h *Handler) broadcastEvent(eventBuilder *events.EventBuilder, provider models.Provider, userID uuid.UUID) {
	event := eventBuilder.Build()
	_ = provider.PersistEvent(event)
	go h.config.EventBroadcaster.Publish(userID, event)
}