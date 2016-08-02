package config

import (
	"github.com/spf13/viper"
)

// Config houses information loaded from the config file.
type Cfg struct {
	HeartbeatInterval string
	HeartbeatLoop string
	BloomfilterSize string
}

// ReadConfig handles opening a file and creating a config object for use
// throughout the application.
func ReadConfig() *Cfg {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &Cfg{
		viper.Get("mysqlpass").(string),
		viper.Get("mysqlport").(string),
		viper.Get("bfsize").(string),
	}
}
