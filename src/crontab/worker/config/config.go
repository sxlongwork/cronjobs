package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	GOL_CONFIG *Config = &Config{}
)

type Config struct {
	Endpoints             []string `json:"endpoints"`
	DialTimeout           int      `json:"dailTimeout"`
	MongodbUrl            string   `json:"mongodbUrl"`
	MongodbTimeout        int      `json:"mongodbTimeout"`
	MongodbName           string   `json:"mongodbName"`
	MongodbCollectionName string   `json:"mongodbCollectionName"`
	LogBatchCount         int      `json:"logBatchCount"`
	AutoCommitLogTime     int      `json:"autoCommitLogTime"`
}

func InitConfig(path string) (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		log.Printf("load config file %s error.\n", path)
		return
	}
	if err = json.Unmarshal(data, GOL_CONFIG); err != nil {
		log.Println("parse data to struct Config failed.")
		return
	}
	// fmt.Println(*GOL_CONFIG)
	log.Println("load config file success.")
	return
}
