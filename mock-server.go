package main

import (
	"net/url"
	"regexp"
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

// Replace all variables in the path to other format
// e.g. /api/v1/users/{id}/info/{data} -> /api/v1/users/:id/info/:data
func findAndReplaceAll(input string) string {
	// Define the regex pattern to find all occurrences of {var}
	pattern := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	// Replace all matches with :var
	output := pattern.ReplaceAllString(input, ":$1")
	return output
}
