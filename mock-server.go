package main

import (
	"net/url"
	"strings"
)

func findBestMatch(urlStr string, responses []Response) int {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return 0
	}

	queryParams := parsedUrl.Query()

	bestMatchIndex := 0
	maxMatches := 0
	for i, response := range responses {
		responseQueryParams, _ := url.ParseQuery(response.Query)
		matches := 0
		for key, values := range responseQueryParams {
			cleanedKey := strings.Trim(key, "?&")
			if queryParams.Get(cleanedKey) == strings.Join(values, "") {
				matches++
			}
		}
		if matches > maxMatches {
			maxMatches = matches
			bestMatchIndex = i
		}
	}

	return bestMatchIndex
}
