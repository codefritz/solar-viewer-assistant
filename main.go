package main

import (
	"log"
	"solar-viewer.de/modules/analytics"
	"solar-viewer.de/modules/nepviewer"
)

func main() {

	log.Println("Hello, Solar Viewer!")
	dayReports := nepviewer.FetchLatestData()
	analytics.UpdateEnergyHistory(dayReports)
	log.Println("Update is done.")
}
