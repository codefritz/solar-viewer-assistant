package nepviewer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"solar-viewer.de/modules/domain"
	"strconv"
	"strings"
	"time"
)

const baseUrl = "https://user.nepviewer.com/pv_monitor/proxy/week_sum/"

func FetchLatestData() []domain.DayReport {
	log.Println("Fetching data from the server...")

	nepUser := os.Getenv("NEP_USER")
	url := fmt.Sprintf("%s%s", baseUrl, nepUser)
	data, err := fetchData(url)
	if err != nil {
		log.Println("Error fetching data, returning fixed DayReport")
		yesterday := time.Now().AddDate(0, 0, -1)
		return []domain.DayReport{
			{ReportDate: yesterday, Energy: 0},
		}
	}
	dataMap := toMap(data)
	report := make([]domain.DayReport, 0)
	// loop through the dataMap
	for _, row := range dataMap {
		date := stringToDate(row[2].(string))
		produced := row[1].(float64)
		log.Printf("Date: %s, Produced: %f\n", date, produced)
		report = append(report, domain.DayReport{ReportDate: date, Energy: produced})
	}
	return report

}

// fetch data from the server via HTTP GET
func fetchData(url string) (string, error) {
	// Fetch the json from the url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
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
