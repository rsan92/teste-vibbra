package configuration

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var Configurations *AppConfigurations

type (
	Security struct {
		SecretKey string `yaml:"secret_key"`
	}
	Application struct {
		ApiPort int `yaml:"api_port"`
	}
	Database struct {
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		Name         string `yaml:"name"`
		UserDatabase string `yaml:"user_database"`
	}
)
type AppConfigurations struct {
	Security    `yaml:"security"`
	Application `yaml:"application"`
	Database    `yaml:"database"`
}

func GetConfigurations() *AppConfigurations {
	var err error
	if Configurations == nil {
		Configurations, err = NewAppConfiguration()
		if err != nil {
			log.Fatal(err)
		}
	}
	return Configurations
}

func NewAppConfiguration() (*AppConfigurations, error) {
	configDir := os.Getenv("CONF_DIR")
	configScope := os.Getenv("SCOPE")

	config := &AppConfigurations{}

	file, err := os.Open(fmt.Sprintf("%v/%v.yml", configDir, configScope))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	yamlDecoder := yaml.NewDecoder(file)

	if err := yamlDecoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
