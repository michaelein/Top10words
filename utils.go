package main

import (
	config2 "awesomeProject/config"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

var wordCounts = make(map[string]int)

type WordCount struct {
	Word  string
	Count int
}

func printTopWords() {
	sortedWordCounts := make([]WordCount, 0, len(wordCounts))
	for word, count := range wordCounts {
		sortedWordCounts = append(sortedWordCounts, WordCount{word, count})
	}
	sort.Slice(sortedWordCounts, func(i, j int) bool {
		return sortedWordCounts[i].Count > sortedWordCounts[j].Count
	})

	fmt.Println("Top 10 words:")
	for i := 0; i < 10 && i < len(sortedWordCounts); i++ {
		fmt.Printf("%s: %d\n", sortedWordCounts[i].Word, sortedWordCounts[i].Count)
	}
}

func downloadFile(url, outputPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	return err
}
func downloadAndReadFile(config config2.Config) ([]string, error) {
	if _, err := os.Stat(config.OutputPath); err == nil {
		removeFile(config)
	}
	err := downloadFile(config.FileURL, config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("error downloading file: %v", err)
	}

	file, err := os.Open(config.OutputPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	defer file.Close()

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return urls, nil
}
func removeFile(config config2.Config) {
	err := os.Remove(config.OutputPath)
	if err != nil {
		fmt.Println("Error removing file:", err)
	}
}
