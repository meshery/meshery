package config

import (
	"os"
)

func MutateConfigIfNeeded(
	configPath string,
	mesheryFolder string,
	token Token,
	ctx Context,
	createConfigFile func() error,
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

	if err := createConfigFile(); err != nil {
		return err
	}

	if err := AddTokenToConfig(token, configPath); err != nil {
		return err
	}

	if err := AddContextToConfig("local", ctx, configPath, true, false); err != nil {
		return err
	}

	return nil
}
