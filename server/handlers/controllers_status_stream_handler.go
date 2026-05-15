// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/internal/sse"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/utils"
)

// controllerStatusTickInterval matches the legacy GraphQL subscription
// resolver's 5-second poll cadence.
const controllerStatusTickInterval = 5 * time.Second

// ControllersStatusStreamHandler streams `kubernetes/controllers/status`
// updates over Server-Sent Events. It replaces the
// `subscribeMesheryControllersStatus` GraphQL subscription.
//
// The handler ticks every controllerStatusTickInterval, polls every requested
// connection's controller handlers, and emits an SSE event only when the
// per-(connection, controller) status snapshot has changed since the last
// emission. The change-detection lives here (not in the UI) so that the
// lodash `_.isEqual` dedupe disappears from the client hot path.
//
// Query params: repeat `connectionID` for every connection to watch, e.g.
//
//	?connectionID=<uuid1>&connectionID=<uuid2>
func (h *Handler) ControllersStatusStreamHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, _ models.Provider) {
	connectionIDs := req.URL.Query()["connectionID"]
	if len(connectionIDs) == 0 {
		writeJSONError(w, "at least one connectionID query parameter is required", http.StatusBadRequest)
		return
	}

	// Detect a non-Flusher writer BEFORE spawning the producer goroutine,
	// so we can return a regular JSON error instead of silently leaking the
	// goroutine. sse.Stream also rejects this case but only after the
	// caller has already spawned work that nothing will consume.
	if _, ok := w.(http.Flusher); !ok {
		writeJSONError(w, "streaming not supported by this connection", http.StatusInternalServerError)
		return
	}

	// Buffer of 4 absorbs one in-flight initial seed plus one tick worth of
	// updates without blocking the producer if the SSE writer briefly stalls.
	events := make(chan any, 4)

	go h.streamControllerStatus(req.Context(), connectionIDs, events)

	// sse.Stream may return ErrFlusherUnavailable (defensively re-checked,
	// but already filtered above) or a write error (typically client
	// disconnect). Either way, draining is handled when streamControllerStatus
	// returns and closes events; the error is informational on the wire-
	// side and is logged inside sse.Stream itself.
	_ = sse.Stream(req.Context(), w, events, h.log)
}

// streamControllerStatus runs the per-tick poll loop, computes diffs, and
// pushes batched updates onto events. It closes events on return so the SSE
// writer exits cleanly.
func (h *Handler) streamControllerStatus(ctx context.Context, connectionIDs []string, events chan<- any) {
	defer close(events)

	// statusMapPerConnection mirrors the GraphQL resolver's map of the last
	// status/version emitted per (connectionID, controller). Used for change
	// detection on every tick.
	statusMapPerConnection := make(map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion)

	// Initial seed: emit the current status for every controller of every
	// requested connection as one batched message — same shape the resolver's
	// first frame had, so the UI doesn't need a special-case for first paint.
	initial := h.collectControllerStatus(connectionIDs, statusMapPerConnection, false)
	if len(initial) > 0 {
		select {
		case events <- initial:
		case <-ctx.Done():
			return
		}
	}

	ticker := time.NewTicker(controllerStatusTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		updates := h.collectControllerStatus(connectionIDs, statusMapPerConnection, true)
		if len(updates) == 0 {
			continue
		}

		select {
		case events <- updates:
		case <-ctx.Done():
			return
		}
	}
}

// controllerObservation is one (connectionID, controller, status+version) tuple
// observed during a tick. Extracted so the diff-vs-snapshot logic can be unit
// tested in isolation from the connection-tracker / state-machine plumbing.
type controllerObservation struct {
	connectionID string
	controller   models.MesheryController
	statusAndVer models.MesheryControllerStatusAndVersion
}

// collectControllerStatus walks every requested connection, gathers the current
// (status, version) for each controller, and delegates the diff-vs-snapshot
// computation to diffControllerObservations.
func (h *Handler) collectControllerStatus(
	connectionIDs []string,
	snapshot map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion,
	diffOnly bool,
) []*model.MesheryControllersStatusListItem {
	observations := make([]controllerObservation, 0)
	for _, connectionID := range connectionIDs {
		inst, ok := h.ConnectionToStateMachineInstanceTracker.Get(uuid.FromStringOrNil(connectionID))
		if !ok || inst == nil {
			continue
		}
		machinectx, err := utils.Cast[*kubernetes.MachineCtx](inst.Context)
		if err != nil {
			h.log.Error(model.ErrMesheryControllersStatusSubscription(err))
			continue
		}
		ctrlHandlers := machinectx.MesheryCtrlsHelper.GetControllerHandlersForEachContext()
		for controller, controllerHandler := range ctrlHandlers {
			version, vErr := controllerHandler.GetVersion()
			if vErr != nil {
				h.log.Error(model.ErrMesheryControllersStatusSubscription(vErr))
			}
			observations = append(observations, controllerObservation{
				connectionID: connectionID,
				controller:   controller,
				statusAndVer: models.MesheryControllerStatusAndVersion{
					Status:  controllerHandler.GetStatus(),
					Version: version,
				},
			})
		}
	}
	return diffControllerObservations(observations, snapshot, diffOnly)
}

// diffControllerObservations updates snapshot in place with the latest
// observations and returns the items that are either new (absent from the
// prior snapshot) or whose (status, version) tuple changed since the prior
// snapshot. When diffOnly is false, every observation is emitted regardless
// of prior state — this is the initial-seed path that mirrors the legacy
// resolver's first frame.
func diffControllerObservations(
	observations []controllerObservation,
	snapshot map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion,
	diffOnly bool,
) []*model.MesheryControllersStatusListItem {
	items := make([]*model.MesheryControllersStatusListItem, 0)
	for _, obs := range observations {
		if _, ok := snapshot[obs.connectionID]; !ok {
			snapshot[obs.connectionID] = make(map[models.MesheryController]models.MesheryControllerStatusAndVersion)
		}
		prev, hasPrev := snapshot[obs.connectionID][obs.controller]
		snapshot[obs.connectionID][obs.controller] = obs.statusAndVer
		// reflect.DeepEqual instead of an or-chain on the fields so the
		// diff survives future additions to MesheryControllerStatusAndVersion.
		if diffOnly && hasPrev && reflect.DeepEqual(prev, obs.statusAndVer) {
			continue
		}
		items = append(items, &model.MesheryControllersStatusListItem{
			ConnectionID: obs.connectionID,
			Controller:   model.GetInternalController(obs.controller),
			Status:       model.GetInternalControllerStatus(obs.statusAndVer.Status),
			Version:      obs.statusAndVer.Version,
		})
	}
	return items
}
