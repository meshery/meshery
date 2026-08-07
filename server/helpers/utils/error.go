package utils

import (
	"github.com/meshery/meshkit/errors"
)

const (
	ErrHelmChartIndexCode = "meshery-server-1473"
)

// ErrHelmChartIndex is returned when a Helm repository's index.yaml cannot be
// fetched or parsed. Without the index there is no way to know which chart
// versions exist, so callers must fail rather than guess at a version.
func ErrHelmChartIndex(url, cause string) error {
	return errors.New(
		ErrHelmChartIndexCode,
		errors.Alert,
		[]string{"Unable to read the Helm chart repository index."},
		[]string{cause},
		[]string{"The Helm repository at " + url + " is unreachable, returned an error, or served an index that is not valid YAML. Meshery Server may have no outbound network access, or the repository may be temporarily unavailable."},
		[]string{"Confirm that " + url + " is reachable from the Meshery Server pod or container, then retry. If Meshery runs behind a proxy, confirm the proxy environment variables are set on the server."},
	)
}
