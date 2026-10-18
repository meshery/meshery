package workspaces

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

const ErrOrganizationMismatchCode = "mesheryctl-1256"

func ErrOrganizationMismatch(workspaceID, actualOrgID, providedOrgID string) error {
	return errors.New(
		ErrOrganizationMismatchCode,
		errors.Alert,
		[]string{"Provided --orgId does not match the workspace's organization"},
		[]string{fmt.Sprintf("Workspace %s belongs to organization %s, not %s", workspaceID, actualOrgID, providedOrgID)},
		[]string{"--orgId is used to verify you are updating a workspace you have access to, not to move it to a different organization"},
		[]string{"Pass the workspace's actual organization ID, or run `mesheryctl workspace view` to look it up"},
	)
}
