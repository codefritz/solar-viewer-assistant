package analytics

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"solar-viewer.de/modules/domain"
)

var db *sql.DB

func UpdateHistory(dayReports []domain.DayReport, weather domain.Weather) {
	connect()
	for _, report := range dayReports {
		updateEnergy(report)
	}

	updateWeather(weather)
}

func updateWeather(weather domain.Weather) {
	log.Printf("Store weather data for date: %s ...", weather.ReportDate.String())

	query := "INSERT INTO `weather_history` (`reporting_date`, `cloudiness`, `day_minutes`) VALUES (?, ?, ?);"
	insert, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Impossible to insert weather data: %s", err)
	}
	resp, err := insert.Exec(weather.ReportDate, weather.Cloudiness, weather.DayTimeMinutes)
	insert.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted row: ", resp)
}
func updateEnergy(entry domain.DayReport) {
	log.Printf("Store analytics data for date: %s ...", entry.ReportDate.String())

	query := `INSERT IGNORE INTO energy_history (reporting_date, energy_kw) 
       			VALUES (?, ?)
       			ON DUPLICATE KEY UPDATE energy_kw = VALUES(energy_kw);`
	insert, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("Impossible to insert energy data: %s", err)
	}
	resp, err := insert.Exec(entry.ReportDate, entry.Energy)
	insert.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Inserted row: ", resp)
}
func connect() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_URL"),
		DBName:               "defaultdb",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// For testing the connection we can query the database.
	if err := db.QueryRow("SELECT reporting_date FROM energy_history limit 1"); err.Err() != nil {
		log.Fatalf("Error while testing connection %s", err.Err())
	}

	log.Println("Connected!")
}
