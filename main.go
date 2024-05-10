package main

import (
	"fmt"
	"solar-viewer.de/modules/analytics"
	"solar-viewer.de/modules/nepviewer"
)

func main() {

	fmt.Println("Hello, Solar Viewer!")
	dayReports := nepviewer.FetchLatestData()
	analytics.UpdateEnergyHistory(dayReports)
	fmt.Println("Update is done.")
}
