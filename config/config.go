package config

import (
	"github.com/spf13/viper"
)

// Config houses information loaded from the config file.
type Cfg struct {
	HeartbeatInterval int
	HeartbeatLoop     int
	BloomfilterSize   int
}

// ReadConfig handles opening a file and creating a config object for use
// throughout the application.
func ReadConfig() *Cfg {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")

	viper.SetDefault("bfsize", 1000)
	viper.SetDefault("Heartbeatloop", 30)
	viper.SetDefault("Heartbeatinterval", 1000)

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &Cfg{
		viper.Get("heartbeatinterval").(int),
		viper.Get("heartbeatloop").(int),
		viper.Get("bfsize").(int),
	}
}
