package config

import (
	"github.com/spf13/viper"
	"log"
)

// Config houses information loaded from the config file.
type Cfg struct {
	HeartbeatInterval int
	HeartbeatLoop     int
	BloomfilterSize   int
	BaseNode          bool
	RemotePeers       []string
	ListenPort        int
	IsTesting         bool
}

// ReadConfig handles opening a file and creating a config object for use
// throughout the application.
func ReadConfig() *Cfg {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")

	viper.SetDefault("bfsize", 1000)
	viper.SetDefault("heartbeatloop", 30)
	viper.SetDefault("heartbeatinterval", 1000)
	viper.SetDefault("basenode", true)
	// By default we assume no peers because we assume we're a base node.
	viper.SetDefault("remotepeers", []string{})
	viper.SetDefault("listenport", 5454)

	err := viper.ReadInConfig()
	if err != nil {
		// If we get an error here, we just fallback to the defaults.
		log.Println("No config file found! Falling back to defaults.")
	}

	return &Cfg{
		viper.Get("heartbeatinterval").(int),
		viper.Get("heartbeatloop").(int),
		viper.Get("bfsize").(int),
		viper.GetBool("basenode"),
		viper.GetStringSlice("remotepeers"),
		viper.GetInt("listenport"),
		false,
	}
}
