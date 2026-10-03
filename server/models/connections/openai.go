package connections

// OpenAIConn mirrors the existing PromConn/GrafanaConn pattern.
type OpenAIConn struct {
	URL  string `json:"url,omitempty"`
	Name string `json:"name,omitempty"`
}

// OpenAICred fixes the secret-exposure gap present in PromCred/GrafanaCred:
// the raw key is never a marshalable field.
type OpenAICred struct {
	Name string `json:"name,omitempty"`
	apiKeyOrBasicAuth string
}

func NewOpenAICred(name, secret string) OpenAICred {
	return OpenAICred{Name: name, apiKeyOrBasicAuth: secret}
}

// Secret returns the raw key for internal callers only (e.g. building the
// real provider client). It is never reachable via json.Marshal.
func (c OpenAICred) Secret() string {
	return c.apiKeyOrBasicAuth
}
