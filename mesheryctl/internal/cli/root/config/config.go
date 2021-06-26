package config

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/layer5io/meshery/mesheryctl/pkg/constants"

	"net/http"
)

// Version unmarshals the json response from the server's version api
type Version struct {
	Build          string `json:"build,omitempty"`
	CommitSHA      string `json:"commitsha,omitempty"`
	ReleaseChannel string `json:"release_channel,omitempty"`
}

// MesheryCtlConfig is configuration structure of mesheryctl with contexts
type MesheryCtlConfig struct {
	Contexts       map[string]Context `mapstructure:"contexts"`
	CurrentContext string             `mapstructure:"current-context"`
	Tokens         []Token            `mapstructure:"tokens"`
}

// Token defines the structure of Token stored in mesheryctl
type Token struct {
	Name     string `mapstructure:"name"`
	Location string `mapstructure:"location"`
}

// Context defines a meshery environment
type Context struct {
	Endpoint string   `mapstructure:"endpoint,omitempty"`
	Token    string   `mapstructure:"token,omitempty"`
	Platform string   `mapstructure:"platform"`
	Adapters []string `mapstructure:"adapters,omitempty"`
	Channel  string   `mapstructure:"channel,omitempty"`
	Version  string   `mapstructure:"version,omitempty"`
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
	if mc.CurrentContext == "" {
		return Context{}, errors.New("current context not set")
	}

	if ctx, exists := mc.Contexts[mc.CurrentContext]; exists {
		err := ctx.ValidateVersion()

		if err != nil {
			return Context{}, err
		}

		if err == nil {
			return ctx, nil
		}
	}

	return Context{}, errors.New("current context " + mc.CurrentContext + " does not exist")
}

// ValidateVersion checks if the version is valid, if empty sets it to default value latest. Returns an error if the version is invalid.
func (ctx *Context) ValidateVersion() error {
	if ctx.Version == "" {
		ctx.Version = "latest"
		return nil
	}

	if ctx.Version == "latest" {
		return nil
	}

	url := "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/" + ctx.Version + "?recursive=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Error(cerr)
		}
	}()

	if resp.StatusCode == 404 {
		log.Fatal("version " + ctx.Version + " is not a valid Meshery release")
	}

	if err != nil {
		return errors.Wrapf(err, "failed to make GET request to %s", url)
	}

	return nil
}

// GetBaseMesheryURL returns the base meshery server URL
func (mc *MesheryCtlConfig) GetBaseMesheryURL() string {
	currentContext, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Warn(err)
	}

	return currentContext.Endpoint
}

// GetCurrentContext returns contents of the current context
func (mc *MesheryCtlConfig) GetCurrentContext() Context {
	currentContext, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Fatal(err)
	}

	return currentContext
}

// SetCurrentContext sets current context and returns contents of the current context
func (mc *MesheryCtlConfig) SetCurrentContext(contextName string) (Context, error) {
	if contextName != "" {
		mc.CurrentContext = contextName
	}
	currCtx, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Errorf(err.Error())
	}

	return currCtx, err
}

// GetBuild returns the build number for the binary
func (v *Version) GetBuild() string {
	return v.Build
}

// GetCommitSHA returns the commit sha for the binary
func (v *Version) GetCommitSHA() string {
	return v.CommitSHA
}
