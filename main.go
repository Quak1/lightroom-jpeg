package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	CatalogPath     string
	DestinationPath string
	StartDate       time.Time
	EndDate         time.Time
	Pick            int
	Rating          int
	Copy            bool
}

type adobeImage struct {
	id               int
	path             string
	filename         string
	format           string
	sidecarExtension string
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	dsn := "file:" + cfg.CatalogPath + "?mode=ro"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	query := `
SELECT imgs.id_local AS id, 
	CONCAT(AgLibraryRootFolder.absolutePath, AgLibraryFolder.pathFromRoot ) AS path,
	AgLibraryFile.originalFilename AS filename,
	imgs.FileFormat AS format,
	AgLibraryFile.sidecarExtensions 
FROM Adobe_images AS imgs
JOIN AgLibraryFile ON imgs.rootFile = AgLibraryFile.id_local
JOIN AgLibraryFolder ON AgLibraryFile.folder = AgLibraryFolder.id_local
JOIN AgLibraryRootFolder ON AgLibraryFolder.rootFolder = AgLibraryRootFolder.id_local
WHERE imgs.captureTime >= date('%s')
	AND imgs.captureTime <  date('%s', '+1 day')
	AND imgs.pick == %d
	AND COALESCE(imgs.rating, 0) >= %d
ORDER BY id;
`
	startDate := cfg.StartDate.Format(time.DateOnly)
	endDate := cfg.EndDate.Format(time.DateOnly)
	query = fmt.Sprintf(query, startDate, endDate, cfg.Pick, cfg.Rating)

	// fmt.Println(query)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var img adobeImage
	for rows.Next() {
		rows.Scan(&img.id, &img.path, &img.filename, &img.format, &img.sidecarExtension)

		if img.format == "RAW" && img.sidecarExtension == "" {
			fmt.Printf("Error: file '%s' doesn't have a sidecar image. Skipping.\n", img.path)
		} else {
			newFilename := replaceExtension(img.filename, img.sidecarExtension)
			src := img.path + newFilename
			dst := cfg.DestinationPath + newFilename
			fmt.Printf("%d: Copying '%s' into '%s'\n", img.id, src, dst)
			if cfg.Copy {
				copyFile(src, dst)
			}
		}
	}
}

func replaceExtension(path, ext string) string {
	idx := strings.LastIndex(path, ".")
	if idx == -1 {
		return path + "." + ext
	}

	return path[:idx] + "." + ext
}

func parseFlags() (*config, error) {
	var cfg config
	flag.StringVar(&cfg.CatalogPath, "catalog", "", "Lightroom catalog path")
	flag.StringVar(&cfg.CatalogPath, "destination", "", "Destination path")
	startDateStr := flag.String("date", "", "Start date. Format: YYYY-MM-DD")
	endDateStr := flag.String("end_date", "", "End date. Format: YYYY-MM-DD")
	pick := flag.Bool("pick", true, "Should pictures be picked or not")
	flag.IntVar(&cfg.Rating, "rating", 0, "Min rating")
	flag.BoolVar(&cfg.Copy, "copy", false, "Copy files")
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
	if *pick {
		cfg.Pick = 1
	} else {
		cfg.Pick = 0
	}

	startDate, err := time.Parse(time.DateOnly, *startDateStr)
	if err != nil {
		return nil, err
	}
	endDate, err := time.Parse(time.DateOnly, *endDateStr)
	if err != nil {
		return nil, err
	}
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("'end date' must be after 'start date'")
	}

	cfg.StartDate = startDate
	cfg.EndDate = endDate

	return &cfg, nil
}
