package config

import (
	"os"
)

func CreateConfigFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}
