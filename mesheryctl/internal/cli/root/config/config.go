package config

import (
	"fmt"
	"os"
	"path/filepath"

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
		return nil, errors.New("invalid meshconfig")
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
		log.Fatal(err)
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

// GetTokenForContext takes in the contextName and returns the token name and path corresponding
// to the given current context
func (mc *MesheryCtlConfig) GetTokenForContext(contextName string) (Token, error) {
	ctx, ok := mc.Contexts[contextName]
	if !ok {
		return Token{}, fmt.Errorf("no token is associated with context: %s", contextName)
	}

	for _, t := range mc.Tokens {
		if t.Name == ctx.Token {
			return t, nil
		}
	}

	return Token{Name: ctx.Token}, fmt.Errorf("no token found for the given context")
}

func (t *Token) GetLocation() string {
	if filepath.IsAbs(t.Location) {
		return t.Location
	}

	// If file path is not absolute then assume that the file
	// is in the .meshery directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warn("failed to get user home directory")
	}

	return filepath.Join(home, ".meshery", t.Location)
}

// GetBuild returns the build number for the binary
func (v *Version) GetBuild() string {
	return v.Build
}

// GetCommitSHA returns the commit sha for the binary
func (v *Version) GetCommitSHA() string {
	return v.CommitSHA
}

// AddTokenToConfig adds token passed to it to mesheryctl config file
func AddTokenToConfig(token Token, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	mctlCfg, err := GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	if mctlCfg.Tokens == nil {
		mctlCfg.Tokens = []Token{}
	}

	for i := range mctlCfg.Tokens {
		if mctlCfg.Tokens[i].Name == token.Name {
			return errors.New("error adding token: a token with same name already exists")
		}
	}

	mctlCfg.Tokens = append(mctlCfg.Tokens, token)

	viper.Set("contexts", mctlCfg.Contexts)
	viper.Set("current-context", mctlCfg.CurrentContext)
	viper.Set("tokens", mctlCfg.Tokens)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// DeleteTokenFromConfig deletes a token passed to it to mesheryctl config file
func DeleteTokenFromConfig(token Token, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	mctlCfg, err := GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	if mctlCfg.Tokens == nil {
		mctlCfg.Tokens = []Token{}
	}

	for i := range mctlCfg.Tokens {
		if mctlCfg.Tokens[i].Name == token.Name {
			mctlCfg.Tokens = append(mctlCfg.Tokens[:i], mctlCfg.Tokens[i+1:]...)
			viper.Set("contexts", mctlCfg.Contexts)
			viper.Set("current-context", mctlCfg.CurrentContext)
			viper.Set("tokens", mctlCfg.Tokens)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}

			return nil
		}
	}
	return errors.New("no such token exists")
}

func SetTokenToConfig(token string, configPath string, ctxName string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	mctlCfg, err := GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}
	for c, context := range mctlCfg.Contexts {
		if c == ctxName {
			ctx := context
			ctx.Token = token
			mctlCfg.Contexts[c] = ctx
			viper.Set("contexts", mctlCfg.Contexts)
			err = viper.WriteConfig()
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("invalid context name")
}

// AddContextToConfig adds context passed to it to mesheryctl config file
func AddContextToConfig(contextName string, context Context, configPath string, set bool) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	mctlCfg, err := GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	if mctlCfg.Contexts == nil {
		mctlCfg.Contexts = map[string]Context{}
	}

	_, exists := mctlCfg.Contexts[contextName]
	if exists {
		return errors.New("error adding context: a context with same name already exists")
	}

	mctlCfg.Contexts[contextName] = context
	if set {
		mctlCfg.CurrentContext = contextName
	}

	viper.Set("contexts", mctlCfg.Contexts)
	viper.Set("current-context", mctlCfg.CurrentContext)
	viper.Set("tokens", mctlCfg.Tokens)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
