package exporter

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func Run(cfg *Config) error {
	dsn := "file:" + cfg.CatalogPath + "?mode=ro"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := queryImages(db, cfg)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		img, err := scanImage(rows)
		if err != nil {
			return err
		}

		if img.format == "RAW" && img.sidecarExtension == "" {
			fmt.Printf("%s doesn't have a sidecar image, skipping\n", img.filename)
			continue
		}

		src, dst := buildPaths(img, cfg)
		fmt.Printf("%d: '%s' -> '%s'\n", img.id, src, dst)

		if cfg.DryRun {
			continue
		}

		if err = copyFile(src, dst); err != nil {
			log.Println(err)
		}
	}

	return rows.Err()
}
