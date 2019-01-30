// Go finder is a command line application for searching text within a file whether
// locally or remotely. The main application is to search for words or sentences
// from a log file that is located locally in your computer or in a remote server.
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
	searchCreteria = "Words, filenames or urls cannot be empty!"
)

var s = flag.Bool("s", true, "Should the search be single sentence or list of words")
var c = flag.Bool("c", true, "Should the be case sensitive or not!")

func main() {

	// Parse the flags, search string and files to be searched.
	flag.Parse()
	if len(flag.Args()) <= 0 {
		fmt.Println(message)
		os.Exit(1)
	}

	words, filenames, urls := extract(flag.Args())
	if len(words) <= 0 || (len(filenames) <= 0 && len(urls) <= 0) {
		panic(searchCreteria)
	}

	if len(filenames) > 0 {
		ch := make(chan *result.SearchResult)
		util.SearchAsync(words, filenames, ch, *c)

		var channelBuffer = len(filenames)
		if len(words) > channelBuffer {
			channelBuffer = len(words)
		}
		for {
			result, ok := <-ch
			if !ok {
				break
			}
			channelBuffer--
			if result.Err != nil {
				fmt.Printf("\nError: %s\n", result.Err)
				continue
			}
			fmt.Printf("word: %s\t counts: %d\t filename: %s\n", result.Word, result.Count, result.Filename)
			if channelBuffer == 0 {
				break
			}
		}
		close(ch)
	}
}

func extract(args []string) (words, filenames, urls []string) {

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
	return
}
