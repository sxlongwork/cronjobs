package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	GOL_CONFIG *Config = &Config{}
)

type Config struct {
	Endpoints   []string `json:"endpoints"`
	DialTimeout int      `json:"dailTimeout"`
}

func InitConfig(path string) (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if err = json.Unmarshal(data, GOL_CONFIG); err != nil {
		return
	}
	// fmt.Println(*GOL_CONFIG)
	return
}
