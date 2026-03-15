package config

import "os"

// NeedsMutation determines whether the Meshery config file needs to be created or updated based on its existence and content.
// needs to be initialized. It returns true if the config file does not exist or exists but is empty, and false otherwise. An error is returned if there is an issue accessing the file.
func NeedsMutation(configPath string) (bool, error) {
	stat, err := os.Stat(configPath)

	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}

	if stat.Size() == 0 {
		return true, nil
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

	if err := os.MkdirAll(mesheryFolder, 0o700); err != nil {
		return err
	}

	if err := createConfigFile(); err != nil {
		return err
	}

	if err := AddTokenToConfig(token, configPath); err != nil {
		return err
	}
	// "local" is the default Meshery context name created during initialization.
	// It represents the local Meshery deployment endpoint and is also used
	// in CI workflows and existing Meshery configuration defaults.
	if err := AddContextToConfig("local", ctx, configPath, true, false); err != nil {
		return err
	}

	return nil
}
