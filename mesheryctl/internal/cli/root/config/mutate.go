package config

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func MutateConfigIfNeeded(
	configPath string,
	mesheryFolder string,
	createConfigFile func(string) error,
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

	if !createDefaultConfig {
		return nil
	}

	if err := os.MkdirAll(mesheryFolder, 0o775); err != nil {
		return err
	}

	if err := createConfigFile(configPath); err != nil {
		return err
	}

	if err := AddTokenToConfig(TemplateToken, configPath); err != nil {
		if err.Error() != "error adding token: a token with same name already exists" {
			return err
		}
	}

	if err := AddContextToConfig("local", TemplateContext, configPath, true, false); err != nil {
		if err.Error() != "error adding context: a context with same name already exists" {
			return err
		}
	}

	log.Debugf("Default config file created at %s", configPath)

	return nil
}
