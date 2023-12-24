package main

import (
	config2 "awesomeProject/config"
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
)

var wordCounts = make(map[string]int)

type WordCount struct {
	Word  string
	Count int
}

func printTopWords() {
	fmt.Println("Top 10 words:")
	for i := 0; i < 10 && len(wordCounts) > 0; i++ {
		maxWord := findMaxWord(wordCounts)
		fmt.Printf("%s: %d\n", maxWord.Word, maxWord.Count)
		delete(wordCounts, maxWord.Word)
	}
}

func findMaxWord(wordCounts map[string]int) WordCount {
	var maxWord WordCount
	maxCount := 0

	for word, count := range wordCounts {
		if count > maxCount {
			maxCount = count
			maxWord = WordCount{word, count}
		}
	}

	return maxWord
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
