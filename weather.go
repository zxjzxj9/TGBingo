package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Location []struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
}

func getWeather(city string, appid string) string {
	// e.g. q = Singapore
	url1 := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s"
	url2 := "https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}"
	ret, err := http.Get(fmt.Sprintf(url1, city, appid))
	if err != nil {
		fmt.Println(err)
		return "geology location request failed"
	}
	location := Location{}
	err = json.NewDecoder(ret.Body).Decode(&location)
	if err != nil {
		fmt.Println(err)
		return "geology request decode failed"
	}
	fmt.Println(url1, url2, ret)
	return "It's sunny in " + city
}
