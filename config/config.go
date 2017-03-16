package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

// FtpConfig is darwin FTP configuration
type FtpConfig struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// PushConfig is darwin push port configuration
type PushConfig struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	Queue    string `json:"queue"`
}

// Config contains program configuration
type Config struct {
	Ftp  FtpConfig  `json:"ftp"`
	Push PushConfig `json:"push"`
}

// ReadConfig reads a JSON configuration file config.json, parses it and if
// successful returns a Config struct
func ReadConfig() Config {
	dir, err := os.Getwd()
	fmt.Println(dir)

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/config/config.json", dir))
	if err != nil || err == io.EOF {
		panic(err)
	}

	configuration := Config{}
	err = json.Unmarshal(data, &configuration)
	if err != nil {
		panic(err)
	}

	return configuration
}
