package analytics

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"solar-viewer.de/modules/domain"
)

var db *sql.DB

func UpdateEnergyHistory(entry domain.DayReport) {
	log.Printf("Store analytics data for date: %s ...", entry.ReportDate.String())

	connect()

	query := "INSERT IGNORE INTO `energy_history` (`reporting_date`, `energy_kw`) VALUES (?, ?);"
	insert, err := db.Prepare(query)
	if err != nil {
		log.Fatalf("impossible insert energy data: %s", err)
	}
	resp, err := insert.Exec(entry.ReportDate, entry.Energy)
	insert.Close()

	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp.LastInsertId())
}
func connect() {
	// ssh -N -L 3306:kbatchdb.k-cloud.io:3306 acharton@shell001.ek-prod.dus1.cloud

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

	// For testing we read user mail here.
	if err := db.QueryRow("SELECT reproting_date FROM energy_history limit 1"); err != nil {
		fmt.Errorf("Error while reading %s", err)
	}

	fmt.Println("Connected!")
}
