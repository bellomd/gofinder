// Package util contain utilities for reading file either locally
// or remotely.
//
// The search can be for a single word or for a complete sentence.
package util

import (
	"knackbrain.com/gofinder/result"
)

// URLSearchAsync searches for the give words in the given urls asynchronously and add the result to the given
// channel, the search can be case sensitive or case insensitive.
func URLSearchAsync(words []string, url string, ch chan *result.SearchResult, sensitive bool) {
	panic("Not implemented!")
}

// URLSearchSync searches for the give words in the given urls synchronously and return the result
// to the client.
func URLSearchSync(words []string, url string, sensitive bool) (r []*result.SearchResult) {
	panic("Not implemented!")
}
