package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAQI(t *testing.T) {
	// test the GetAQI function
	ConfigData, err := loadConfig("./config.json")
	assert.Equal(t, err, nil, "Load config file error!")
	ret := GetAQI("Shanghai", ConfigData.WeatherToken)
	fmt.Println(ret)
}
