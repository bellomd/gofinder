// Go finder is a command line application for searching text within
// a file whether locally or remotely. The main application is to search
// for words or sentences from a log file that is located locally in
// your computer or in a remote server.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"knackbrain.com/gofinder/result"

	"knackbrain.com/gofinder/util"
)

const (
	message        = "Please enter your search!"
	searchCreteria = "Invalid search creteria words, filenames or urlnames cannot be empty!"
)

var s = flag.Bool("s", true, "Indicate whether the given input should be searched as single string or list of words")
var c = flag.Bool("c", true, "Indicate whether the search should be case sensitive or not!")

func main() {

	flag.Parse() // parse the flags, search string and files to be searched.
	if len(flag.Args()) <= 0 {
		fmt.Println(message)
		os.Exit(1)
	}

	words := make([]string, 0)
	filenames := make([]string, 0)
	urlnames := make([]string, 0)

	for _, arg := range flag.Args() {
		if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
			urlnames = append(urlnames, arg)
			continue
		}
		if strings.LastIndex(arg, ".") > -1 {
			filenames = append(filenames, arg)
			continue
		}
		words = append(words, arg)
	}

	if len(words) <= 0 || (len(filenames) <= 0 && len(urlnames) <= 0) {
		panic(searchCreteria)
	}

	if len(filenames) > 0 {
		ch := make(chan *result.SearchResult)
		util.SearchAsync(words, filenames, ch, *c)

		var processedCount int
		var exitCount = len(filenames)

		if len(words) > exitCount {
			exitCount = len(words)
		}

		for {

			result, ok := <-ch
			if !ok {
				break
			}

			processedCount++

			if result.Err != nil {
				fmt.Printf("\nError: %s\n", result.Err)
				continue
			}

			fmt.Printf("word: %s\t counts: %d\t filename: %s\n", result.Word, result.Count, result.Filename)
			if processedCount == exitCount {
				break
			}
		}
		close(ch)
	}

	// NOTE:
	// The search urls will be added here.
}
