package main

import (
	config2 "awesomeProject/config"
	"fmt"
	"time"
)

func main() {
	startTime := time.Now()
	config, err := config2.LoadConfig("./config/config.yaml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}

	if err := processUrls(config); err != nil {
		fmt.Println("Error:", err)
		return
	}

	elapsedTime := time.Now().Sub(startTime)
	fmt.Println("Total time taken:", elapsedTime)
}

func processUrls(config config2.Config) error {
	batchSize := config.BatchSize
	urls, err := downloadAndReadFile(config)
	if err != nil {
		return err
	}

	processBatches(urls, config, batchSize)
	printTopWords()
	removeFile(config)

	return nil
}

func processBatches(urls []string, config config2.Config, batchSize int) {
	totalCount := 0
	for i := 0; i < len(urls); i += batchSize {
		end := i + batchSize
		if end > len(urls) {
			end = len(urls)
		}
		batch := urls[i:end]
		totalCount += len(batch)
		fmt.Printf("totalCount: %v\n", totalCount)
		CountWordsParallel(batch, config)
	}
}
