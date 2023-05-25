package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"

	"github.com/layer5io/meshkit/models/meshmodel/core/policies"
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

	// TODO: remove this hardcoding. get this from API request params
   res, _ := h.registryManager.GetEntities(&v1alpha1.PolicyFilter{
		Kind: "Network",
		SubType: "OPA",
	})
	logrus.Debug("res: ", res)

	// var policies []v1alpha1.PolicyDefinition
	var policy v1alpha1.PolicyDefinition
	for _, r := range res {
		policy, _ = r.(v1alpha1.PolicyDefinition)
		// if ok {
		// 	m := make(map[string]interface{})
		// 	_ = json.Unmarshal([]byte(policy.M))
		// }
	}

	logrus.Debugf("policy: %+v", policy)


	// evaluate all the rego policies in the policies directory
	networkPolicy, err := policies.RegoPolicyHandler(context.Background(), policy.Expression, "data.network_policy", body)
	if err != nil {
		h.log.Error(ErrResolvingRegoRelationship(err))
		http.Error(rw, ErrResolvingRegoRelationship(err).Error(), http.StatusInternalServerError)
		return
	}

	// write the response
	ec := json.NewEncoder(rw)
	err = ec.Encode(networkPolicy)
	if err != nil {
		h.log.Error(ErrEncoding(err, "networkPolicy response"))
		http.Error(rw, ErrEncoding(err, "networkPolicy response").Error(), http.StatusInternalServerError)
		return
	}
}
