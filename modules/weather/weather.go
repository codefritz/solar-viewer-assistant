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
	log.Println("Fetching weather jsonData from the server...")

	openWeatherKey := os.Getenv("OW_API_KEY")

	url := fmt.Sprintf("%s?lat=%s&lon=%s&appid=%s", baseurl, "52.5901", "13.3729", openWeatherKey)
	jsonData := fetchData(url)

	var weatherData WeatherData
	err := json.Unmarshal([]byte(jsonData), &weatherData)
	if err != nil {
		log.Fatal(err)
	}

	clouds := 0
	if weatherData.Clouds != nil {
		clouds = weatherData.Clouds.All
	}

	dt := time.Unix(weatherData.Dt, 0).UTC()
	date := dt.Format("2006-01-02")

	fmt.Printf("Dt: %s\n", dt)
	fmt.Printf("Date: %s\n", date)
	fmt.Printf("Clouds: %d\n", clouds)

	return domain.Weather{Cloudiness: clouds, ReportDate: dt}
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
}

type Clouds struct {
	All int `json:"all"`
}
