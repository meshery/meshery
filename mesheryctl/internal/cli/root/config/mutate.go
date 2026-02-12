package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func MutateConfigIfNeeded(
	configPath string,
	mesheryFolder string,
	token Token,
	ctx Context,
) error {

	stat, err := os.Stat(configPath)
	createDefaultConfig := false

	switch {
	case os.IsNotExist(err):
		createDefaultConfig = true

	case err == nil && stat.Size() == 0:
		createDefaultConfig = true

	case err != nil:
		return err
	}

	if createDefaultConfig {
		if err := os.MkdirAll(mesheryFolder, 0o775); err != nil {
			return err
		}

		if err := CreateConfigFile(configPath); err != nil {
			return err
		}

		if err := AddTokenToConfig(token, configPath); err != nil {
			return err
		}

		if err := AddContextToConfig("local", ctx, configPath, true, false); err != nil {
			return err
		}

		log.Debugf("Default config file created at %s", configPath)
	}

	return nil
}
