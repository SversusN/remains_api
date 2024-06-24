package config

import (
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var appConfig Config

// Config struct to hold the app config
type Config struct {
	Server   string `yaml:"server"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// InitConfig initializes the AppConfig
func initConfig() {
	log.Println("initilizing db configuration....")
	path, _ := os.Getwd()
	appConfig = Config{}
	f := fmt.Sprint(path, "/config.yaml")
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	filebytes, _ := io.ReadAll(file)
	err = yaml.Unmarshal(filebytes, &appConfig)
	if err != nil {
		log.Fatal(err)
	}

}

// AppConfig returns the current AppConfig
func GetConfig() *Config {
	if appConfig == (Config{}) {
		initConfig()
	}
	return &appConfig
}
