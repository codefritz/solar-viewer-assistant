package main

import (
	"fmt"
	"log"
	"solar-viewer.de/modules/analytics"
	"solar-viewer.de/modules/weather"
	"strconv"
	"time"
)

func main() {

	log.Println("Hello, Solar Viewer! Start updating the data...")
	// dayReports := nepviewer.FetchLatestData()
	analytics.Connect()
	now := time.Now().UTC()
	for i := 0; i < 180; i++ {
		day := now.AddDate(0, 0, -i)
		eightAM := time.Date(day.Year(), day.Month(), day.Day(), 8, 0, 0, 0, time.UTC)
		fmt.Printf("Date: %s, Unix Timestamp: %d\n", eightAM.Format("2006-01-02 15:04:05"), eightAM.Unix())
		formatInt := strconv.FormatInt(eightAM.Unix(), 10)
		log.Println(formatInt)
		weatherReport := weather.FetchWeather(formatInt)
		log.Println(weatherReport)
		analytics.UpdateWeather(weatherReport)
		// Sleep for 3 seconds
		time.Sleep(3 * time.Second)
	}

	log.Println("Update is done.")
}
