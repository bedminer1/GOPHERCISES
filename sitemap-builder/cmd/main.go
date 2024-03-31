package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	link "example.com/parse"
)
const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Urls []loc `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "url that you want to builld a sitemap for")
	maxDepth := flag.Int("depth", 5, "number of layers of BFS search")
	 
	flag.Parse()
 
	pages := bfs(*urlFlag, *maxDepth)
	toXml := urlSet {
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	// set of visited pages
	seen := make(map[string]struct{})
	var q map[string]struct{}
	// next layer
	nq := map[string]struct{} {
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{}) 
		// optimize breaking at empty queues
		if len(q) == 0 {
			break
		}
		for url := range q {
			// already visited
			if _, ok := seen[url]; ok {
				continue
			}

			// set to visited
			seen[url] = struct{}{}

			// add children to next queue
			for _, link := range getWebPage(url) {
				// check for child in visited set
				if _, ok := seen[link]; ok {
					continue
				}
				nq[link] = struct{}{}
			} 
		}
	}

	// convert seen hashSet into a slice
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
 }

func getWebPage(urlStr string) []string {
	// GET the webpage
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	
	// get the base url to prefix it to internal links
	reqUrl := resp.Request.URL
	baseUrl := &url.URL {
		Scheme: reqUrl.Scheme,
		Host: reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(getHrefs(resp.Body, base), withPrefix(base))
}

func getHrefs(r io.Reader, base string) []string {
	// parse the links
	links, _ := link.Parse(r)
	var hrefs []string

	// filter out stuff like fragments, fixup internal links so
	// get request can be sent
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base + l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	return hrefs 
}

// filter function to only add links that has allowed prefix
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}