package main

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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

	_, urls, err := Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	visited[url_hash_hex] = true
	fmt.Printf("links found on page: %s %s\n", url, url_hash_hex)

	for e := urls.Front(); e != nil; e = e.Next() {
		u := e.Value.(string)
		fmt.Printf("-- %s\n", u)
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

func Fetch(url string) (string, *list.List, error) {
	fmt.Printf("Fetch %s\n", url)
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
	links := GetLinks(url, resp_body, body_len)
	return "", links, nil
}

func GetLinks(url string, body string, body_len int64) *list.List {
	links := list.New()
	if strings.HasSuffix(url, "/") {
		// strip trailing slash
		url = url[0 : len(url)-1]
	}
	re := regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)
	r2 := re.FindAllStringSubmatch(body, -1)

	for _, m := range r2 {
		// fmt.Printf("%s\n", m[1])
		if strings.HasPrefix(m[1], "//") {
			url_value := "http:"
			url_value += m[1]
			// fmt.Println("->", url_value)
			links.PushBack(url_value)
		} else if strings.HasPrefix(m[1], "/") {
			url_value := url
			url_value += m[1]
			// fmt.Println("->", url_value)
			links.PushBack(url_value)
		} else if strings.HasPrefix(m[1], "#") {
			// fmt.Println("skip", m[1])
		} else {
			// fmt.Println("ok", m[1])
			links.PushBack(m[1])

		}
	}

	return links
}
