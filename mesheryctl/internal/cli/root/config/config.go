package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	Endpoint   string   `mapstructure:"endpoint,omitempty"`
	Token      string   `mapstructure:"token,omitempty"`
	Platform   string   `mapstructure:"platform"`
	Components []string `mapstructure:"components,omitempty"`
	Channel    string   `mapstructure:"channel,omitempty"`
	Version    string   `mapstructure:"version,omitempty"`
}

// GetMesheryCtl returns a reference to the mesheryctl configuration object
func GetMesheryCtl(v *viper.Viper) (*MesheryCtlConfig, error) {
	c := &MesheryCtlConfig{}
	// Load the config data into the object
	err := v.Unmarshal(&c)
	if err != nil {
		return nil, errors.New("invalid meshconfig")
	}
	return c, err
}

// UpdateContextInConfig write the given context in meshconfig
func UpdateContextInConfig(v *viper.Viper, context *Context, name string) error {
	viper.Set("contexts."+name, context)
	err := viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// CheckIfCurrentContextIsValid checks if current context is valid
func (mc *MesheryCtlConfig) CheckIfCurrentContextIsValid() (*Context, error) {
	if mc.CurrentContext == "" {
		return &Context{}, errors.New("Valid context is not available in meshconfig")
	}

	ctx, exists := mc.Contexts[mc.CurrentContext]
	if exists {
		return &ctx, nil
	}

	return &Context{}, errors.New("current context " + mc.CurrentContext + " does not exist")
}
func (mc *MesheryCtlConfig) CheckIfGivenContextIsValid(name string) (*Context, error) {
	ctx, exists := mc.Contexts[name]
	if exists {
		return &ctx, nil
	}

	return &Context{}, errors.New("context " + name + " does not exist")
}

// GetBaseMesheryURL returns the base meshery server URL
func (mc *MesheryCtlConfig) GetBaseMesheryURL() string {
	currentContext, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Fatal(err)
	}

	return currentContext.Endpoint
}

func (mc *MesheryCtlConfig) GetCurrentContextName() string {
	return mc.CurrentContext
}

// GetCurrentContext returns contents of the current context
func (mc *MesheryCtlConfig) GetCurrentContext() (*Context, error) {
	currentContext, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Fatal(err)
	}

	return currentContext, err
}

// Get any context
func (mc *MesheryCtlConfig) GetContext(name string) (*Context, error) {
	context, err := mc.CheckIfGivenContextIsValid(name)
	if err != nil {
		log.Fatal(err)
	}

	return context, err
}

// SetCurrentContext sets current context and returns contents of the current context
func (mc *MesheryCtlConfig) SetCurrentContext(contextName string) error {
	if contextName != "" {
		mc.CurrentContext = contextName
	}
	_, err := mc.CheckIfCurrentContextIsValid()
	if err != nil {
		log.Errorf(err.Error())
	}

	return err
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

// GetTokens returns the tokens present in the config file
func (mc *MesheryCtlConfig) GetTokens() *[]Token {
	return &mc.Tokens
}

// GetEndpoint returns the endpoint of the current context
func (ctx *Context) GetEndpoint() string {
	return ctx.Endpoint
}

// SetEndpoint sets the endpoint of the current context
func (ctx *Context) SetEndpoint(endpoint string) {
	ctx.Endpoint = endpoint
}

// GetToken returns the token of the current context
func (ctx *Context) GetToken() string {
	return ctx.Token
}

// SetToken sets the token of the current context
func (ctx *Context) SetToken(token string) {
	ctx.Token = token
}

// GetPlatform returns the platform  of the current context
func (ctx *Context) GetPlatform() string {
	return ctx.Platform
}

// SetPlatform sets the platform of the current context
func (ctx *Context) SetPlatform(platform string) {
	ctx.Platform = platform
}

// GetComponents returns the components in the current context
func (ctx *Context) GetComponents() []string {
	return ctx.Components
}

// GetChannel returns the channel of the current context
func (ctx *Context) GetChannel() string {
	return ctx.Channel
}

// SetChannel sets the channel of the current context
func (ctx *Context) SetChannel(channel string) {
	ctx.Channel = channel
}

// GetVersion returns the version of the current context
func (ctx *Context) GetVersion() string {
	return ctx.Version
}

// SetVersion sets the version of the current context
func (ctx *Context) SetVersion(version string) {
	ctx.Version = version
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

	url := "https://github.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/releases/tag/" + ctx.Version
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
		log.Fatal("version '" + ctx.Version + "' is not a valid Meshery release.")
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatal("failed to validate Meshery release version " + ctx.Version)
	}

	if err != nil {
		return errors.Wrapf(err, "failed to make GET request to %s", url)
	}

	return nil
}

// GetName returns the token name
func (t *Token) GetName() string {
	return t.Name
}

// GetLocation returns the location of the token
func (t *Token) GetLocation() string {
	if filepath.IsAbs(t.Location) {
		return t.Location
	}

	// If file path is relative, then it has to be expanded
	if strings.HasPrefix(t.Location, "~/") {
		usr, err := os.UserHomeDir()
		if err != nil {
			log.Warn("failed to get user home directory")
		}
		return filepath.Join(usr, t.Location[2:])
	}

	// If file path is not absolute, then assume that the file
	// is in the .meshery directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Warn("failed to get user home directory")
	}

	return filepath.Join(home, ".meshery", t.Location)
}

// SetName sets the token name
func (t *Token) SetName(name string) {
	t.Name = name
}

// SetLocation sets the location of the token
func (t *Token) SetLocation(location string) {
	t.Location = location
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
func DeleteTokenFromConfig(tokenName string, configPath string) error {
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
		if mctlCfg.Tokens[i].Name == tokenName {
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

func SetTokenToConfig(tokenName string, configPath string, ctxName string) error {
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
	context, err := mctlCfg.GetContext(ctxName)
	if err != nil {
		return err
	}
	context.Token = tokenName
	err = UpdateContextInConfig(viper.GetViper(), context, ctxName)
	if err != nil {
		return err
	}
	return nil
}

// AddContextToConfig adds context passed to it to mesheryctl config file. If overwrite is set to true, existing
// context with the contextName is overwritten
func AddContextToConfig(contextName string, context Context, configPath string, set bool, overwrite bool) error {
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
	if exists && !overwrite {
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
