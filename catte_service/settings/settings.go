package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Settings struct {
	PrivateKeyPath     string
	PublicKeyPath      string
	JWTExpirationDelta int
}

var settings Settings = Settings{}

func Init() {
	LoadSettings()
}

func LoadSettings() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadFile("settings/env.json")
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
