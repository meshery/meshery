package models

type MesheryCtlConfig struct {
	Contexts       map[string]Context `json:"contexts"`
	CurrentContext string             `json:"current-context"`
	Tokens         map[string]Token   `json:"tokens"`
}

type Token struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Context struct {
	Endpoint string   `json:"endpoint"`
	Token    string   `json:"token"`
	Platform string   `json:"platform"`
	Adapters []string `json:"adapters,omitempty"`
}
