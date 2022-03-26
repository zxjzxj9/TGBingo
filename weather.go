package main

import (
	"fmt"
	"net/http"
)

func getWeather(city string, appid string) string {
	url1 := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s"
	url2 := "https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}"
	ret, err := http.Get(fmt.Sprintf(url1, city, appid))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(url1, url2, ret)
	return "It's sunny in " + city
}
