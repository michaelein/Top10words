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
	//TODO validate urls
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
		CountWordsParallel(batch, config)
	}
	printTopWords()
	elapsedTime := time.Now().Sub(startTime)
	removeFile(config)
	fmt.Println("Total time taken:", elapsedTime)
}
