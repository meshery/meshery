package models

type MesheryCtlConfig struct {
	Contexts       map[string]Context `yaml:"contexts"`
	CurrentContext string             `yaml:"current-context"`
	Tokens         map[string]Token   `yaml:"tokens"`
}

type Token struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

type Context struct {
	Endpoint string   `yaml:"endpoint"`
	Token    Token    `yaml:"token"`
	Platform string   `yaml:"platform"`
	Adapters []string `yaml:"adapters,omitempty"`
}
