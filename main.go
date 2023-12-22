package main

import (
	config2 "awesomeProject/config"
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	// Record the start time
	startTime := time.Now()
	// Load configuration from file
	config, err := config2.LoadConfig("./config/config.yaml") // Adjust the file extension based on your configuration file format
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	batchSize := config.BatchSize
	err = downloadFile(config.FileURL, config.OutputPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	filePath := "urls.txt"
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
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
