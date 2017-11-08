package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/blang/semver"

	"golang.org/x/net/html"
)

func fatal(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func main() {
	url := "https://golang.org/dl/"
	re := regexp.MustCompile(`go(\d+)\.(\d+)\.(\d+)\.linux-amd64\.tar\.gz$`)

	client := http.Client{
		Timeout: 2 * time.Minute,
	}
	resp, err := client.Get(url)
	if err != nil {
		fatal(fmt.Sprintf("Request failed: %v", err))
	}
	defer resp.Body.Close()

	var versions []semver.Version

	t := html.NewTokenizer(resp.Body)
loop:
	for {
		tt := t.Next()
		switch tt {

		case html.StartTagToken:
			tok := t.Token()

			if tok.Data == "a" {
				for _, a := range tok.Attr {
					if a.Key == "href" {
						m := re.FindStringSubmatch(a.Val)
						if len(m) == 4 {
							ver, err := semver.Make(fmt.Sprintf("%v.%v.%v", m[1], m[2], m[3]))
							if err != nil {
								fmt.Printf("Failed to parse version in link: %v\n", a.Val)
							} else {
								versions = append(versions, ver)
							}
						}
						break
					}
				}
			}
		case html.ErrorToken:
			break loop
		}
	}
	if len(versions) == 0 {
		fatal("No versions found!")
	}
	semver.Sort(versions)
	fmt.Printf("%v\n", versions[len(versions)-1])
}
