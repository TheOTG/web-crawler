package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return []string{}, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	htmlReader := strings.NewReader(htmlBody)
	doc, err := html.Parse(htmlReader)
	if err != nil {
		return []string{}, fmt.Errorf("couldn't parse HTML: %v", err)
	}

	urls := []string{}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						return []string{}, fmt.Errorf("couldn't parse href: %v", err)
					}
					resolvedURL := baseURL.ResolveReference(href)
					urls = append(urls, resolvedURL.String())
					// fmt.Printf("found URL: %s\n", resolvedURL.String())
					break
				}
			}
		}
	}
	return urls, nil
}
