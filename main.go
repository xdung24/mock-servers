package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create a channel to capture the interrupt signal (Ctrl+C)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Create a cache manager
	ctx := context.TODO()
	cacheManager := newCacheManager(ctx)

	// Setup mock servers
	appsToMock := listSubfolders("./data")
	for _, appToMock := range appsToMock {
		setupMockServer(appToMock, cacheManager)
	}

	// Wait for the interrupt signal
	fmt.Println("App is running. Press Ctrl+C to terminate.")
	<-interrupt
	fmt.Println("Terminating the application.")
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
