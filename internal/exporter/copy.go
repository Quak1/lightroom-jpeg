package exporter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

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

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL|os.O_WRONLY, srcInfo.Mode())
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("Error: file '%s' already exists.", dst)
		}
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		os.Remove(dst)
		return err
	}

	return dstFile.Sync()
}

func buildPaths(img *Image, cfg *Config) (src, dst string) {
	newFilename := img.filename
	if img.format == "RAW" {
		newFilename = replaceExtension(img.filename, img.sidecarExtension)
	}

	src = filepath.Join(img.path, newFilename)
	dst = filepath.Join(cfg.DestinationPath, newFilename)
	return src, dst
}

func replaceExtension(path, ext string) string {
	idx := strings.LastIndex(path, ".")
	if idx == -1 {
		return path + "." + ext
	}
	return path[:idx] + "." + ext
}
