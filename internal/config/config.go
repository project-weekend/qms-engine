package config

import "github.com/project-weekend/qms-engine/server/config"

func LoadConfig() *config.Config {
	v := NewViper()

	var conf config.Config
	if err := v.Unmarshal(&conf); err != nil {
		panic(err)
	}

	return &conf
}
