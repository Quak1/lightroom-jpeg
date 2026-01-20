package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type config struct {
	CatalogPath     string
	DestinationPath string
	StartDate       time.Time
	EndDate         time.Time
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	catalog, err := os.Open(cfg.CatalogPath)
	if err != nil {
		log.Fatal(err)
	}
	catalogInfo, err := catalog.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(catalogInfo.IsDir())

}

func parseFlags() (*config, error) {
	var cfg config
	flag.StringVar(&cfg.CatalogPath, "catalog", "", "Lightroom catalog path")
	flag.StringVar(&cfg.CatalogPath, "destination", "", "Destination path")
	startDateStr := flag.String("date", "", "Start date. Format: YYYY-MM-DD")
	endDateStr := flag.String("end_date", "", "End date. Format: YYYY-MM-DD")
	flag.Parse()

	if cfg.CatalogPath == "" {
		return nil, fmt.Errorf("'catalog' path is required.")
	}
	if cfg.DestinationPath == "" {
		cfg.DestinationPath = "./"
	}
	if *startDateStr == "" {
		return nil, fmt.Errorf("'date' is required.")
	}
	if *endDateStr == "" {
		endDateStr = startDateStr
	}

	startDate, err := time.Parse(time.DateOnly, *startDateStr)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(time.DateOnly, *endDateStr)
	if err != nil {
		return nil, err
	}

	cfg.StartDate = startDate
	cfg.EndDate = endDate

	return &cfg, nil
}
