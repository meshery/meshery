package handlers

import (
	"context"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshery/server/models"
)

type Connection interface {
	Register(res http.ResponseWriter, req *http.Request)
	Status(res http.ResponseWriter, req *http.Request)
	Verify(res http.ResponseWriter, req *http.Request)
	Details(res http.ResponseWriter, req *http.Request)
	AddMetadata(res http.ResponseWriter, req *http.Request)
	Configure(res http.ResponseWriter, req *http.Request)
}

func NewConnectionInstance(ctx context.Context, provider models.Provider, log *logger.Handler, eb *models.Broadcast) Connection {
	connectionType, _ := ctx.Value(models.ConnectionKindKey).(string)
	systemID, _ := ctx.Value(models.SystemID).(*uuid.UUID)

	userID, _ := ctx.Value(models.UserID).(string)
	userUUID := uuid.FromStringOrNil(userID)	

	switch connectionType {
	case "helm":
		return &Helm{
			systemID: systemID,
			userID: &userUUID,
			log: *log,
			eb: eb,
			provider: provider,
		}
	}
	return nil
}