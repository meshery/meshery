package utils

type DockerCompose struct {
	Version  string             `yaml:"version,omitempty"`
	Services map[string]Service `yaml:"services,omitempty"`
	Volumes  Volumes            `yaml:"volumes,omitempty"`
}
type Service struct {
	Image       string   `yaml:"image,omitempty"`
	Labels      []string `yaml:"labels,omitempty"`
	Environment []string `yaml:"environment,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
}

type Volumes struct {
	MesheryConfig interface{} `yaml:"meshery-config,omitempty"`
}
