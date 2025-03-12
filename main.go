package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("no website provided")
	} else if len(args) > 3 {
		log.Fatal("too many arguments provided")
	}

	rawBaseURL := args[0]
	maxConcurrency := 5
	maxPages := 100
	if args[1] != "" {
		num, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatalf("invalid max concurrency: %s", args[1])
		}
		maxConcurrency = num
	}
	if args[2] != "" {
		num, err := strconv.Atoi(args[2])
		if err != nil {
			log.Fatalf("invalid max page: %s", args[2])
		}
		maxPages = num
	}
	fmt.Printf("starting crawler of: %s\n", rawBaseURL)

	cfg, err := newConfig(rawBaseURL, maxConcurrency, maxPages)
	if err != nil {
		log.Fatalf("unable to configure: %v", err)
	}

	cfg.crawlPage(rawBaseURL)
	cfg.wg.Wait()

	cfg.printReport()
}
