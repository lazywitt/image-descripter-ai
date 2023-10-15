package db

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	User     string
	Password string
	Dbname   string
	Port     string
}

const (
	pgdbConfigName string = "pgdb-dev"
	configDirectory   string = "configs"
)

func GetConf() *Config {
	conf := &Config{}

	currDirectory, err := os.Getwd()
	if err != nil {
		fmt.Printf("error getting the current directory: %v", err)
	}
	configDirectory := filepath.Join(currDirectory, configDirectory)

	viper.AddConfigPath(configDirectory)
	viper.SetConfigName(pgdbConfigName)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("%v", err)
	}

	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into Config struct, %v", err)
	}

	return conf
}