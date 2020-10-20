package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	configFilePath = "resources/config.json"
)

//Config type represents the json struct to hold this applications Config.
//Example:
// {
//	"cronSchedules":["*/5 * * * *"],
//	"parentDirID": "craaaazy id here",
//  "persistentWords" : ["keep", "these", "words", "in", "names"],
//	"nameDelimiter" : "_",
//	"fileExtensions":["mp4", ".MOV"],
//	"logLevel": "3",
//	"logLocation":"C:\Users\amazingUser\allOfMyLogs",
//	"credentialsPath": "resources/superSecret/credentials.json",
//	"tokenPath": "resources/superSuperSecret/token.json",
//	"runAtLaunch": true,
//	"logToConsole": true
// }

type Config struct {
	CronSchedules   []string `json:"cronSchedules"`
	ParentDirID     string   `json:"parentDirID"`
	PersistentWords []string `json:"persistentWords"`
	NameDelimiter	string	 `json:"nameDelimiter"`
	FileExtensions  []string `json:"fileExtensions"`
	LogLevel        string   `json:"logLevel"`
	LogLocation     string   `json:"logLocation"`
	CredentialsPath string   `json:"credentialsPath"`
	TokenPath       string   `json:"tokenPath"`
	RunAtLaunch     bool     `json:"RunAtLaunch"`
	LogToConsole    bool     `json:"logToConsole"`
}

//GetConfig returns a config struct after reading config.json
func GetConfig() (*Config, error) {
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}
	config := Config{}
	err = json.Unmarshal([]byte(configFile), &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
