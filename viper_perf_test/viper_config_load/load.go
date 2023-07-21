package viper_config_load

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viperLoad()
}

func viperLoad() {
	viper.SetConfigFile("../config.json")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("err in viperLoad: " + err.Error())
	}
}
