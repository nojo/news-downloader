package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		log.Println("usage: news-downloader <baseUrl> <tempDir> <redisHostNameAndPort>")
		return
	}
	baseURL := os.Args[1]
	tempDir := os.Args[2]
	redisHostNamePort := os.Args[3]
	if err := InitRedisStorage(redisHostNamePort); err != nil {
		log.Printf("Unable to connect to redis: %v\n", err.Error())
		return
	}
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
		if uploaded, _ := IsFileUploaded(file); !uploaded {
			fileContents, err := readFile(file)
			if err == nil {
				AddFileToList(file, fileContents)
			} else {
				log.Printf("unable to read file: %v\n", err.Error())
			}
		}
	}
}

func readFile(fileName string) (string, error) {
	dat, err := ioutil.ReadFile(fileName)
	return string(dat), err
}
