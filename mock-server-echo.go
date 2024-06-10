package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func setupMockServerEcho(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)

	// Show server info
	fmt.Println("Running mock server for: ", setting.Name)

	// mock all requests
	e := echo.New()
	for _, request := range setting.Requests {
		// Handle the request
		e.Add(request.Method, request.Path, func(c echo.Context) error {
			// Find the most match response
			matched_index := findBestMatch(c.Request().URL.String(), request.Responses)

			matched_response := request.Responses[matched_index]

			// Write server headers
			for _, header := range setting.Headers {
				c.Response().Header().Set(header.Name, header.Value)
			}

			// write response headers
			for _, header := range matched_response.Headers {
				c.Response().Header().Set(header.Name, header.Value)
			}

			// Return response body
			if matched_response.FilePath != "" {
				res, ok := cacheManager.read(matched_response.FilePath)
				if ok {
					c.Response().WriteHeader(matched_response.Code)
					c.Response().Write(res)
				}
			}
			return nil
		})
	}

	go e.Start(fmt.Sprintf("%s:%v", setting.Host, setting.Port))
}
