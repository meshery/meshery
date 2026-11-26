package prometheus

import (
	"github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshkit/logger"
)

// newTestClient builds a Client whose transport is permitted to reach private
// addresses, for tests that drive a local httptest server. Production code must
// use New, which installs the SSRF-guarded client.
func newTestClient(baseURL, secret string, log logger.Handler) *Client {
	return NewWithHTTPClient(baseURL, secret, log, utils.NewSafeHTTPClientWithAllowPrivate(defaultTimeout, true))
}
