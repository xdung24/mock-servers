package main

import (
	"fmt"
	"io/fs"
	"log"

	"github.com/labstack/echo/v4"
)

func setupMockServerEcho(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)

	// Show server info
	fmt.Printf("Serving mock server for: %s on port %v\n", setting.Name, setting.Port)

	// Serve swagger-ui as static files from embedded resources
	static, err := fs.Sub(swaggerUiFolder, "swagger-ui")
	if err != nil {
		panic(err)
	}

	e := echo.New()

	// mock all requests
	for _, request := range setting.Requests {
		// Handle the request
		path := findAndReplaceAll(request.Path)
		e.Add(request.Method, path, func(c echo.Context) error {
			if len(request.Responses) == 0 {
				// write 501 if no response is configured
				c.Response().WriteHeader(501)
				return nil
			}

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
		e.StaticFS("/swagger-ui", fs.FS(static))
	}

	go e.Start(fmt.Sprintf("%s:%v", setting.Host, setting.Port))
}
