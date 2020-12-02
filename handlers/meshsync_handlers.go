package handlers

import (
	"net/http"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// MeshSyncHandler - handles that parses meshsync response
func (h *Handler) MeshSyncHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {

	data := `{
		"clusters":[{
			"name":"",
			"nodes":[{
				"name":"",
				"cpu":"",
				"ram":"",
				"status":""
			}],
			"namespaces":[{
				"name":"",
				"labels":[],
				"annotations":[]
			}],
			"service-meshes":[{
				"name":""
			}]
		}]
	}`

	w.Header().Set("Content-Type", "application/json")

	// err := json.NewEncoder(w).Encode(data)
	// if err != nil {
	// 	logrus.Errorf("unable to send data: %v", err)
	// }
	_, err := w.Write([]byte(data))
	if err != nil {
		logrus.Errorf("unable to send data: %v", err)
	}
}
