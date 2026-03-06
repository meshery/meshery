package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// PrepareConfig loads config into Viper WITHOUT mutating state.
// Safe for read-only commands.
func PrepareConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	stat, err := os.Stat(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // read-only: missing config is OK
		}
		return err
	}

	if stat.Size() == 0 {
		return nil // empty config is OK
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	log.Debug("Using config file:", viper.ConfigFileUsed())
	return nil
}
