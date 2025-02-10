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

// const baseUrl = "https://user.nepviewer.com/pv_monitor/proxy/week_sum/"
const baseUrl = "https://api.nepviewer.net/v2/site/statistics/echarts"

type RequestPayload struct {
	Types     int    `json:"types"`
	RangeDate string `json:"rangeDate"`
	Sid       string `json:"sid"`
}

func FetchLatestData() []domain.DayReport {
	log.Println("Fetching data from the server...")

	url := fmt.Sprintf(baseUrl)
	data, err := fetchData(url)
	// we add by pass here, to return a fixed DayReport
	// the call is not getting error, but the data is not correct
	// and will otherwise cause the program to crash in deserialize step.
	if err != nil {
		log.Println("Error fetching data, returning fixed DayReport")
		yesterday := time.Now().AddDate(0, 0, -1)
		return []domain.DayReport{
			{ReportDate: yesterday, Energy: 0},
		}
	}
	var response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Legend    []string `json:"legend"`
			XAxisData []string `json:"xAxisData"`
			Series    []struct {
				Stack string    `json:"stack"`
				Name  string    `json:"name"`
				Data  []float64 `json:"data"`
			} `json:"series"`
		} `json:"data"`
	}

	err = json.Unmarshal([]byte(data), &response)
	if err != nil {
		log.Printf("Error parsing data, returning fixed DayReport: %s", err)
		yesterday := time.Now().AddDate(0, 0, -1)
		return []domain.DayReport{
			{ReportDate: yesterday, Energy: 0},
		}
	}

	// Create DayReport instances from the parsed data
	report := make([]domain.DayReport, 0)
	for i, dateStr := range response.Data.XAxisData {
		date := stringToDate(dateStr)
		energy := response.Data.Series[0].Data[i]
		report = append(report, domain.DayReport{ReportDate: date, Energy: energy})
	}

	return report
}

func stringToDate(data string) time.Time {
	// Split data into day and month
	dateParts := strings.Split(data, "/")
	day, _ := strconv.Atoi(dateParts[0])
	month, _ := strconv.Atoi(dateParts[1])
	year := time.Now().Year() // Use the current year

	// Transform the data into a time.Time object
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// fetch data from the server via HTTP GET
func fetchData(url string) (string, error) {
	authToken := os.Getenv("AUTH_TOKEN")
	nepUser := os.Getenv("NEP_USER")
	sign := os.Getenv("SIGN")
	// Create the request payload
	payload := RequestPayload{
		Types:     3,
		RangeDate: "",
		Sid:       nepUser,
	}
	// Serialize the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		log.Printf("Error while creating request: %s", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf(authToken))
	req.Header.Set("sign", sign)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Origin", "https://user.nepviewer.com")
	req.Header.Set("Referer", "https://user.nepviewer.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("app", "0")
	req.Header.Set("client", "web")
	req.Header.Set("lan", "1")
	req.Header.Set("oem", "NEP")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error while fetching data: %s", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading response body: %s", err)
		return "", err
	}

	return string(body), nil
}
