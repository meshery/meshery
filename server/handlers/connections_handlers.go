package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

func (h *Handler) SaveConnection(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(fmt.Errorf("error reading request body: %v", err))
		http.Error(w, "unable to read result data", http.StatusInternalServerError)
		return
	}

	connection := models.Connection{}
	err = json.Unmarshal(bd, &connection)
	if err != nil {
		h.log.Error(fmt.Errorf("error unmarshal request body: %v", err))
		http.Error(w, "unable to parse connection data", http.StatusInternalServerError)
		return
	}

	err = provider.SaveConnection(req, &connection, "", false)
	if err != nil {
		h.log.Error(fmt.Errorf("error saving connection: %v", err))
		http.Error(w, "unable to save connection", http.StatusInternalServerError)
		return
	}

	h.log.Info("connection saved successfully")
	w.WriteHeader(http.StatusCreated)
}
