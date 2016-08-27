package main

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Unzip unzips a file and returns a []string of files created
func Unzip(zipFile string, unzipTo string) []string {
	result := []string{}
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Println(err.Error())
		return result
	}

	for _, f := range r.File {
		if file, err := extractFileFromZip(f, unzipTo); err != nil {
			log.Println(err)
		} else {
			result = append(result, file)
		}
	}

	return result
}

func extractFileFromZip(f *zip.File, unzipTo string) (string, error) {
	path := filepath.Join(unzipTo, f.Name)

	if f.FileInfo().IsDir() {
		os.MkdirAll(path, f.Mode())
	}

	fileReader, openErr := f.Open()
	if openErr != nil {
		return "", openErr
	}
	defer fileReader.Close()

	targetFile, targetErr := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if targetErr != nil {
		return "", targetErr
	}
	defer targetFile.Close()

	if _, copyErr := io.Copy(targetFile, fileReader); copyErr != nil {
		return "", copyErr
	}
	return path, nil
}
