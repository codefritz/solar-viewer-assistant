package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"solar-viewer.de/modules/domain"
	"time"
)

const baseurl = "https://api.openweathermap.org/data/2.5/weather"

func FetchWeather() domain.Weather {
	log.Printf("Fetching weather jsonData from the server %s...", baseurl)

	openWeatherKey := os.Getenv("OW_API_KEY")

	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s", baseurl, "52.5901", "13.3729", openWeatherKey)
	jsonData := fetchData(url)

	log.Printf("jsonData: %s\n", jsonData)

	var weatherData WeatherData
	err := json.Unmarshal([]byte(jsonData), &weatherData)
	if err != nil {
		log.Fatal(err)
	}

	clouds := 0
	if weatherData.Clouds != nil {
		clouds = weatherData.Clouds.All
	}

	dayTimeMinutes := (weatherData.Sys.Sunset - weatherData.Sys.Sunrise) / 60
	log.Printf("Sunrise: %d, Sunset: %d, DayTimeMinutes: %d\n", weatherData.Sys.Sunrise, weatherData.Sys.Sunset, dayTimeMinutes)

	dt := time.Unix(weatherData.Dt, 0).UTC()
	log.Printf("WeatherData: %s\n", weatherData)

	return domain.Weather{Cloudiness: clouds, ReportDate: dt, DayTimeMinutes: dayTimeMinutes}
}

// fetch data from the server via HTTP GET
func fetchData(url string) string {
	// Fetch the json from the url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

type WeatherData struct {
	Clouds *Clouds `json:"clouds"`
	Dt     int64   `json:"dt"`
	Sys    Sys     `json:"sys"`
}

type Clouds struct {
	All int `json:"all"`
}

type Sys struct {
	Sunrise int `json:"sunrise"`
	Sunset  int `json:"sunset"`
}
