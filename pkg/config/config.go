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
//     "credentialsPath": "",
//     "cronSchedules":"",
//     "fileRenamePaths":"captures/**/*.mp4",
//	   "logLocation":"logs/"
// }
type Config struct {
	CredentialsPath string   `json:"credentialsPath"`
	TokenPath		string	 `json:"tokenPath"`
	CronSchedules   []string `json:"cronSchedules"`
	FileRenamePaths []string `json:"FileRenamePaths"`
	LogLevel        string   `json:"logLevel"`
	LogLocation     string   `json:"logLocation"`
	LogToConsole	bool	 `json:"logToConsole"`
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
