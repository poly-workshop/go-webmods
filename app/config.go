package app

import (
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultConfigName = "default"
)

var config *viper.Viper

func Config() *viper.Viper {
	return config
}

func initConfig(configPath string) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(defaultConfigName)
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}

	viper.SetConfigName(mode)
	err = viper.MergeInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config = viper.GetViper()
}
