package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// Version unmarshals the json response from the server's version api
type Version struct {
	Build          string `json:"build,omitempty"`
	CommitSHA      string `json:"commitsha,omitempty"`
	ReleaseChannel string `json:"release_channel,omitempty"`
}

// MesheryCtlConfig is configuration structure of mesheryctl with contexts
type MesheryCtlConfig struct {
	Contexts       map[string]Context `yaml:"contexts"`
	CurrentContext string             `yaml:"current-context"`
	Tokens         map[string]Token   `yaml:"tokens"`
}

// Token defines the structure of Token stored in mesheryctl
type Token struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
}

// Context defines a meshery environment
type Context struct {
	Endpoint string   `yaml:"endpoint"`
	Token    Token    `yaml:"token"`
	Platform string   `yaml:"platform"`
	Adapters []string `yaml:"adapters,omitempty"`
}

// GetMesheryCtl returns a reference to the mesheryctl configuration object.
func GetMesheryCtl(v *viper.Viper) (*MesheryCtlConfig, error) {
	c := &MesheryCtlConfig{}
	// Load the config data into the object
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}
	return c, err
}

// CheckIfCurrentContextIsValid checks if current context is valid
func (mc *MesheryCtlConfig) CheckIfCurrentContextIsValid() (Context, error) {
	currentContext := viper.GetString("current-context")
	if currentContext == "" {
		return Context{}, errors.New("current context not set")
	}

	return mc.Contexts[currentContext], nil
}

// GetBaseMesheryURL returns the base meshery server URL
func (mc *MesheryCtlConfig) GetBaseMesheryURL() string {
	currentContext, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Fatal(err)
	}

	return currentContext.Endpoint
}

// GetContextContent returns contents of the current context
func (mc *MesheryCtlConfig) GetContextContent() Context {
	currentContext, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Fatal(err)
	}

	return currentContext
}

// GetBuild returns the build number for the binary
func (v *Version) GetBuild() string {
	return v.Build
}

// GetCommitSHA returns the commit sha for the binary
func (v *Version) GetCommitSHA() string {
	return v.CommitSHA
}
