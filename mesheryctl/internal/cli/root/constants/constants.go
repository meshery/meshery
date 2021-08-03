package constants

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	version        = "Not Set"
	commitsha      = "Not Set"
	releasechannel = "Not Set"
	tokenPath      = "Not Set" // the token path includes location to the token for current-context
)

func GetMesheryctlVersion() string {
	return version
}

func GetMesheryctlCommitsha() string {
	return commitsha
}

func GetMesheryctlReleaseChannel() string {
	return releasechannel
}

// Function checks the location of token and returns appropriate location of the token
func GetTokenLocation(token config.Token) (string, error) {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errors.Wrap(err, "failed to get users home directory")
	}

	location := token.Location
	// if location contains /home/nectar in path we return exact location
	ok := strings.Contains(location, home)
	if ok {
		return location, nil
	}

	// else we add the home directory with the path
	return filepath.Join(utils.MesheryFolder, location), nil
}

// GetCurrentAuthToken returns location of current context token
func GetCurrentAuthToken() string {
	// get config.yaml struct
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		log.Fatal(err)
	}
	// Get token of current-context
	token, err := mctlCfg.GetTokenForContext(mctlCfg.CurrentContext)
	if err != nil {
		// Attempt to create token if it doesn't already exists
		token.Location = utils.AuthConfigFile

		// Write new entry in the config
		if err := config.AddTokenToConfig(token, utils.DefaultConfigPath); err != nil {
			log.Fatal("failed to find token path for the current context")
		}
	}
	// grab actual token location with home directory
	tokenPath, err = GetTokenLocation(token)
	if err != nil {
		log.Fatal(err)
	}

	return tokenPath
}
