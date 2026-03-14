package config

import "os"

// NeedsMutation determines whether the Meshery config file needs to be created or updated based on its existence and content.
// needs to be initialized. It returns true if the config file
// does not exist or exists but is empty.
func NeedsMutation(configPath string) (bool, error) {
	stat, err := os.Stat(configPath)

	switch {
	case os.IsNotExist(err):
		return true, nil

	case err == nil && stat.Size() == 0:
		return true, nil

	case err != nil:
		return false, err
	}

	return false, nil
}

// InitDefaultConfig creates the default Meshery configuration,
// including the config directory, default token, and default context.
func InitDefaultConfig(
	configPath string,
	mesheryFolder string,
	token Token,
	ctx Context,
	createConfigFile func() error,
) error {

	stat, err := os.Stat(configPath)

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err == nil && stat.Size() > 0 {
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
