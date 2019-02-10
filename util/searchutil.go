package util

import (
	"strings"

	"knackbrain.com/gofinder/result"
)

const (
	// BufferSize default size of the channels
	BufferSize = 100
	// SearchCreteria error displayed when words, files or urls are empty
	SearchCreteria = "Words, filenames or urls cannot be empty!"
)

// Async searches the given words in the given files or urls aysnchronously.
func Async(args []string, sensitive bool, sentence bool) (ch []chan *result.SearchResult) {

	words, filenames, urls := extract(args, sentence)

	if len(filenames) > 0 {
		for _, filename := range filenames {
			channel := make(chan *result.SearchResult, BufferSize)
			ch = append(ch, channel)
			SearchAsync(words, filename, channel, sensitive)
		}
	}

	if len(urls) > 0 {
		for _, url := range urls {
			channel := make(chan *result.SearchResult, BufferSize)
			ch = append(ch, channel)
			URLSearchAsync(words, url, channel, sensitive)
		}
	}
	return
}

// Sync searches the given words synchronously in the file
func Sync(args []string, sensitive bool, sentence bool) (r []*result.SearchResult) {

	words, filenames, urls := extract(args, sentence)
	if len(filenames) > 0 {
		r = append(r, SearchSync(words, filenames, sensitive)...)
	}
	if len(urls) > 0 {
		r = append(r, SearchSync(words, filenames, sensitive)...)
	}
	return
}

// Extract extract words, files and urls from a given slice of strings
func extract(args []string, sentence bool) (words, filenames, urls []string) {

	words = make([]string, 0)
	filenames = make([]string, 0)
	urls = make([]string, 0)

	for _, arg := range args {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			urls = append(urls, arg)
			continue
		}
		if strings.LastIndex(arg, ".") > -1 {
			filenames = append(filenames, arg)
			continue
		}
		words = append(words, arg)
	}

	if len(words) <= 0 || (len(filenames) <= 0 && len(urls) <= 0) {
		panic(SearchCreteria)
	}

	if sentence {
		words = []string{strings.Join(words, " ")}
	}
	return
}
