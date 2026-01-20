# lightroom-jpeg

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

### AgParsedImportHash

File id global, filename, filesize, filetime

## DB Relations

- Adobe_images.rootFile = AgLibraryFile.id_local
- AgLibraryFile.folder = AgLibraryFolder.id_local
- AgLibraryFolder.rootFolder = AgLibraryRootFolder.id_local

## TODO

- Handle 3 cases:
  - RAW
  - RAW + JPEG
  - JPEG
- Require start date
- If no end date is provided use today
