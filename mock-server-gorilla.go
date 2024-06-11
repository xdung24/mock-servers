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
	fmt.Println("Running mock server for: ", setting.Name)

	// mock all requests
	r := mux.NewRouter()
	for _, request := range setting.Requests {
		// Handle the request
		r.HandleFunc(request.Path, func(w http.ResponseWriter, r *http.Request) {
			// Find the most match response
			matched_index := findBestMatch(r.URL.String(), request.Responses)

			matched_response := request.Responses[matched_index]

			// Write server headers
			for _, header := range setting.Headers {
				w.Header().Set(header.Name, header.Value)
			}

			// write response headers
			for _, header := range matched_response.Headers {
				w.Header().Set(header.Name, header.Value)
			}

			// Return response body
			if matched_response.FilePath != "" {
				res, ok := cacheManager.read(matched_response.FilePath)
				if ok {
					w.WriteHeader(matched_response.Code)
					w.Write(res)
				}
			}
		}).Methods(request.Method)
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
