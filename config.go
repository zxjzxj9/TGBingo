package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Server  string `json:"server""`
	TGToken string `json:"tg_token"`
}

func loadConfig(fname string) *Config {
	var config Config
	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("Error in opening config file: %v", err.Error())
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error in reading data: %v", err.Error())
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Json unmarshal failed: %v", err.Error())
	}
	return &config
}
