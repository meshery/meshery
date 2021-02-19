package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
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
	Contexts       map[string]Context `mapstructure:"contexts"`
	CurrentContext string             `mapstructure:"current-context"`
	Tokens         map[string]Token   `mapstructure:"tokens"`
}

// Token defines the structure of Token stored in mesheryctl
type Token struct {
	Name     string `mapstructure:"name"`
	Location string `mapstructure:"location"`
}

// Context defines a meshery environment
type Context struct {
	Endpoint string   `mapstructure:"endpoint,omitempty"`
	Token    Token    `mapstructure:"token,omitempty"`
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
		if err := ctx.SetDefaultVersion(); err != nil {
			return ctx, err
		}
		return ctx, nil
	}
	return Context{}, errors.New("current context:" + mc.CurrentContext + "does not exist")
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

// SetDefaultVersion sets version of config if not specified
func (ctx *Context) SetDefaultVersion() error {
	if ctx.Version == "" {
		// if no version is specified, pick latest imgaes based on channel
		if ctx.Channel == "edge" {
			// if channel is edge, pick the edge-latest images
			ctx.Version = "latest"
		} else if ctx.Channel == "stable" {
			// if the channel is stable, pick the version of the latest stable release
			version, err := GetLatestStableReleaseTag()
			ctx.Version = version
			if err != nil {
				return errors.Wrapf(err, fmt.Sprintf("failed to fetch latest stable release tag"))
			}
		} else {
			return errors.Errorf("unknown channel %s", ctx.Channel)
		}
	}
	return nil
}

// GetLatestStableReleaseTag fetches and returns the latest release tag from GitHub
func GetLatestStableReleaseTag() (string, error) {
	url := "https://api.github.com/repos/layer5io/meshery/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer resp.Body.Close()

	var dat map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read response body")
	}
	if err := json.Unmarshal(body, &dat); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal json into object")
	}

	return dat["tag_name"].(string), nil
}

// GetCurrentContext returns the current context name and context struct.
// If the user mentions a temporary context(tempCtxName) with -c flag, change the current-context and proceed to temporary-context
func GetCurrentContext(tempCtxName string) (string, Context, error) {
	mctlCfg, err := GetMesheryCtl(viper.GetViper())
	if err != nil {
		return "", Context{}, errors.Wrap(err, "error processing config")
	}

	if tempCtxName != "" {
		mctlCfg.CurrentContext = tempCtxName
	}

	currCtx, err := mctlCfg.CheckIfCurrentContextIsValid()
	if err != nil {
		// if the user specifies a context that is not in the config.yaml file, throw an error and show the available contexts
		log.Errorf("\n\"%s\" context does not exist. The available contexts are:", mctlCfg.CurrentContext)
		for context := range mctlCfg.Contexts {
			log.Errorf("%s", context)
		}
		return "", Context{}, errors.New("context does not exist")
	}
	return mctlCfg.CurrentContext, currCtx, nil
}
