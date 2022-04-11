package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
	return ""
}
