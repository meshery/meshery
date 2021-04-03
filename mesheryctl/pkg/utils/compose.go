package utils

type K8sCompose struct {
	APIVersion interface{} `yaml:"apiVersion,omitempty"`
	Kind       string      `yaml:"kind,omitempty"`
	Status     interface{} `yaml:"status,omitempty"`
	Metadata   Metadata    `yaml:"metadata,omitempty"`
	Spec       Spec        `yaml:"spec,omitempty"`
}

type Spec struct {
	Selector Selector    `yaml:"selector,omitempty"`
	Replicas int         `yaml:"replicas,omitempty"`
	Strategy interface{} `yaml:"strategy,omitempty"`
	Template Template    `yaml:"template,omitempty"`
}

type Template struct {
	Metadata Metadata     `yaml:"metadata,omitempty"`
	Spec     TemplateSpec `yaml:"spec,omitempty"`
}

type TemplateSpec struct {
	ServiceAccount string       `yaml:"serviceAccount,omitempty"`
	Containers     []Containers `yaml:"containers,omitempty"`
	RestartPolicy  string       `yaml:"restartPolicy,omitempty"`
}

type Containers struct {
	Env       []interface{}    `yaml:"env,omitempty"`
	Image     string           `yaml:"image,omitempty"`
	Name      string           `yaml:"name,omitempty"`
	Ports     []map[string]int `yaml:"ports,omitempty"`
	Resources interface{}      `yaml:"resources,omitempty"`
}

type Selector struct {
	MatchLabels Labels `yaml:"matchLabels,omitempty"`
}

type Metadata struct {
	Annotations       map[string]interface{} `mapstructure:"annotations,omitempty"`
	CreationTimestamp string                 `yaml:"creationTimestamp,omitempty"`
	Labels            Labels                 `yaml:"labels,omitempty"`
	Name              string                 `yaml:"name,omitempty"`
}

type Annotations struct {
	KomposeCmd     interface{} `yaml:"kompose.cmd,omitempty"`
	KomposeVersion interface{} `yaml:"kompose.version,omitempty"`
}

type Labels struct {
	KomposeService interface{} `yaml:"io.kompose.service,omitempty"`
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
