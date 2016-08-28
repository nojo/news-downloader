package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const httpGetTimeoutSeconds = 300

// GetDirectoryListing gets a listing of the children of the given url.
func GetDirectoryListing(urlstring string) []string {
	client := http.Client{Timeout: time.Duration(httpGetTimeoutSeconds * time.Second)}
	resp, err := client.Get(urlstring)
	if err != nil {
		log.Printf("Unable to get %v: %v\n", urlstring, err)
		return nil
	}
	if resp.StatusCode >= 400 {
		log.Printf("response code: %v\n", resp.StatusCode)
		return nil
	}
	urlFetched := resp.Request.URL

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	return findLinks(urlFetched.String(), body)
}

func fileNameFromURL(file string) string {
	parts := strings.Split(file, "/")
	return parts[len(parts)-1]
}

// DownloadFile downloads a file via a GET.  Any status >= 400 is considered an error.
func DownloadFile(file string, downloadTo string) (bool, error) {
	log.Println(file)
	client := http.Client{Timeout: time.Duration(httpGetTimeoutSeconds * time.Second)}
	resp, err := client.Get(file)
	if err != nil || resp.StatusCode >= 400 {
		return false, err
	}
	defer resp.Body.Close()

	out, err := os.Create(downloadTo)
	if err != nil {
		return false, err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return false, err
	}

	return true, nil
}

func findLinks(rootURL string, htmlBytes []byte) (urls []string) {
	urls = []string{}
	z := html.NewTokenizer(bytes.NewReader(htmlBytes))
	for {
		tt := z.Next()

		switch tt {
		case html.ErrorToken:
			// End of the document
			return
		case html.StartTagToken:
			t := z.Token()

			isA := t.Data == "a"
			if isA {
				// this is an "a" tag.  Find href attribute
				for _, a := range t.Attr {
					if a.Key == "href" {
						urls = append(urls, appendToURL(rootURL, a.Val))
					}
				}
			}
		}
	}
}

func appendToURL(rootURL string, child string) string {
	if !strings.HasSuffix(rootURL, "/") {
		rootURL = rootURL + "/"
	}
	return rootURL + child
}
