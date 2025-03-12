package main

import (
	"cmp"
	"fmt"
	"net/url"
	"slices"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

type Page struct {
	key string
	val int
}

func (cfg *config) getPagesLen() int {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages)
}

func (cfg *config) printReport() {
	fmt.Println("=============================")
	fmt.Printf("REPORT for %s\n", cfg.baseURL.String())
	fmt.Println("=============================")

	var pages []Page
	for k, v := range cfg.pages {
		page := Page{
			key: k,
			val: v,
		}
		pages = append(pages, page)
	}

	slices.SortFunc(pages, sortPageFunc)

	for _, page := range pages {
		fmt.Printf("Found %d internal links to  %s\n", cfg.pages[page.key], page.key)
	}
}

func newConfig(rawBaseURL string, maxConcurrency, maxPages int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	return &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}, nil
}

func sortPageFunc(a, b Page) int {
	cmpVal := cmp.Compare(a.val, b.val)
	if cmpVal == 0 {
		return cmp.Compare(a.key, b.key)
	}
	return cmpVal * -1
}
