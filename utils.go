package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
)

var wordCounts = make(map[string]int)

// WordCount represents a word and its count
type WordCount struct {
	Word  string
	Count int
}

func printTopWords() {
	// Sort the word counts by frequency in descending order
	sortedWordCounts := make([]WordCount, 0, len(wordCounts))
	for word, count := range wordCounts {
		sortedWordCounts = append(sortedWordCounts, WordCount{word, count})
	}
	sort.Slice(sortedWordCounts, func(i, j int) bool {
		return sortedWordCounts[i].Count > sortedWordCounts[j].Count
	})

	// Print the top 10 words
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
