package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	"gopkg.in/yaml.v2"

	"github.com/sirupsen/logrus"
)

// swagger:route POST /api/policies/run_policy GetRegoPolicyForDesignFile idGetRegoPolicyForDesignFile
// Handle POST request for running the set of policies on the design file, the policies are picked from the policies directory and query is sent to find all the relationships around the services in the given design file
//
// responses:
// 200
func (h *Handler) GetRegoPolicyForDesignFile(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	_ models.Provider,
) {
	defer func() {
		_ = r.Body.Close()
	}()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error(ErrRequestBody(err))
		http.Error(rw, ErrRequestBody(err).Error(), http.StatusBadRequest)

		rw.WriteHeader((http.StatusBadRequest))
		return
	}

	var input core.Pattern
	err = yaml.Unmarshal((body), &input)

	if err != nil {
		http.Error(rw, ErrDecoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}

	for _, svc := range input.Services {
		svc.Settings = core.Format.DePrettify(svc.Settings, false)
	}

	data, err := yaml.Marshal(input)
	if err != nil {
		http.Error(rw, models.ErrEncoding(err, "design file").Error(), http.StatusInternalServerError)
		return
	}
	// evaluate all the rego policies in the policies directory
	networkPolicy, err := h.Rego.RegoPolicyHandler("data.meshmodel_policy", data)
	if err != nil {
		h.log.Error(ErrResolvingRegoRelationship(err))
		http.Error(rw, ErrResolvingRegoRelationship(err).Error(), http.StatusInternalServerError)
		return
	}

	// write the response
	ec := json.NewEncoder(rw)
	err = ec.Encode(networkPolicy)
	if err != nil {
		h.log.Error(models.ErrEncoding(err, "networkPolicy response"))
		http.Error(rw, models.ErrEncoding(err, "networkPolicy response").Error(), http.StatusInternalServerError)
		return
	}
}
