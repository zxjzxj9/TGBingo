package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AQI struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	List []struct {
		Main struct {
			Aqi int `json:"aqi"`
		} `json:"main"`
		Components struct {
			Co   float64 `json:"co"`
			No   int     `json:"no"`
			No2  float64 `json:"no2"`
			O3   float64 `json:"o3"`
			So2  float64 `json:"so2"`
			Pm25 float64 `json:"pm2_5"`
			Pm10 float64 `json:"pm10"`
			Nh3  float64 `json:"nh3"`
		} `json:"components"`
		Dt int `json:"dt"`
	} `json:"list"`
}

func GetAQI(city string, appid string) string {
	url1 := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s"
	url2 := "http://api.openweathermap.org/data/2.5/air_pollution?lat=%f&lon=%f&appid=%s"
	ret, err := http.Get(fmt.Sprintf(url1, city, appid))
	if err != nil {
		fmt.Println(err)
		return "geology location request failed"
	}

	// Location part is same as the one in the GetWeather function
	location := Location{}
	err = json.NewDecoder(ret.Body).Decode(&location)
	if err != nil {
		fmt.Println(err)
		return "geology request decode failed"
	}
	ret, err = http.Get(fmt.Sprintf(url2, location[0].Lat, location[0].Lon, appid))
	aqi := AQI{}
	err = json.NewDecoder(ret.Body).Decode(&aqi)
	return fmt.Sprintf("Air quality index: %d, pm2.5 index: ", aqi.List[0].Main.Aqi, aqi.List[0].Components.Pm25)
}
