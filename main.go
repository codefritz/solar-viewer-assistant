package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"solar-viewer.de/modules/analytics"
	"strconv"
	"strings"
	"time"
)

func main() {
	nepUser := os.Getenv("NEP_USER")
	fmt.Println("Hello, Solar Viewer!")
	fmt.Println("Fetching data from the server...")

	data := fetchData("https://user.nepviewer.com/pv_monitor/proxy/week_sum/" + nepUser)
	dataMap := toMap(data)
	// loop through the dataMap
	for _, row := range dataMap {
		date := stringToDate(row[2].(string))
		produced := row[1].(float64)
		fmt.Printf("Date: %s, Produced: %f\n", date, produced)
		analytics.UpdateEnergyHistory(date, produced)
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

func toMap(data string) DataMap {
	var dataMap DataMap
	err := json.Unmarshal([]byte(data), &dataMap)
	if err != nil {
		log.Fatal(err)
	}
	return dataMap
}

type DataMap [][]any

func stringToDate(data string) time.Time {
	// split data into day, month, year
	dateParts := strings.Split(data, ".")
	// dateParts[2] to int
	year, _ := strconv.Atoi(dateParts[2])
	year += 2000
	day, _ := strconv.Atoi(dateParts[1])
	month, _ := strconv.Atoi(dateParts[0])
	// transform the data into a time.Time object

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
