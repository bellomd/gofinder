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
func SearchAsync(words []string, files []string, ch chan *result.SearchResult, sensitive bool) {
	for _, filename := range files {
		go search(words, filename, ch, sensitive)
	}
}

// SearchSync searches for the given word or string in a local file.
func SearchSync(words []string, files []string, sensitive bool) (r []*result.SearchResult) {

	var channel chan *result.SearchResult

	for _, filename := range files {
		go search(words, filename, channel, sensitive)
	}

	for ch := range channel {
		r = append(r, ch)
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
}

// searchWord search the given slice of words from the given line of string
// and add the result to result map. The map hold words as keys and occurrences
// as the value of every word.
func searchWord(words []string, line string, results map[string]int, sensitive bool) {
	for _, word := range words {
		if sensitive {
			results[word] += strings.Count(line, word)
		} else {
			results[word] += strings.Count(strings.ToUpper(line), strings.ToUpper(word))
		}
	}
}
