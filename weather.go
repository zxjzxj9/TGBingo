package main

import "fmt"

func getWeather(city string) string {
	url1 := "http://api.openweathermap.org/geo/1.0/direct?q={city name},{state code},{country code}&limit={limit}&appid={API key}"
	url2 := "https://api.openweathermap.org/data/2.5/weather?lat={lat}&lon={lon}&appid={API key}"
	fmt.Println(url1, url2)
	return "It's sunny in " + city
}
