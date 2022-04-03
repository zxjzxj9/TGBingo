package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWeather(t *testing.T) {
	ConfigData, err := loadConfig("./config.json")
	assert.Equal(t, err, nil, "Load config file error!")
	ret := GetWeather("London", ConfigData.WeatherToken)
	fmt.Println(ret)
}
