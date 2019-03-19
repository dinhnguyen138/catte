package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Settings struct {
	ServiceIp      string
	ServerCertPath string
	ServerKeyPath  string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
}

var environments = map[string]string{
	"prod": "~/settings/prod-backend.json",
	"dev":  "settings/dev.json",
}

var settings Settings = Settings{}
var env = "dev"

func Init() {
	env = os.Getenv("ENV")
	if env == "" {
		fmt.Println("Warning: Undefined environment, use dev")
		env = "dev"
	}
	LoadSettings(env)
}

func LoadSettings(env string) {
	content, err := ioutil.ReadFile(environments[env])
	if err != nil {
		fmt.Println("Error while reading config file", err)
	}
	settings = Settings{}
	jsonErr := json.Unmarshal(content, &settings)
	if jsonErr != nil {
		fmt.Println("Error while parsing config file", jsonErr)
	}
}

func Get() Settings {
	if &settings == nil {
		Init()
	}
	return settings
}
