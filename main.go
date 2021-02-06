package main

import (
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func main() {
	for _, url := range os.Args[1:] {
		links, err := getLinks(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse: %v/n", err)
		}

		for _, link := range links {
			fmt.Println(link)
		}
	}
}

func getLinks(url string) ([]string, error) {
	response, err := makeRequest(url)

	if err != nil {
		return nil, fmt.Errorf("Something bad happen")
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)

	if err != nil {
		return nil, fmt.Errorf("Error parse html")
	}

	return validateLink(nil, doc), nil
}

func makeRequest(url string) (*http.Response, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status not OK")
	}

	return response, nil
}

func validateLink(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = validateLink(links, c)
	}

	return links
}
