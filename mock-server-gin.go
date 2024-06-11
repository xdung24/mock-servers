package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func setupMockServerGin(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
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
