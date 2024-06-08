package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

// watchDirectory sets up a watcher on the specified directory path.
func watchDirectory(path string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// Start a goroutine to handle events and errors.
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	// Walk through the directory tree and add each directory to the watcher.
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Wait for the done signal. In a real application, you might use a more robust mechanism
	// for managing goroutine lifecycle.
	<-done
	return nil
}

// calculateFolderChecksum calculates the MD5 checksum for all files in a folder and its subfolders.
func calculateFolderChecksum(folderPath string) (string, error) {
	hash := md5.New() // Initialize the hash function

	// Walk the directory tree
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // Skip directories
		}

		// Read file content
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// Update the hash with the file content
		_, err = hash.Write(data)
		return err
	})

	if err != nil {
		return "", err
	}

	// Convert the hash to a hexadecimal string
	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, nil
}

func pollingDirectory(folderPath string, interval time.Duration) {
	pollingInterval := interval * time.Second // Polling interval

	var lastChecksum string = ""
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()

	for range ticker.C {
		checksum, err := calculateFolderChecksum(folderPath)
		if err != nil {
			log.Printf("Error calculating folder checksum: %v", err)
			continue
		}
		log.Printf("Folder checksum: %s", checksum)
		var changed bool = false
		if checksum != lastChecksum {
			if lastChecksum != "" {
				changed = true
			}
			lastChecksum = checksum // Update the last known checksum
		}

		if changed {
			// Perform your action on change here
			log.Println("Folder has changed")
		}
	}
}
