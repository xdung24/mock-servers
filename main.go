package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Get the configuration
	config := getEnvConfig()
	if err := config.Validate(); err != nil {
		fmt.Println("error validating the configuration: ", err)
		return
	} else {
		fmt.Println("data folder: ", config.DataFolder)
		fmt.Println("use fsnotify: ", config.UseFsNotify)
		fmt.Println("use polling: ", config.UsePolling)
		fmt.Println("polling time: ", config.PollingTime)
		fmt.Println("web engine: ", config.WebEngine)
	}

	// Create a channel to capture the interrupt signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Create a cache manager
	cacheManager := newCacheManager(context.TODO())

	// Setup mock servers
	appsToMock := listSubfolders(config.DataFolder)
	for _, appToMock := range appsToMock {
		if config.WebEngine == "gin" {
			setupMockServerGin(appToMock, cacheManager)
		} else if config.WebEngine == "gorilla" {
			setupMockServerGorilla(appToMock, cacheManager)
		} else if config.WebEngine == "echo" {
			setupMockServerEcho(appToMock, cacheManager)
		} else if config.WebEngine == "fiber" {
			setupMockServerFiber(appToMock, cacheManager)
		} else {
			fmt.Println("error: web engine not supported")
			return
		}
	}

	// Watch for changes in the data folder
	if config.UseFsNotify {
		go watchDirectory(config.DataFolder)
	} else if config.UsePolling {
		go pollingDirectory(config.DataFolder, time.Duration(config.PollingTime))
	}

	// pprof
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Wait for the interrupt signal
	<-interrupt
	fmt.Println("terminating the application.")
}

func listSubfolders(rootDir string) []string {
	var subfolders []string
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		fmt.Println(err)
		return subfolders
	}

	// Append subfolder to result
	for _, entry := range entries {
		if entry.IsDir() {
			subfolders = append(subfolders, entry.Name())
		}
	}
	return subfolders
}
