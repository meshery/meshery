package utils

type K8sCompose struct {
	Name      string   `yaml:"spec.template.spec.containers.name,omitempty"`
	Image     string   `yaml:"spec.template.spec.containers.image,omitempty"`
	Ports     []string `yaml:"spec.template.spec.containers.ports,omitempty"`
	Resources []string `yaml:"spec.template.spec.containers.resources,omitempty"`
}

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
