# lr-exporter

Small CLI tool that queries a Lightroom Classic catalog (`.lrcat`) and copies selected images.

- Filter images by:
  - Capture date
  - Pick flag
  - Minimum rating
- Fetches original image paths from catalog
- Handles RAW + JPG and JPG only workflows

## Usage

### Notes

- Lightroom must be close when running this tool
- Destination files are not overwritten
- Folder structure is not preserved

```bash
./lr-exporter \
  -catalog /path/to/catalog.lrcat \
  -destination ./export \
  -date 2024-01-01 \
  -date_end 2024-01-07 \
  -rating 3 \
  -pick=true \
```

### Flags

| Flag           | Description                     | Default         |
| -------------- | ------------------------------- | --------------- |
| `-catalog`     | Path to Lightroom `.lrcat` file |                 |
| `-destination` | Output directory                | `.`             |
| `-date`        | Start date (`YYYY-MM-DD`)       |                 |
| `-date_end`    | End date (`YYYY-MM-DD`)         |                 |
| `-pick`        | Only include picked images      | `true`          |
| `-rating`      | Minimum rating                  | `0`             |
| `-dry-run`     | Perform a dry run               | `false`         |

## DB Tables

### Adobe_images

File ids with picture information like file format, aspect ratio, color labels,
wether picked or not, ratins, etc.

- id_local | `8700`
- captureTime | `2025-11-27T12:21:42`
- colorLabels | `""`
- fileFormat | `RAW` or `JPG`
- pick | `0.0` or `1.0`
- rating | `NULL` or `1.0`...
- rootFile | `12752`

### AgLibraryFile

File ids with original filenames and folder ids

- id_local | `12752`
- folder | `12744`
- originalFilename | `ABC00138.ARW`
- sidecarExtensions | `JPG` or empty

### AgLibraryFolder

Folder ids with paths from a given Root folder

- id_local | `12744`
- pathFromRoot | `2025/11/2025-11-27/`
- rootFolder | `8692`

### AgLibraryRootFolder

Root folder ids to absolute path

- id_local | `8692`
- absolutePath | `Z:/photos/`
