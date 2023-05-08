package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/server/models"

	"github.com/ghodss/yaml"
	"github.com/open-policy-agent/opa/rego"
	"github.com/sirupsen/logrus"
)

// Define a simple policy struct
type Policy struct {
	Rules string `json:"rules"`
}

type RelationObject struct {
	DestinationID   string                 `json:"destination_id,omitempty"`
	DestinationName string                 `json:"destination_name,omitempty"`
	SourceId        string                 `json:"source_id,omitempty"`
	SourceName      string                 `json:"source_name,omitempty"`
	Port            map[string]interface{} `json:"port,omitempty"`
}

type NetworkPolicyRegoResponse struct {
	ServicePodRelationships        []RelationObject `json:"service_pod_relationships,omitempty"`
	ServiceDeploymentRelationships []RelationObject `json:"service_deployment_relationships,omitempty"`
}

// PolicyRelationshipRegoHandler gets the designFile in the bytes of Yaml and resolves the network relationship as of now.
func PolicyRelationshipRegoHandler(ctx context.Context, designFile []byte) (*NetworkPolicyRegoResponse, error) {
	// Load the policy
	policyFile, err := ioutil.ReadFile("../meshmodel/policies/network-policy.rego") // absolute path, needs to be changed
	if err != nil {
		logrus.Fatal("error reading rego file", err.Error())
		return nil, err
	}

	// Initialize the policy
	policy := Policy{string(policyFile)}

	// Initialize the rego engine with the policy
	engine, err := rego.New(
		rego.Query("data.network_policy"),
		rego.Module("policy", policy.Rules),
	).PrepareForEval(ctx)

	if err != nil {
		return nil, err
	}

	var input map[string]interface{}
	err = yaml.Unmarshal(designFile, &input)
	if err != nil {
		return nil, err
	}

	eval_reponse, err := engine.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return nil, err
	}

	// Check the result of the policy
	if len(eval_reponse) > 0 && len(eval_reponse[0].Expressions) > 0 {
		if eval_resp, ok := (eval_reponse[0].Expressions[0].Value).(map[string]interface{}); ok {
			return NewRelationPolicy(eval_resp)
		}
	}

	return nil, fmt.Errorf("failed to evaluate rego policy %s", eval_reponse)
}

func (h *Handler) HandleNetworkRelationship(
	rw http.ResponseWriter,
	r *http.Request,
	_ *models.Preference,
	_ *models.User,
	provider models.Provider,
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

	// evaluate the rego policy
	networkPolicy, err := PolicyRelationshipRegoHandler(context.Background(), body)
	if err != nil {
		logrus.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(fmt.Sprintf("failed to evaluate policy, err: %s", err)))
		return
	}

	// write the response
	ec := json.NewEncoder(rw)
	_ = ec.Encode(networkPolicy)
}

// NewRelationPolicy unmarhals the json response to the struct defined as NetworkPolicyRegoResponse or throws error in failiure, enforces type struct obey
func NewRelationPolicy(networkResponse map[string]interface{}) (*NetworkPolicyRegoResponse, error) {
	var result NetworkPolicyRegoResponse

	b, err := json.Marshal(networkResponse)
	if err != nil {
		return nil, fmt.Errorf("error marhsalling json %s", err.Error())
	}

	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json into network response struct %s", err.Error())
	}

	return &result, nil
}
