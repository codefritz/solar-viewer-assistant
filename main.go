package main

import (
	"fmt"
	"solar-viewer.de/modules/analytics"
	"solar-viewer.de/modules/nepviewer"
)

func main() {

	fmt.Println("Hello, Solar Viewer!")
	for _, report := range nepviewer.FetchLatestData() {
		analytics.UpdateEnergyHistory(report.ReportDate, report.Energy)
	}
	fmt.Println("Update is done.")
}
