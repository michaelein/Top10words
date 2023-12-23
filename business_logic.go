package main

import (
	config2 "awesomeProject/config"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"unicode"
)

var mutex sync.Mutex
var mutexes []sync.Mutex
var numShards = 4
var shardedMaps = make([]map[string]int, numShards)

func init() {
	for i := 0; i < numShards; i++ {
		shardedMaps[i] = make(map[string]int)
	}

	// Initialize mutexes for each shard
	mutexes = make([]sync.Mutex, numShards)
}

func extractWordsFromHTML(htmlContent []byte) []string {
	var words []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(htmlContent)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return words
	}
	doc.Find(":not(script)").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		wordList := strings.Fields(text)
		validWords := filterInvalidWords(wordList)
		words = append(words, validWords...)
	})
	return words
}

func filterInvalidWords(words []string) []string {
	var validWords []string

	for _, word := range words {
		if isValidWord(word) {
			validWords = append(validWords, word)
		}
	}

	return validWords
}

func isValidWord(word string) bool {
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

	words := extractWordsFromHTML(body)

	tokenizeAndCount(words)
}
func tokenizeAndCount(words []string) {
	shardIndex := getShardIndex(strings.Join(words, ""))

	mutexes[shardIndex].Lock()
	defer mutexes[shardIndex].Unlock()

	for _, word := range words {
		if isValidWord(word) {
			shardedMaps[shardIndex][word]++
		}
	}

	combineShardedMap(shardedMaps[shardIndex])
}

func combineShardedMap(shardMap map[string]int) {
	mutex.Lock()
	defer mutex.Unlock()

	for word, count := range shardMap {
		wordCounts[word] += count
	}
}

func getShardIndex(content string) int {
	hash := 0
	for _, char := range content {
		hash += int(char)
	}
	return hash % numShards
}

func CountWordsParallel(essayURLs []string, config config2.Config) {
	numGoroutines := config.NumGoroutinesMultiplier

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, numGoroutines)

	for _, url := range essayURLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			// Acquire semaphore to control the number of concurrent goroutines
			semaphore <- struct{}{}

			fetchEssay(url, nil)
			<-semaphore
		}(url)
	}

	wg.Wait()
}
