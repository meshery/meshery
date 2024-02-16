package helpers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshsync/pkg/model"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	once                      sync.Once
	autoRegistrationSingleton *AutoRegistrationHelper
)

type AutoRegistrationHelper struct {
	dbHandler         *database.Handler
	log               logger.Handler
	smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker
	eventBroadcast    *models.Broadcast
}

func InitRegistrationHelperSingleton(dbHandler *database.Handler, log logger.Handler, smInstanceTracker *machines.ConnectionToStateMachineInstanceTracker, eventBroadcast *models.Broadcast) {
	once.Do(func() {
		autoRegistrationSingleton = &AutoRegistrationHelper{
			smInstanceTracker: smInstanceTracker,
			log:               log,
			dbHandler:         dbHandler,
			eventBroadcast:    eventBroadcast,
		}
		go autoRegistrationSingleton.processRegistration()
	})
}

func GetAutoRegistrationSingleton() *AutoRegistrationHelper {
	return autoRegistrationSingleton
}

func (arh *AutoRegistrationHelper) processRegistration() {
	if arh == nil {
		return
	}
	sysID := viper.Get("INSTANCE_ID").(*uuid.UUID)
	regChan := models.GetMeshSyncRegistrationQueue().RegChan

	for regData := range regChan {
		go func(data models.MeshSyncRegistrationData) {
			// Ideally iterate all Connection defs, extract fingerprint composite key and try to match with the given obj,
			// For all connections that match the fingerprint and autoRegsiter is set to true, try to do auto registration.

			userID := data.MeshsyncDataHandler.UserID
			compMetadata := data.Obj.ComponentMetadata
			compCapabilites, ok := compMetadata["capabilities"].(map[string]interface{})
			if !ok {
				return
			}

			// MeshSync annoates the object for which connection can be created.
			isConnectionFeasible, ok := compCapabilites["connection"].(bool)
			if !ok || !isConnectionFeasible {
				return
			}

			connType := getTypeOfConnection(&data.Obj)
			if connType != "" {
				connectionDefs := arh.getConnectionDefinitions(connType)
				connectionName := FormatToTitleCase(connType)
				for _, connectionDef := range connectionDefs {
					connCapabilities, err := utils.Cast[string](connectionDef.Metadata["capabilities"])

					if err != nil {
						arh.log.Error(err)
						continue
					}
					var capabilities map[string]interface{}
					err = utils.Unmarshal(connCapabilities, &capabilities)
					if err != nil {
						arh.log.Error(models.ErrUnmarshal(err, fmt.Sprintf("Connection Definition \"%s\" capabilities", connectionDef.Kind)))
						continue
					}
					autoRegister, ok := capabilities["autoRegister"].(bool)
					if ok && autoRegister {
						urls, err := utils.Cast[[]interface{}](compCapabilites["urls"])
						if err != nil {
							return
						}

						ctx := context.WithValue(context.Background(), models.UserCtxKey, &models.User{ID: userID.String()})
						ctx = context.WithValue(ctx, models.SystemIDKey, sysID)
						ctx = context.WithValue(ctx, models.TokenCtxKey, data.MeshsyncDataHandler.Token)

						for _, url := range urls {
							connMetadata := map[string]interface{}{
								"name": data.Obj.KubernetesResourceMeta.Name,
								"url":  url,
							}

							connectionPayload, id := getConnectionPayload(connType, data.Obj.KubernetesResourceMeta.Name, data.Obj.ID, url, userID, &connectionDef, connMetadata)

							machineInst, err := InitializeMachineWithContext(connectionPayload, ctx, id, data.MeshsyncDataHandler.UserID, arh.smInstanceTracker, arh.log, data.MeshsyncDataHandler.Provider, machines.DISCOVERED, connType, nil)
							if err != nil {
								arh.log.Error(ErrAutoRegister(err, connType))
							}

							_, err = machineInst.SendEvent(ctx, machines.Register, connectionPayload)
							// If connection does not exist, transition to next states because in the "connect" event the connection will get created.
							if err != nil && !IsConnectionUpdateErr(err) {
								continue
							}
							event, err := machineInst.SendEvent(ctx, machines.Connect, connectionPayload)

							if err != nil {
								event.Description = fmt.Sprintf("Failed to auto register \"%s\" connection at %s", connectionName, url)
								// Do not publish the event if auto registration fails.
								_ = data.MeshsyncDataHandler.Provider.PersistEvent(event)
								continue
							}

							// Delete the meshsync resource which has been upgraded to Connection.
							_ = arh.dbHandler.Model(&model.KubernetesResource{}).Delete(&model.KubernetesResource{ID: data.Obj.ID})

							event = events.NewEvent().WithCategory("connection").WithAction("register").FromUser(data.MeshsyncDataHandler.UserID).ActedUpon(data.MeshsyncDataHandler.ConnectionID).WithDescription(fmt.Sprintf("Auto Registered connection of type \"%s\" at %s", connectionName, url)).Build()

							go arh.eventBroadcast.Publish(data.MeshsyncDataHandler.UserID, event)
							_ = data.MeshsyncDataHandler.Provider.PersistEvent(event)
						}
					}
				}

			}
		}(regData)
	}
}

func getConnectionPayload(connType, objName, objID string, identifier interface{}, userID uuid.UUID, connectionDef *v1alpha1.ComponentDefinition, connMetadata map[string]interface{}) (models.ConnectionPayload, uuid.UUID) {

	id, _ := generateUUID(map[string]interface{}{
		"name":       objName,
		"user_id":    userID,
		"identifier": identifier,
	})

	subCategory, _ := connectionDef.Metadata["subCategory"].(string)
	return models.ConnectionPayload{
		Kind:                       connType,
		Name:                       objName,
		Type:                       connectionDef.Model.Category.Name,
		SubType:                    subCategory,
		SkipCredentialVerification: false,
		MetaData:                   connMetadata,
		ID:                         id,
	}, id
}

func (arh *AutoRegistrationHelper) getConnectionDefinitions(connType string) []v1alpha1.ComponentDefinition {
	connectionCompFilter := &v1alpha1.ComponentFilter{
		Name:       fmt.Sprintf("%sConnection", connType),
		APIVersion: "meshery.layer5.io/v1alpha1",
		Greedy:     true,
	}

	connectionDefs, _, _ := v1alpha1.GetMeshModelComponents(arh.dbHandler, *connectionCompFilter)
	return connectionDefs
}

// Improve this fingerprinting
func getTypeOfConnection(obj *model.KubernetesResource) string {
	if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "grafana") && obj.Kind == "Service" {
		return "grafana"
	} else if strings.Contains(strings.ToLower(obj.KubernetesResourceMeta.Name), "prometheus") && obj.Kind == "Service" {
		return "prometheus"
	}
	return ""
}

func FormatToTitleCase(s string) string {
	c := cases.Title(language.English)
	return c.String(s)
}

func generateUUID(data map[string]interface{}) (uuid.UUID, error) {
	marshalledData, _ := utils.Marshal(data)
	hash := md5.Sum([]byte(marshalledData))
	return uuid.FromString(hex.EncodeToString(hash[:]))
}
