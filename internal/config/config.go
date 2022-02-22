package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
)

// LogConfig holds log information
type LogConfig struct {
	Format string `toml:"format"`
	Output string `toml:"output"`
	Dir    string `toml:"log-directory"`
	Level  string `toml:"level"`
}

// Config holds all the details from config.toml passed to application
type Config struct {
	Port                int                               `toml:"port"`
	JWTSecret           string                            `toml:"jwt-secret"`
	JWTExpiryMins       int                               `toml:"jwt-expiry-mins"`
	APIKeyJWTDetailsMap map[string]map[string]interface{} `toml:"api-key-jwt-map"`
	Log                 LogConfig                         `toml:"log"`
}

// AppConf package global has values parsed from config.toml
var AppConf Config

// InitConfig Initializes AppConf
// It reads in the Config file at configPath and populates AppConf
func InitConfig(configPath string) {
	log.WithFields(log.Fields{
		"file": configPath,
	}).Info("Reading in Config File")

	if _, err := toml.DecodeFile(configPath, &AppConf); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("unable to parse config toml file")
		panic(fmt.Errorf("unable to parse config toml file"))
	}

	log.Infof("Config: %+v", AppConf)
}
