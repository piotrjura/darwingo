package config

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type FtpConfig struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type PushConfig struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
	Queue    string `json:"queue"`
}

type Config struct {
	Ftp  FtpConfig  `json:"ftp"`
	Push PushConfig `json:"push"`
}

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
