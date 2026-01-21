package exporter

import (
	"database/sql"
	"path/filepath"
	"time"
)

func queryImages(db *sql.DB, cfg *Config) (*sql.Rows, error) {
	stmt, err := db.Prepare(`
SELECT imgs.id_local AS id, 
	CONCAT(rfolder.absolutePath, folder.pathFromRoot) AS path,
	file.originalFilename AS filename,
	imgs.FileFormat AS format,
	file.sidecarExtensions 
FROM Adobe_images AS imgs
JOIN AgLibraryFile AS file ON imgs.rootFile = file.id_local
JOIN AgLibraryFolder AS folder ON file.folder = folder.id_local
JOIN AgLibraryRootFolder AS rfolder ON folder.rootFolder = rfolder.id_local
WHERE imgs.captureTime >= date(?)
	AND imgs.captureTime <  date(?, '+1 day')
	AND imgs.pick == ?
	AND COALESCE(imgs.rating, 0) >= ?
ORDER BY imgs.id_local;
`)
	if err != nil {
		return nil, err
	}

	return stmt.Query(
		cfg.StartDate.Format(time.DateOnly),
		cfg.EndDate.Format(time.DateOnly),
		cfg.Pick,
		cfg.Rating,
	)
}

func scanImage(rows *sql.Rows) (*Image, error) {
	var img Image
	err := rows.Scan(&img.id, &img.path, &img.filename, &img.format, &img.sidecarExtension)
	img.path = filepath.FromSlash(img.path)
	return &img, err
}
