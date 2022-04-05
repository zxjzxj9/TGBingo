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

type Weather struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func GetWeather(city string, appid string) string {
	// e.g. q = Singapore
	url1 := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s"
	url2 := "https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s"
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
	ret, err = http.Get(fmt.Sprintf(url2, location[0].Lat, location[0].Lon, appid))
	weather := Weather{}
	err = json.NewDecoder(ret.Body).Decode(&weather)
	if err != nil {
		fmt.Println(err)
		return "weather request failed"
	}
	fmt.Println(url1, url2, ret)
	return fmt.Sprintf("It's %s in %s, temperature %3.1f °C ~ %3.1f °C, feels like %3.1f °C, pressure %d hPa, humidity %d%%",
		weather.Weather[0].Description, city+", "+weather.Sys.Country,
		weather.Main.TempMin-273.15, weather.Main.TempMax-273.15,
		weather.Main.FeelsLike-273.15, weather.Main.Pressure, weather.Main.Humidity)
}
