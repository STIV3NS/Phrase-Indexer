package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	threadURL, exclude, start, end, limit := getArguments()



	fmt.Printf("%s [ %u - %u ]\n", threadURL, start, end)
	fmt.Printf("Limit to %u, exclude: %s\n", limit, exclude)
}

func getArguments() (threadURL, exclude string, start, end, limit uint) {
	const sREQUIRED = ""
	const iREQUIRED = 0

	flag.StringVar(&threadURL, "threadURL", sREQUIRED,
		"[REQUIRED] URL to threadURL that is meant to be indexed")
	flag.UintVar(&start, "start", 1,
		"[OPTIONAL] Page number on which to start indexing")
	flag.UintVar(&end, "end", iREQUIRED,
		"[REQUIRED] Page number on which to end indexing")

	flag.StringVar(&exclude, "exclude", "",
		"[OPTIONAL] Path to file that contains phrases to exclude from output")
	flag.UintVar(&limit, "limit", math.MaxUint32,
		"[OPTIONAL] Limit output to top #{value} entries")

	flag.Parse()

	if end == 0 || threadURL == "" {
		fmt.Fprintf(os.Stderr, "Missing arguments; --help for more information\n")
		os.Exit(1)
	}

	return
}