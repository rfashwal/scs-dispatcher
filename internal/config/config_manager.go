package config

import (
	"github.com/rfashwal/scs-utilities/config"
)

type registerConfig struct {
	config.Manager
	dbPath string
}

var instance *registerConfig

func Config() *registerConfig {
	if instance == nil {
		instance = new(registerConfig)
		instance.Init()
	}
	return instance
}
