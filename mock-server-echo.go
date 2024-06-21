package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

func setupMockServerEcho(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)

	// Show server info
	fmt.Printf("Serving mock server for: %s on port %v\n", setting.Name, setting.Port)

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

	// Server swagger-ui as static files from embedded resources
	if setting.SwaggerEnabled {
		// serve openapi file if file exists
		openApiFiles := []string{"openapi.json", "openapi.yml", "openapi.yaml"}
		loaded := false
		for _, file := range openApiFiles {
			filePath := fmt.Sprintf("data/%s/%s", setting.Folder, file)
			if openapi, ok := cacheManager.read(filePath); ok {
				e.GET("/"+file, func(c echo.Context) error {
					c.Response().Header().Set("Content-Type", "application/json")
					c.Response().WriteHeader(200)
					c.Response().Write(openapi)
					return nil
				})
				loaded = true
				break
			}
		}

		if !loaded {
			log.Panicf("OpenAPI file not found (openapi.json/openapi.yml/openapi.yaml) in folder: %s", setting.Folder)
		}

		// serve swagger-ui
		e.Static("/swagger-ui", "swagger-ui")
	}

	go e.Start(fmt.Sprintf("%s:%v", setting.Host, setting.Port))
}
