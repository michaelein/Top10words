package main

import (
	config2 "awesomeProject/config"
	"bufio"
	"fmt"
	"os"
	"time"
)

func downloadAndReadFile(config config2.Config) ([]string, error) {
	// Download file
	err := downloadFile(config.FileURL, config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %v", err)
	}

	// Read URLs from file
	file, err := os.Open(config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return urls, nil
}
func main() {
	startTime := time.Now()
	config, err := config2.LoadConfig("./config/config.yaml") // Adjust the file extension based on your configuration file format
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	batchSize := config.BatchSize
	urls, err := downloadAndReadFile(config)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	totalCount := 0
	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}

		batch := urls[i:end]
		totalCount = totalCount + len(batch)
		fmt.Printf("totalCount: %v\n", totalCount)
		CountWordsParallel(batch)
	}
	printTopWords()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	// Remove the file after processing
	err = os.Remove(config.OutputPath)
	if err != nil {
		fmt.Println("Error removing file:", err)
	}
	fmt.Println("Total time taken:", elapsedTime)
}
