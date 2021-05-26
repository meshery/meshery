// Package handlers Meshery API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
// This should demonstrate all the possible comment annotations
// that are available to turn go code into a fully compliant swagger 2.0 spec
//
//
//     Schemes: http
//     Host: localhost:30611
//     BasePath: /
//     Version: 0.4.27
//     License: Apache-2.0 http://www.apache.org/licenses/LICENSE-2.0.txt
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - token:
//
//     SecurityDefinitions:
//     token:
//          type: JWT
//          name: token
//          in: cookie
//
//
// swagger:meta
package handlers

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/models"
)

// Returns all meshery patterns
// swagger:response mesheryPatternsResponseWrapper
type mesheryPatternsResponseWrapper struct {
	// in: body
	Body models.PatternsAPIResponse
}

// Returns a single meshery pattern
// swagger:response mesheryPatternResponseWrapper
type mesheryPatternResponseWrapper struct {
	// in: body
	Body models.MesheryPattern
}

// swagger:response noContentWrapper
type noContentWrapper struct {
}

// swagger:parameters idGetMesheryPattern idDeleteMesheryPattern
type IDParameterWrapper struct {
	// id to fetch the item
	// in: path
	// required: true
	ID uuid.UUID `json:"id"`
}
