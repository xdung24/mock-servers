package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func setupMockServerFiber(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "fiber",
		AppName:       appName,
	})

	// Show server info
	fmt.Printf("Serving mock server for: %s on port %v\n", setting.Name, setting.Port)

	// mock all requests
	for _, request := range setting.Requests {
		// Create a local copy of request for the closure
		localRequest := request
		// Handle the request
		app.All(localRequest.Path, func(c *fiber.Ctx) error {
			if len(request.Responses) == 0 {
				// write 501 if no response is configured
				c.Status(http.StatusNotImplemented)
				return nil
			}

			// Find the most match response
			matched_index := findBestMatch(c.Request().URI().String(), localRequest.Responses)

			matched_response := localRequest.Responses[matched_index]

			// Write server headers
			for _, header := range setting.Headers {
				c.Set(header.Name, header.Value)
			}

			// write response headers
			for _, header := range matched_response.Headers {
				c.Set(header.Name, header.Value)
			}

			// Return response body
			if matched_response.FilePath != "" {
				res, ok := cacheManager.read(matched_response.FilePath)
				if ok {
					return c.SendString(string(res))
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
				app.Get("/"+file, func(c *fiber.Ctx) error {
					return c.SendString(string(openapi))
				})
				loaded = true
				break
			}
		}
		if !loaded {
			log.Panicf("OpenAPI file not found (openapi.json/openapi.yml/openapi.yaml) in folder: %s", setting.Folder)
		}
		// serve swagger-ui files
		app.Use("/swagger-ui", filesystem.New(filesystem.Config{
			Root:       http.FS(swaggerUiFolder),
			PathPrefix: "swagger-ui",
		}))
	}

	go app.Listen(fmt.Sprintf("%s:%v", setting.Host, setting.Port))
}
