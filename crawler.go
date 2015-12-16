package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	//"regexp"
	"sync"
)

// Crawl pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, visited map[string]bool, wg *sync.WaitGroup) {

	fmt.Printf("Crawl %s\n", url)
	if depth <= 0 {
		fmt.Println("depth 0")
		return
	}

	url_hash := md5.Sum([]byte(url))
	url_hash_hex := hex.EncodeToString(url_hash[:])
	if visited[url_hash_hex] == true {
		fmt.Printf("skipping %s\n", url)
		return
	}

	body, urls, err := Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q %s\n", url, body, url_hash_hex)
	visited[url_hash_hex] = true

	for _, u := range urls {
		wg.Add(1)
		go func(the_url string) {
			defer wg.Done()
			Crawl(the_url, depth-1, visited, wg)
		}(u)
	}
	return
}

func main() {
	var wg sync.WaitGroup
	var visited = make(map[string]bool)
	Crawl("http://golang.org/", 4, visited, &wg)
	wg.Wait()
}

type httpResult struct {
	body string
	urls []string
}

func Fetch(url string) (string, []string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", nil, fmt.Errorf("not found: %s", url)
	}
	defer resp.Body.Close()
	body_len := resp.ContentLength
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", nil, fmt.Errorf("not found: %s", url)
	}
	// get response body as string
	resp_body := string(body[:body_len])
	GetLinks(resp_body)
	return "", nil, fmt.Errorf("not found: %s", url)
}

func GetLinks(body string) []string {
	fmt.Println("body", body)
	links := []string{}
	return links
}

/*
func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}
*/
