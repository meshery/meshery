package utils

type K8sCompose struct {
	APIVersion interface{} `yaml:"apiVersion,omitempty"`
	Kind       string      `yaml:"kind,omitempty"`
	Metadata   interface{} `yaml:"metadata,omitempty"`
	Spec       Spec        `yaml:"spec,omitempty"`
	Status     interface{} `yaml:"status,omitempty"`
}

type Spec struct {
	Replicas int         `yaml:"replicas,omitempty"`
	Selector interface{} `yaml:"selector,omitempty"`
	Strategy interface{} `yaml:"strategy,omitempty"`
	Template Template    `yaml:"template,omitempty"`
}

type Template struct {
	Metadata interface{}  `yaml:"metadata,omitempty"`
	Spec     TemplateSpec `yaml:"spec,omitempty"`
}

type TemplateSpec struct {
	ServiceAccount     string       `yaml:"serviceAccount,omitempty"`
	ServiceAccountName string       `yaml:"serviceAccountName,omitempty"`
	Containers         []Containers `yaml:"containers,omitempty"`
	RestartPolicy      string       `yaml:"restartPolicy,omitempty"`
}

type Containers struct {
	Env             interface{}      `yaml:"env,omitempty"`
	Image           string           `yaml:"image,omitempty"`
	ImagePullPolicy string           `yaml:"imagePullPolicy"`
	Name            string           `yaml:"name,omitempty"`
	Ports           []map[string]int `yaml:"ports,omitempty"`
	Resources       interface{}      `yaml:"resources,omitempty"`
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
