package main

import (
	"flag"
	"fmt"
	"net/http"

	link "example.com/parse"
)

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "url that you want to builld a sitemap for")
	flag.Parse()
 
	fmt.Println(*urlFlag)
	// GET the webpage
	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	links, _ := link.Parse(resp.Body)
}