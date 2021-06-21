package config

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"os"
)

type ConfigStruct struct {
	SecretKey string `json:"secret_key"         env:"ML_SECRET_KEY"`

	Scheme string `json:"scheme" env:"ML_SCHEME"`
	Host   string `json:"host"   env:"ML_HOST"`
	Port   string `json:"port"   env:"ML_PORT"`

	CoreScheme string `json:"core_scheme" env:"ML_CORE_SCHEME"`
	CoreHost   string `json:"core_host" env:"ML_CORE_HOST"`
	CorePort   string `json:"core_port" env:"ML_CORE_PORT"`
}

var Config ConfigStruct

func init() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	config := ConfigStruct{
		SecretKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjowfQ.3YL3HomqyVouwYO9TYPd0DKM9xcTQkU8zyVkQDLK7dM",
		Host:      "127.0.0.1",
		Port:      "6001",

		CoreHost: "127.0.0.1",
		CorePort: "6001",
	}

	Config = config
}

func parseJSON(filePath string, cfg *ConfigStruct) {
	jsonFile, err := os.Open(filePath)
	defer jsonFile.Close()

	if err != nil {
		log.Println(err)
		return
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &cfg)
}
