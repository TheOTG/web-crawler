package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string) {

	cfg.wg.Add(1)

	go func() {
		cfg.concurrencyControl <- struct{}{}
		defer func() {
			<-cfg.concurrencyControl
			cfg.wg.Done()
		}()

		if cfg.getPagesLen() >= cfg.maxPages {
			return
		}

		currentURL, err := url.Parse(rawCurrentURL)
		if err != nil {
			fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
			return
		}
		if currentURL.Hostname() != cfg.baseURL.Hostname() {
			return
		}

		normalized, err := normalizeURL(rawCurrentURL)
		if err != nil {
			fmt.Printf("unable to normalize current url: %v\n current URL: %v", err, rawCurrentURL)
			return
		}

		isFirst := cfg.addPageVisit(normalized)
		if !isFirst {
			return
		}

		fmt.Printf("crawling %s\n", rawCurrentURL)
		html, err := getHTML(rawCurrentURL)
		if err != nil {
			fmt.Printf("unable to get html: %v", err)
			return
		}

		nextURLs, err := getURLsFromHTML(html, cfg.baseURL.String())
		if err != nil {
			fmt.Printf("unable to get urls: %v", err)
			return
		}

		for _, nextURL := range nextURLs {
			// time.Sleep(50 * time.Millisecond)
			cfg.crawlPage(nextURL)
		}
	}()
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, ok := cfg.pages[normalizedURL]
	if ok {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}
