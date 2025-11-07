package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

// NewViper is a function to load config from service-config.json in config_files directory
// You can change the implementation, for example load from env file, consul, etcd, etc
func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName("service-config")
	config.SetConfigType("json")

	// Add multiple config paths to handle different working directories
	config.AddConfigPath("./config_files")                   // From project root
	config.AddConfigPath("../config_files")                  // From subdirectory
	config.AddConfigPath("../../config_files")               // From nested subdirectory
	config.AddConfigPath(filepath.Join(".", "config_files")) // Alternative format

	err := config.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %w", err))
	}

	return config
}
