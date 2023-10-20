package config

import (
	"encoding/json"
	"log"
	"os"
)

var Config Configuration

type Configuration struct {
	BaseFolder  string
	Database    DBConfig
	Server      Server
	StartupJobs bool
}

type DBConfig struct {
	DBType       string
	DBAddr       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPass       string
	MaxOpenConns int
}

type Server struct {
	Port string
}

func SetupConfig() {
	var raw []byte
	var err error

	if raw, err = os.ReadFile("/conf/config.json"); err != nil {
		if raw, err = os.ReadFile("config.json"); err != nil {
			log.Fatal("Unable to read configuration file: ", err)
		}
	}

	if err = json.Unmarshal(raw, &Config); err != nil {
		log.Fatal("Unable to parse configuration file: ", err)
	}
}
