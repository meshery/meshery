package config

import "os"

// What needs to be mutated
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

// Mutation logic
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

	if err := AddContextToConfig("local", ctx, configPath, true, false); err != nil {
		return err
	}

	return nil
}
