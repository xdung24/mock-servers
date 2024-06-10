package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
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

func setupMockServerGin(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// Show server info
	fmt.Println("Running mock server for: ", setting.Name)

	// mock all requests
	for _, request := range setting.Requests {
		// Handle the request
		r.Handle(request.Method, request.Path, func(c *gin.Context) {
			// Find the most match response
			matched_index := findBestMatch(c.Request.URL.String(), request.Responses)

			matched_response := request.Responses[matched_index]

			// Write server headers
			for _, header := range setting.Headers {
				c.Header(header.Name, header.Value)
			}

			// write response headers
			for _, header := range matched_response.Headers {
				c.Header(header.Name, header.Value)
			}

			// Return response body
			if matched_response.FilePath != "" {
				res, ok := cacheManager.read(matched_response.FilePath)
				if ok {
					c.Data(matched_response.Code, "", res)
				}
			}
		})
	}

	go r.Run(fmt.Sprintf("%s:%v", setting.Host, setting.Port))
}
