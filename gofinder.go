// Go finder is a command line application for searching text within a file whether
// locally or remotely. The main application is to search for words or sentences
// from a log file that is located locally in your computer or in a remote server.
package main

import (
	"flag"
	"fmt"
	"os"

	"knackbrain.com/gofinder/result"

	"knackbrain.com/gofinder/util"
)

const (
	message = "Please enter your search!"
)

var s = flag.Bool("s", true, "Single sentence or list of words")
var c = flag.Bool("c", true, "Case sensitive or insensitive!")

func main() {

	// Parse the flags, search string and files to be searched.
	flag.Parse()
	if len(flag.Args()) <= 0 {
		fmt.Println(message)
		os.Exit(1)
	}
	// Searches that return channels to read the result.
	processChannels(util.Async(flag.Args(), *c, *s))

	// Searches that return slices of results.
	//processResults(util.Sync(flag.Args(), *c, *s))
}

func processResults(r []*result.SearchResult) {
	for _, result := range r {
		if result.Err != nil {
			fmt.Printf("Error => %s\n", result.Err)
		} else {
			fmt.Printf("word: %s\t counts: %d\t filename: %s\n", result.Word, result.Count, result.Filename)
		}
	}
}

func processChannels(chs []chan *result.SearchResult) {
	for {
		closedChannel := 0
		for _, channel := range chs {
			result, ok := <-channel
			if !ok {
				closedChannel++
				continue
			}
			if result.Err != nil {
				fmt.Printf("Error => %s\n", result.Err)
			} else {
				fmt.Printf("word: %s\t counts: %d\t filename: %s\n", result.Word, result.Count, result.Filename)
			}
		}
		if closedChannel == len(chs) {
			break
		}
	}
}
