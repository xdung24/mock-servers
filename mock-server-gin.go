package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupMockServerGin(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetTrustedProxies(nil)

	// Show server info
	fmt.Printf("Serving mock server for: %s on port %v\n", setting.Name, setting.Port)

	// Server swagger-ui as static files from embedded resources
	if setting.SwaggerEnabled {
		static, err := fs.Sub(swaggerUiFolder, "swagger-ui")
		if err != nil {
			panic(err)
		}

		// Serve openapi file if file exists
		openApiFiles := []string{"openapi.json", "openapi.yml", "openapi.yaml"}
		loaded := false
		for _, file := range openApiFiles {
			filePath := fmt.Sprintf("data/%s/%s", setting.Folder, file)
			if openapi, ok := cacheManager.read(filePath); ok {
				r.GET("/"+file, func(c *gin.Context) {
					c.Data(200, "application/json", openapi)
				})
				loaded = true
				break
			}
		}
		if !loaded {
			log.Panicf("OpenAPI file not found (openapi.json/openapi.yml/openapi.yaml) in folder: %s", setting.Folder)
		}

		// serve swagger-ui files
		r.StaticFS("/swagger-ui", http.FS(static))
	}

	// mock all requests
	for _, request := range setting.Requests {
		// Handle the request
		r.Handle(request.Method, request.Path, func(c *gin.Context) {
			if len(request.Responses) == 0 {
				// write 501 if no response is configured
				c.Status(http.StatusNotImplemented)
				return
			}

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
