package config

import (
	"encoding/json"
	"os"
	"sync"
)

/* Configuration settings */
type config struct {
	Dir struct {
		Base      string `json:"base"`
		Templates string `json:"templates"`
		Uploads   string `json:"uploads"`
		Avatars   string `json:"avatars"`
	} `json:"dir"`
	Url struct {
		Gravatar  string `json:"gravatar"`
		Templates string `json:"templates"`
	} `json:"url"`
	Auth struct {
		Cookie string `json:"cookie"`
	} `json:"auth"`
	Db struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		Username string `json:"user"`
		Password string `json:"passwd"`
		Database string `json:"database"`
	} `json:"db"`
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
}

var (
	once     sync.Once
	instance *config
)

/* Load configuration file */
func LoadConfig(filename string) error {
	cfg := &config{}
	// Open config file
	configFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	// close config file
	defer configFile.Close()
	// decoder for reading json files
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&cfg)

	// Creates a new Config instance
	once.Do(func() {
		instance = cfg
	})
	return nil
}

func GetInstance() *config {
	return instance
}
