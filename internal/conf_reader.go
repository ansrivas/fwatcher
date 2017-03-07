package internal

import (
	"errors"

	"github.com/spf13/viper"
)

// GetConfig returns the config object for this project
func GetConfig(configPath string) (*viper.Viper, error) {
	viper.SetConfigFile(configPath)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, errors.New("Config file not found: " + configPath)
	}
	return viper.GetViper(), nil
}
