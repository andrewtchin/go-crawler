package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, visited map[string]bool, wg *sync.WaitGroup) {

	fmt.Printf("Crawl %s\n", url)
	if depth <= 0 {
		fmt.Println("depth 0")
		return
	}
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	url_hash := md5.Sum([]byte(url))
	url_hash_hex := hex.EncodeToString(url_hash[:])

	if visited[url_hash_hex] == true {
		fmt.Printf("skipping %s\n", url)
		return
	}

	fmt.Printf("found: %s %q %s\n", url, body, url_hash_hex)
	visited[url_hash_hex] = true

	for _, u := range urls {
		wg.Add(1)
		go func(the_url string) {
			defer wg.Done()
			Crawl(the_url, depth-1, fetcher, visited, wg)
		}(u)
	}
	return
}

func main() {
	var wg sync.WaitGroup
	var visited = make(map[string]bool)
	Crawl("http://golang.org/", 4, fetcher, visited, &wg)
	wg.Wait()
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
