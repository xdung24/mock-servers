package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func setupMockServerGorilla(appName string, cacheManager *CacheManager) {
	setting := parseSetting(appName)
	setting.loadResources(cacheManager)

	// Show server info
	fmt.Printf("Serving mock server for: %s on port %v\n", setting.Name, setting.Port)

	r := mux.NewRouter()

	// mock all requests
	for _, request := range setting.Requests {
		// Handle the request
		r.HandleFunc(request.Path, func(w http.ResponseWriter, r *http.Request) {
			if len(request.Responses) == 0 {
				// write 501 if no response is configured
				w.WriteHeader(http.StatusNotImplemented)
				return
			}

			// Find the most match response
			matched_index := findBestMatch(r.URL.String(), request.Responses)

			matched_response := request.Responses[matched_index]

			// Write server headers
			if setting.Headers != nil {
				for _, header := range *setting.Headers {
					w.Header().Set(header.Name, header.Value)
				}
			}

			// write response headers
			if matched_response.Headers != nil {
				for _, header := range *matched_response.Headers {
					w.Header().Set(header.Name, header.Value)
				}
			}

			// write status code
			w.WriteHeader(matched_response.Code)

			// Return response body
			if matched_response.FilePath != nil && *matched_response.FilePath != "" {
				res, ok := cacheManager.read(*matched_response.FilePath)
				if ok {
					w.Write(res)
				}
			}
		}).Methods(request.Method)
	}

	// Server swagger-ui as static files from embedded resources
	if setting.SwaggerEnabled {
		// serve openapi file if file exists
		openApiFiles := []string{"openapi.json", "openapi.yml", "openapi.yaml"}
		loaded := false
		for _, file := range openApiFiles {
			filePath := fmt.Sprintf("data/%s/%s", setting.Folder, file)
			if openapi, ok := cacheManager.read(filePath); ok {
				r.HandleFunc("/"+file, func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					w.Write(openapi)
				})
				loaded = true
				break
			}

		}
		if !loaded {
			log.Panicf("OpenAPI file not found (openapi.json/openapi.yml/openapi.yaml) in folder: %s", setting.Folder)
		}

		// serve swagger-ui
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.FS(swaggerUiFolder))))
	}

	srv := &http.Server{
		Addr: fmt.Sprintf("%s:%v", setting.Host, setting.Port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
}
