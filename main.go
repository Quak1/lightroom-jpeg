package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
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
		log.Println("Error:", err)
		return
	}

	dsn := "file:" + cfg.CatalogPath + "?mode=ro"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Error:", err)
		return
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

	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	var img adobeImage
	for rows.Next() {
		rows.Scan(&img.id, &img.path, &img.filename, &img.format, &img.sidecarExtension)
		img.path = filepath.FromSlash(img.path)

		if img.format == "RAW" && img.sidecarExtension == "" {
			log.Printf("Error: file '%s' doesn't have a sidecar image. Skipping.\n", img.path)
		} else {
			newFilename := replaceExtension(img.filename, img.sidecarExtension)
			src := filepath.Join(img.path, newFilename)
			dst := filepath.Join(cfg.DestinationPath, newFilename)

			fmt.Printf("%d: Copying '%s' into '%s'\n", img.id, src, dst)

			if cfg.Copy {
				err = copyFile(src, dst)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}
	if !srcInfo.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_EXCL, srcInfo.Mode())
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("Error: file '%s' already exists.", dst)
		}
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		dstFile.Close()
		os.Remove(dst)
		return err
	}

	return dstFile.Sync()
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
	flag.StringVar(&cfg.DestinationPath, "destination", "", "Destination path")
	startDateStr := flag.String("date", "", "Start date. Format: YYYY-MM-DD")
	endDateStr := flag.String("date_end", "", "End date. Format: YYYY-MM-DD")
	pick := flag.Bool("pick", true, "Should pictures be picked or not")
	flag.IntVar(&cfg.Rating, "rating", 0, "Min rating")
	flag.BoolVar(&cfg.Copy, "copy", false, "Copy files")
	flag.Parse()

	if cfg.CatalogPath == "" {
		return nil, fmt.Errorf("'catalog' path is required.")
	}
	if cfg.DestinationPath == "" {
		cfg.DestinationPath = "."
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
