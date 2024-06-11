package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func setupMockServerFiber(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "fiber",
		AppName:       appName,
	})

	// Show server info
	fmt.Println("Running mock server for: ", setting.Name)

	// mock all requests
	for _, request := range setting.Requests {
		// Create a local copy of request for the closure
		localRequest := request
		// Handle the request
		app.All(localRequest.Path, func(c *fiber.Ctx) error {
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

	go app.Listen(fmt.Sprintf("%s:%v", setting.Host, setting.Port))
}
