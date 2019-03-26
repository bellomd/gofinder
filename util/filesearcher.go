// Package util contain utilities for reading file either locally
// or remotely.
//
// The search can be for a single word or for a complete sentence.
package util

import (
	"bufio"
	"os"
	"strings"

	"knackbrain.com/gofinder/result"
)

// SearchAsync searches for the given word or string in a local file.
func SearchAsync(words []string, filename string, ch chan *result.SearchResult, sensitive bool) {
	go search(words, filename, ch, sensitive)
}

// SearchSync searches for the given word or string in a local file.
func SearchSync(words []string, files []string, sensitive bool) (r []*result.SearchResult) {

	var channels []chan *result.SearchResult
	for _, filename := range files {
		channel := make(chan *result.SearchResult, BufferSize)
		channels = append(channels, channel)
		go search(words, filename, channel, sensitive)
	}

	for {
		closedChannel := 0
		for _, channel := range channels {
			result, ok := <-channel
			if !ok {
				closedChannel++
				continue
			}
			r = append(r, result)
		}
		if closedChannel == len(channels) {
			break
		}
	}
	return
}

func search(words []string, filename string, ch chan *result.SearchResult, sensitive bool) {

	file, err := os.Open(filename)
	if err != nil {
		ch <- &result.SearchResult{
			Filename: filename,
			Err:      err,
		}
		close(ch)
		return
	}
	defer file.Close()

	results := make(map[string]int)

	scanner := bufio.NewScanner(file) // create a scanner to read the content of the file.
	for scanner.Scan() {
		searchWord(words, scanner.Text(), results, sensitive) // read content line by line.
	}

	// Add the result of the search to channel.
	for key, value := range results {
		ch <- &result.SearchResult{
			Word:     key,
			Count:    value,
			Filename: filename,
			Err:      nil,
		}
	}
	close(ch)
}

// searchWord search the given slice of words from the given line of string
// and add the result to result map. The map hold words as keys and occurrences
// as the value of every word.
func searchWord(words []string, line string, results map[string]int, sensitive bool) {
	for _, word := range words {
		if sensitive {
			results[word] += strings.Count(line, word)
		} else {
			results[word] += strings.Count(strings.ToLower(line), strings.ToLower(word))
		}
	}
}
