package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

type configManager struct {
}

func (c configManager) Setup() {
	// Load initial config
	loadConfig()

	// Watch for changes in the config file
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
		loadConfig()
	})
}

func loadConfig() {

	log.Print("loadConfig function was triggered")

	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	viper.AutomaticEnv()

	log.Print("loadConfig function was completed")
}
