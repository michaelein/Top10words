package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

var mutex sync.Mutex
var mutex2 sync.Mutex
var numShards = 4
var shardedMaps = make([]map[string]int, numShards)

func init() {
	for i := 0; i < numShards; i++ {
		shardedMaps[i] = make(map[string]int)
	}
}
func extractWordsFromHTML(htmlContent []byte) []string {
	var words []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return words
	}

	// Find all text nodes and extract words
	doc.Find(":not(script)").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		// Split the text into words using spaces
		wordList := strings.Fields(text)

		// Filter out invalid words
		validWords := filterInvalidWords(wordList)

		words = append(words, validWords...)
	})

	return words
}

func filterInvalidWords(words []string) []string {
	var validWords []string

	for _, word := range words {
		// Check if the word is valid
		if isValidWord(word) {
			validWords = append(validWords, word)
		}
	}

	return validWords
}

func isValidWord(word string) bool {
	// Check if the word contains at least 3 alphabetic characters
	if len(word) >= 3 {
		for _, char := range word {
			if !unicode.IsLetter(char) {
				return false
			}
		}
		return true
	}
	return false
}

func fetchEssay(url string, wg *sync.WaitGroup) {
	defer func() {
		if wg != nil {
			wg.Done()
		}
	}()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching essay:", err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading essay:", err)
		return
	}

	// Extract words from HTML content
	words := extractWordsFromHTML(body)

	// Tokenize and count words
	tokenizeAndCount(words)
}
func tokenizeAndCount(words []string) {
	shardIndex := getShardIndex(strings.Join(words, ""))

	// Acquire the mutex before writing to the map
	mutex2.Lock()
	defer mutex2.Unlock()

	for _, word := range words {
		// Check if the word is valid
		if isValidWord(word) {
			shardedMaps[shardIndex][word]++
		}
	}

	// Combine the specific sharded map after the loop finishes
	combineShardedMap(shardedMaps[shardIndex])
}

func combineShardedMap(shardMap map[string]int) {
	mutex.Lock()
	defer mutex.Unlock()

	// Combine the specific sharded map into the global wordCounts map
	for word, count := range shardMap {

		wordCounts[word] += count

	}
}

func getShardIndex(content string) int {
	// Use a simple hash function for sharding
	hash := 0
	for _, char := range content {
		hash += int(char)
	}
	return hash % numShards
}

func CountWordsParallel(essayURLs []string) {
	numGoroutines := 32
	//numCPU := runtime.NumCPU()
	//numGoroutines := numCPU * 2 // Adjust this multiplier based on experimentation
	//fmt.Println("nnnnnnnnnn", len(essayURLs))
	var wg sync.WaitGroup

	// Use a buffered channel to limit the number of concurrent goroutines
	// to numGoroutines
	semaphore := make(chan struct{}, numGoroutines)

	// Launch goroutines for fetching essays and tokenization/counting
	for _, url := range essayURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Acquire semaphore to control the number of concurrent goroutines
			semaphore <- struct{}{}

			fetchEssay(url, nil)
		}(url)
	}

	// Wait for all goroutines to complete
	wg.Wait()
}
