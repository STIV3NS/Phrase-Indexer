package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	thread, start, end, limit, exclude := getArguments()

	fmt.Printf("%s [ %u - %u ]\n", thread, start, end)
	fmt.Printf("Limit to %u, exclude: %s\n", limit, exclude)
}

func getArguments() (string, uint, uint, uint, string) {
	const sREQUIRED  = ""
	const iREQUIRED = 0

	var thread, exclude string
	var start, end, limit uint

	flag.StringVar(&thread, "thread", sREQUIRED, "[REQUIRED] URL to thread that is meant to be indexed")
	flag.UintVar(&start, "start", 1, "[OPTIONAL] Page number on which to start indexing")
	flag.UintVar(&end, "end", iREQUIRED, "[REQUIRED] Page number on which to end indexing")

	flag.StringVar(&exclude, "exclude", "", "[OPTIONAL] Path to file that contains phrases to exclude from output")
	flag.UintVar(&limit, "limit", math.MaxUint32, "[OPTIONAL] Limit output to top #{value} entries")


	flag.Parse()
	if end == 0 || thread == "" {
		fmt.Fprintf(os.Stderr, "Missing arguments; --help for more information")
		os.Exit(1)
	}

	return thread, start, end, limit, exclude
}