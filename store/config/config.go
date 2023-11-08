package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/dexconv/hasin/store/common"
)

type Config struct {
	DbArgs      string `json:"DbArgs"`
	Port        string `json:"Port"`
	StorageSize int64  `json:"StorageSize"`
}

var (
	log         = common.Log
	GLB         = getConfig()
	filename    = "config/file/config.json"
	defaultConf = Config{
		DbArgs:      "host=postgres user=postgres password=postgres dbname=store port=5432 sslmode=disable",
		Port:        ":8080",
		StorageSize: 50,
	}
)

func getConfig() Config {
	var Conf Config
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Info("creating config file")
		_, err := os.Create(filename)
		if err != nil {
			log.Warn("error creating config file", "error", err)
		}
		configStr, err := json.Marshal(defaultConf)
		if err != nil {
			log.Warn("error marshalling config", "error", err)
		}
		if err := ioutil.WriteFile(filename, configStr, 644); err != nil {
			log.Warn("error Writing to config file", "error", err)
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		log.Warn("error reading config", "error", err)
	}

	if err := json.NewDecoder(f).Decode(&Conf); err != nil {
		log.Warn("error Decoding config", "error", err)
	}
	return Conf
}
