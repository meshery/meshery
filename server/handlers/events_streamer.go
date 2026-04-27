package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/meshes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	_events "github.com/meshery/meshkit/utils/events"
	"github.com/meshery/schemas/models/core"
)

// subscribeFunc registers ch to receive events published on eb. Pulled out
// behind a type so tests can inject a hook (e.g. a readiness signal) without
// rewriting a package-level var, which would race under t.Parallel().
type subscribeFunc func(eb *_events.EventStreamer, ch chan interface{})

func defaultSubscribeToEventStream(eb *_events.EventStreamer, ch chan interface{}) {
	eb.Subscribe(ch)
}

// eventStreamDrainTimeout bounds how long listenForCoreEvents keeps draining
// its subscriber channel after Unsubscribe. meshkit's EventStreamer.Publish
// snapshots the subscriber list under a mutex and then fans out sends in
// fresh goroutines; if any of those goroutines won the scheduler race before
// Unsubscribe ran, they will still try to send into our now-unused channel
// and block forever without a reader. Draining on exit absorbs that in-flight
// traffic; the timeout keeps the handler from waiting on an idle publisher.
const eventStreamDrainTimeout = 100 * time.Millisecond

type eventStatusPayload struct {
	Status    string       `json:"status"`
	StatusIDs []*core.Uuid `json:"ids"`
}

type statusIDs struct {
	IDs []*core.Uuid `json:"ids"`
}

func (h *Handler) GetAllEvents(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	userID := user.ID
	page, offset, limit,
		search, order, sortOnCol, status := getPaginationParams(req)
	// eventCategory :=
	filter, err := getEventFilter(req)
	if err != nil {
		h.log.Warn(err)
	}
	filter.Limit = limit
	filter.Offset = offset
	filter.Order = order
	filter.SortOn = sortOnCol
	filter.Search = search
	filter.Status = events.EventStatus(status)

	ctx := req.Context()
	token, _ := ctx.Value(models.TokenCtxKey).(string)

	e, err := provider.GetEvents(token, filter, page, userID, *h.SystemID)

	if err != nil || e == nil {
		h.log.Error(ErrGetEvents(err))
		http.Error(w, ErrGetEvents(err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(e)

	if err != nil {
		h.log.Error(models.ErrMarshal(err, "events response"))
		http.Error(w, models.ErrMarshal(err, "events response").Error(), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.Write(e)

	// eventsResult, err := provider.GetAllEvents(filter, userID, *h.SystemID)
	// eventsR ,err := provider.GetEvents(token string, page string, pageSize string, search string, order string)
	// //
	// if err != nil {
	// 	h.log.Error(ErrGetEvents(err))
	// 	http.Error(w, ErrGetEvents(err).Error(), http.StatusInternalServerError)
	// 	return
	// }

	// eventsResult.Page = pa
	// err = json.NewEncoder(w).Encode(eventsResult)
	// if err != nil {
	// 	h.log.Error(models.ErrMarshal(err, "events response"))
	// 	http.Error(w, models.ErrMarshal(err, "events response").Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func (h *Handler) GetEventTypes(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	userID := user.ID
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	eventTypes, err := provider.GetEventTypes(token, userID, *h.SystemID)

	if err != nil {
		msg := "error retrieving event categories and actions"
		h.log.Error(fmt.Errorf("%s: %w", msg, err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(eventTypes)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "event types response"))
		http.Error(w, models.ErrMarshal(err, "event types response").Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateEventStatus(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	eventID := uuid.FromStringOrNil(mux.Vars(req)["id"])
	token, _ := req.Context().Value(models.TokenCtxKey).(string)

	defer func() {
		_ = req.Body.Close()
	}()

	var reqBody map[string]interface{}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	_ = json.Unmarshal(body, &reqBody)
	status, ok := reqBody["status"].(string)
	if !ok {
		h.log.Error(ErrUpdateEvent(fmt.Errorf("unable to parse provided event status %s", status), eventID.String()))
		http.Error(w, ErrUpdateEvent(fmt.Errorf("unable to parse provided event status %s", status), eventID.String()).Error(), http.StatusInternalServerError)
		return
	}
	err = provider.UpdateEventStatus(token, eventID, status)
	if err != nil {
		_err := ErrUpdateEvent(err, eventID.String())
		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) BulkUpdateEventStatus(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {

	defer func() {
		_ = req.Body.Close()
	}()

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	var reqBody eventStatusPayload
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	_ = json.Unmarshal(body, &reqBody)
	err = provider.BulkUpdateEventStatus(token, reqBody.StatusIDs, reqBody.Status)
	if err != nil {
		_err := ErrBulkUpdateEvent(err)
		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *Handler) BulkDeleteEvent(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	defer func() {
		_ = req.Body.Close()
	}()

	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	var reqBody statusIDs
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	_ = json.Unmarshal(body, &reqBody)
	err = provider.BulkDeleteEvent(token, reqBody.IDs)
	if err != nil {
		_err := ErrBulkDeleteEvent(err)
		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	eventID := uuid.FromStringOrNil(mux.Vars(req)["id"])
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	err := provider.DeleteEvent(token, eventID)

	if err != nil {
		_err := ErrDeleteEvent(err, eventID.String())
		h.log.Error(_err)
		http.Error(w, _err.Error(), http.StatusInternalServerError)
		return
	}
}

func getEventFilter(req *http.Request) (*events.EventsFilter, error) {
	urlValues := req.URL.Query()
	category := urlValues.Get("category")
	action := urlValues.Get("action")
	severity := urlValues.Get("severity")
	acted_upon := urlValues.Get("acted_upon")

	eventFilter := &events.EventsFilter{}
	if category != "" {
		err := json.Unmarshal([]byte(category), &eventFilter.Category)
		if err != nil {
			return eventFilter, models.ErrUnmarshal(err, "event category filter")
		}
	}

	if action != "" {
		err := json.Unmarshal([]byte(action), &eventFilter.Action)
		if err != nil {
			return eventFilter, models.ErrUnmarshal(err, "event action filter")
		}
	}

	if severity != "" {
		err := json.Unmarshal([]byte(severity), &eventFilter.Severity)
		if err != nil {
			return eventFilter, models.ErrUnmarshal(err, "event severity filter")
		}
	}

	if acted_upon != "" {
		err := json.Unmarshal([]byte(acted_upon), &eventFilter.ActedUpon)
		if err != nil {
			return eventFilter, models.ErrUnmarshal(err, "event acted upon filter")
		}
	}

	return eventFilter, nil
}

// EventStreamHandler endpoint is used for streaming events to the frontend
func (h *Handler) EventStreamHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	flusher, ok := w.(http.Flusher)

	if !ok {
		h.log.Error(ErrEventStreamingNotSupported)
		http.Error(w, "Event streaming is not supported at the moment.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	notify := req.Context()

	var err error

	localMeshAdapters := map[string]*meshes.MeshClient{}
	localMeshAdaptersLock := &sync.Mutex{}

	respChan := make(chan []byte, 100)

	newAdaptersChan := make(chan *meshes.MeshClient)

	go func() {
		for mClient := range newAdaptersChan {
			h.log.Debug("received a new mesh client, listening for events")
			go func(mClient *meshes.MeshClient) {
				listenForAdapterEvents(req.Context(), mClient, respChan, h.log, p, h.config.EventBroadcaster, *h.SystemID, user.ID.String())
				_ = mClient.Close()
			}(mClient)
		}

		h.log.Debug("new adapters channel closed")
	}()
	go listenForCoreEvents(req.Context(), h.EventsBuffer, respChan, h.log, p, defaultSubscribeToEventStream)
	go writeEventStream(req.Context(), w, respChan, h.log, flusher)

	defer func() {
		closeAdapterConnections(localMeshAdaptersLock, localMeshAdapters)
		h.log.Debug("events handler closed")
	}()

STOP:
	for {
		select {
		case <-notify.Done():
			h.log.Debug("received signal to close connection and channels")
			close(newAdaptersChan)
			break STOP
		default:
			meshAdapters := prefObj.MeshAdapters
			if meshAdapters == nil {
				meshAdapters = []*models.Adapter{}
			}

			adaptersLen := len(meshAdapters)
			if adaptersLen == 0 {
				// Clear the adapter cache
				localMeshAdapters = closeAdapterConnections(localMeshAdaptersLock, localMeshAdapters)
			} else {
				localMeshAdaptersLock.Lock()
				for _, ma := range meshAdapters {
					mClient, ok := localMeshAdapters[ma.Location]
					if !ok {
						mClient, err = meshes.CreateClient(req.Context(), ma.Location)
						if err == nil {
							localMeshAdapters[ma.Location] = mClient
						}
					}
					if mClient != nil {
						_, err = mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
						if err != nil {
							_ = mClient.Close()
							delete(localMeshAdapters, ma.Location)
						} else {
							if !ok { // reusing the map check, only when ok is false a new entry will be added
								newAdaptersChan <- mClient
							}
						}
					}
				}
				localMeshAdaptersLock.Unlock()
			}
		}

		// Yield to the scheduler without going unresponsive to client
		// disconnection: a bare time.Sleep would leave the handler pinned
		// for up to 5 seconds after the request context is cancelled.
		select {
		case <-notify.Done():
		case <-time.After(5 * time.Second):
		}
	}
}

func writeEventStream(ctx context.Context, w io.Writer, respChan <-chan []byte, log logger.Handler, flusher http.Flusher) {
	for {
		select {
		case data, ok := <-respChan:
			if !ok {
				log.Debug("response channel closed")
				return
			}

			log.Debug("received new data on response channel")
			if _, err := fmt.Fprintf(w, "data: %s\n\n", data); err != nil {
				// A write failure here almost always means the client
				// disconnected (broken pipe) — stop the loop so we don't
				// spin publishing into a dead socket.
				log.Error(fmt.Errorf("failed to write event stream: %w", err))
				return
			}
			if flusher != nil {
				flusher.Flush()
				log.Debug("Flushed the messages on the wire...")
			}
		case <-ctx.Done():
			log.Debug("client disconnected, stopping flusher loop")
			return
		}
	}
}

func sendStreamEvent(ctx context.Context, respChan chan<- []byte, data []byte) bool {
	select {
	case respChan <- data:
		return true
	case <-ctx.Done():
		return false
	}
}

func listenForCoreEvents(ctx context.Context, eb *_events.EventStreamer, resp chan []byte, log logger.Handler, _ models.Provider, subscribe subscribeFunc) {
	// Subscribe synchronously so the subscription is registered before any
	// events can be published to datach — running Subscribe in a goroutine
	// left a window in which early Publish calls could be dropped.
	datach := make(chan interface{}, 10)
	subscribe(eb, datach)
	defer func() {
		eb.Unsubscribe(datach)
		// Drain any sender goroutines already past Publish's snapshot of
		// the subscriber list. Bounded by drainTimeout so we never block
		// the handler's return if the publisher is idle.
		drainTimer := time.NewTimer(eventStreamDrainTimeout)
		defer drainTimer.Stop()
		for {
			select {
			case <-datach:
			case <-drainTimer.C:
				return
			}
		}
	}()
	for {
		select {
		case datap, ok := <-datach:
			if !ok {
				return
			}
			event, ok := datap.(*meshes.EventsResponse)
			if !ok {
				continue
			}
			data, err := json.Marshal(event)
			if err != nil {
				log.Error(models.ErrMarshal(err, "event"))
				continue
			}
			if !sendStreamEvent(ctx, resp, data) {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}
func listenForAdapterEvents(ctx context.Context, mClient *meshes.MeshClient, respChan chan []byte, log logger.Handler, p models.Provider, ec *models.Broadcast, systemID core.Uuid, userID string) {
	log.Debug("Received a stream client...")
	token, _ := ctx.Value(models.TokenCtxKey).(string)
	userUUID := uuid.FromStringOrNil(userID)
	streamClient, err := mClient.MClient.StreamEvents(ctx, &meshes.EventsRequest{})
	if err != nil {
		log.Error(ErrStreamEvents(err))
		// errChan <- err
		// http.Error(w, "There was an error connecting to the backend to get events", http.StatusInternalServerError)
		return
	}

	for {
		log.Debug("Waiting to receive events.")
		event, err := streamClient.Recv()
		if err != nil {
			if err == io.EOF {
				log.Error(ErrStreamClient(err))
				return
			}
			log.Error(ErrStreamClient(err))
			return
		}
		// log.Debugf("received an event: %+#v", event)
		log.Debug("Received an event.")
		eventType := event.EventType.String()
		eventBuilder := events.NewEvent().FromSystem(uuid.FromStringOrNil(event.Component)).
			WithSeverity(events.Informational).WithDescription(event.Summary).WithCategory(event.ComponentName).WithAction("deploy").FromUser(userUUID)
		if strings.Contains(event.Summary, "removed") {
			eventBuilder.WithAction("undeploy")
		}

		if eventType == "ERROR" {
			err := errors.New(event.ErrorCode, errors.Alert, []string{event.Summary}, []string{event.Details}, []string{event.ProbableCause}, []string{event.SuggestedRemediation})
			eventBuilder.WithMetadata(map[string]interface{}{
				"error": err,
			})
		}

		_event := eventBuilder.Build()
		_ = p.PersistEvent(*_event, token)
		ec.Publish(userUUID, _event)

		data, err := json.Marshal(event)
		if err != nil {
			log.Error(models.ErrMarshal(err, "event"))
			return
		}
		if !sendStreamEvent(ctx, respChan, data) {
			return
		}
	}
}

func closeAdapterConnections(localMeshAdaptersLock *sync.Mutex, localMeshAdapters map[string]*meshes.MeshClient) map[string]*meshes.MeshClient {
	localMeshAdaptersLock.Lock()
	for _, mcl := range localMeshAdapters {
		_ = mcl.Close()
	}
	localMeshAdaptersLock.Unlock()

	return map[string]*meshes.MeshClient{}
}

func (h *Handler) ClientEventHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	userID := user.ID
	token, err := provider.GetProviderToken(req)
	if err != nil {
		h.log.Error(ErrRetrieveUserToken(err))
		http.Error(w, ErrRetrieveUserToken(err).Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		_ = req.Body.Close()
	}()

	var evt events.Event
	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		http.Error(w, ErrRequestBody(err).Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &evt)
	if err != nil {
		h.log.Error(models.ErrUnmarshal(err, "event"))
		http.Error(w, models.ErrUnmarshal(err, "event").Error(), http.StatusInternalServerError)
		return
	}

	if evt.ActedUpon.IsNil() || evt.Action == "" || evt.Category == "" || evt.Severity == "" {
		h.log.Error(models.ErrInvalidEventData())
		http.Error(w, models.ErrInvalidEventData().Error(), http.StatusBadRequest)
		return
	}

	eventBuilder := events.NewEvent().FromUser(userID).FromSystem(*h.SystemID).
		WithCategory(evt.Category).WithAction(evt.Action).WithSeverity(events.EventSeverity(evt.Severity)).
		WithDescription(evt.Description).WithMetadata(evt.Metadata).ActedUpon(evt.ActedUpon)

	event := eventBuilder.Build()
	err = provider.PersistEvent(*event, token)
	if err != nil {
		h.log.Error(models.ErrPersistEvent(err))
		http.Error(w, models.ErrPersistEvent(err).Error(), http.StatusInternalServerError)
		return

	}
	go h.config.EventBroadcaster.Publish(userID, event)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(event)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "event response"))
		http.Error(w, models.ErrMarshal(err, "event response").Error(), http.StatusInternalServerError)
		return
	}
}
