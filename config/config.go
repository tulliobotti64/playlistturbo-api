package config

import (
	"encoding/json"
	"log"
	"os"
)

var Config Configuration

type Configuration struct {
	BaseURL             string
	Server              Server
	Database            DBConfig
	StartupJobs         bool
	SupportedExtensions []string
	DlnaGenreUrl        string
}

type DBConfig struct {
	DBAddr       string
	DBPort       string
	DBName       string
	DBUser       string
	DBPass       string
	DBType       string
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
