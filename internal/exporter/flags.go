package exporter

import (
	"flag"
	"fmt"
	"time"
)

func ParseFlags() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.CatalogPath, "catalog", "", "Lightroom catalog path")
	flag.StringVar(&cfg.DestinationPath, "destination", "", "Destination path")
	startDateStr := flag.String("date", "", "Start date: YYYY-MM-DD")
	endDateStr := flag.String("date_end", "", "End date: YYYY-MM-DD")
	pick := flag.Bool("pick", true, "Picked images only")
	flag.IntVar(&cfg.Rating, "rating", 0, "Minimum rating")
	flag.BoolVar(&cfg.DryRun, "dry-run", false, "Dry run")
	flag.Parse()

	if cfg.CatalogPath == "" || *startDateStr == "" {
		return nil, fmt.Errorf("'catalog' and 'date' are required.")
	}

	if *pick {
		cfg.Pick = 1
	} else {
		cfg.Pick = 0
	}

	if cfg.DestinationPath == "" {
		cfg.DestinationPath = "."
	}

	var err error
	cfg.StartDate, err = time.Parse(time.DateOnly, *startDateStr)
	if err != nil {
		return nil, err
	}

	if *endDateStr == "" {
		cfg.EndDate = cfg.StartDate
	} else {
		cfg.EndDate, err = time.Parse(time.DateOnly, *endDateStr)
		if err != nil {
			return nil, err
		}
	}

	if cfg.EndDate.Before(cfg.StartDate) {
		return nil, fmt.Errorf("'end date' must be after 'start date'")
	}

	return &cfg, nil
}
