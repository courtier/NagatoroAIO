package utils

import (
	"os"

	"github.com/courtier/NagatoroAIO/logger"
	"github.com/pelletier/go-toml"
)

var (
	//Config holds config values
	Config *toml.Tree
)

//LoadConfigToMemory loads config variables
func LoadConfigToMemory() {
	path := "config.toml"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.LogFatal("config not found")
	} else {
		Config, err = toml.LoadFile(path)
		if err != nil {
			logger.LogFatal("error loading config", err.Error())
		}
	}
}
