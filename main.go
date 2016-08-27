package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	baseURL := os.Args[1]
	tempDir := os.Args[2]
	dirListing := GetDirectoryListing(baseURL)
	ProcessZipFiles(dirListing, tempDir)
}

// ProcessZipFiles determines which files in files are zip files, then downloads them to downloadDir
func ProcessZipFiles(files []string, downloadDir string) {
	for _, file := range files {
		if strings.HasSuffix(file, ".zip") {
			// it's a zip file
			ProcessZipFile(file, downloadDir)
		}
	}
}

// ProcessZipFile unzips a file and processes each XML file contained within.
func ProcessZipFile(file string, downloadDir string) {
	fileName := fileNameFromURL(file)
	filePath := filepath.Join(downloadDir, fileName)
	if ok, err := DownloadFile(file, filePath); !ok {
		log.Println(err)
		return
	}
	// remove the zip after we're done with it.
	defer os.Remove(filePath)
	// unzip
	xmlFiles := Unzip(filePath, downloadDir)
	ProcessXMLFiles(xmlFiles)
}

// ProcessXMLFiles processes each xml file
func ProcessXMLFiles(xmlFiles []string) {
	for _, file := range xmlFiles {
		defer os.Remove(file)
		log.Println(file)
	}
}
