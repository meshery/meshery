package resolver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/events"
)

func (r *Resolver) eventsResolver (ctx context.Context, provider models.Provider, user models.User) (<- chan *model.Event, error){
	userID, _ := uuid.FromString(user.ID)
	ch, unsubscribe := r.Config.EventBroadcaster.Subscribe(userID)

	eventsChan := make(chan *model.Event)
	go func(userID uuid.UUID) {
		r.Log.Info("Events Subscription started for %s", user.ID)
		for {
			select {
			case ech := <- ch:
				event := ech.(*events.Event)
				m := make(map[string]interface{})
				_ = json.Unmarshal(event.Metadata, &m)
				_event := &model.Event{
					ID: event.ID.String(),
					UserID: event.UserID.String(),
					ActedUpon: event.ActedUpon.String(),
					OperationID: event.OperationID.String(),
					Description: event.Description,
					Category: event.Category,
					Action: event.Action,
					CreatedAt: event.CreatedAt, // change graphql type for date to scalar data
					DeletedAt: event.DeletedAt,
					UpdatedAt: event.UpdatedAt,
					Metadata: m,

				}
				
				r.Log.Info(fmt.Sprintf("event received for id %s %v: ", userID, ch))
				eventsChan <- _event
			case <-ctx.Done():
				unsubscribe()
				close(eventsChan)
				
				r.Log.Info("Events Subscription stopped for %s", userID)
				return
			}
		}
	}(userID)
	return eventsChan, nil
}
