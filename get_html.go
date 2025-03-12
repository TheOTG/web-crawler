package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	res, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't fetch webpage: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return "", fmt.Errorf("couldn't read body: %v", err)
	}
	if res.StatusCode >= 400 {
		return "", fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if !strings.Contains(res.Header.Get("content-type"), "text/html") {
		return "", fmt.Errorf("content type is not html: %v", res.Header.Get("content-type"))
	}

	return string(body), nil
}
