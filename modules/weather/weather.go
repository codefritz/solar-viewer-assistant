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

const baseurl = "https://api.openweathermap.org/data/3.0/onecall/timemachine"

func FetchWeather(dt string) domain.Weather {
	log.Printf("Fetching weather jsonData from the server %s...", baseurl)

	openWeatherKey := os.Getenv("OW_API_KEY")

	url := fmt.Sprintf("%s?lat=%s&lon=%s&dt=%s&appid=%s", baseurl, "52.5901", "13.3729", dt, openWeatherKey)
	log.Println(url)
	jsonData := fetchData(url)

	log.Printf("jsonData: %s\n", jsonData)

	var weatherData WeatherData
	err := json.Unmarshal([]byte(jsonData), &weatherData)
	if err != nil {
		log.Fatal(err)
	}

	if len(weatherData.Data) == 0 {
		log.Fatal("No weather data available")
	}

	data := weatherData.Data[0]

	clouds := data.Clouds

	return domain.Weather{
		Cloudiness:     clouds,
		ReportDate:     time.Unix(data.Dt, 0).UTC(),
		DayTimeMinutes: (data.Sunset - data.Sunrise) / 60,
	}
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
	Lat            float64 `json:"lat"`
	Lon            float64 `json:"lon"`
	Timezone       string  `json:"timezone"`
	TimezoneOffset int     `json:"timezone_offset"`
	Data           []Data  `json:"data"`
}

type Data struct {
	Dt      int64 `json:"dt"`
	Sunrise int   `json:"sunrise"`
	Sunset  int   `json:"sunset"`
	Clouds  int   `json:"clouds"`
}
