// Package result hold the result of the search
package result

// SearchResult as a result for the search.
type SearchResult struct {
	Word     string
	Count    int
	Filename string
	Err      error
}
