package models

import (
	"net/url"
)

func getRedirectURLForNavigatorExtension(remoteProviderProperties *ProviderProperties) string {
	redirectURL := "/"
	// This is not ideal as it only considers for 1st navigator extension.
	// The navigator extension to which redirection should happen must be available, either pre-configured with server/ via user selection
	if len(remoteProviderProperties.Extensions.Navigator) > 0 {
		href := remoteProviderProperties.Extensions.Navigator[0].Href
		redirectURL = href.URI
		if href.External != nil && !*href.External {
			redirectURL, _ = url.JoinPath("/extension", href.URI)
		}
	}
	return redirectURL
}
