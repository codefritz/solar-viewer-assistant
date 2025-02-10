package main

import (
	"log"
	"solar-viewer.de/modules/analytics"
	"solar-viewer.de/modules/nepviewer"
	"solar-viewer.de/modules/weather"
)

func main() {

	log.Println("Hello, Solar Viewer! Start updating the data...")
	dayReports := nepviewer.FetchLatestData()
	weatherReport := weather.FetchWeather()
	analytics.UpdateHistory(dayReports, weatherReport)

	log.Println(weatherReport)
	log.Println("Update is done.")
}
