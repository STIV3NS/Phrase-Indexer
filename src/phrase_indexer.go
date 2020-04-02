package main

import (
	"flag"
	"fmt"
	"math"
	"os"
)

func main() {
	thread, exclude, start, end, limit := getArguments()

	fmt.Printf("%s [ %u - %u ]\n", thread, start, end)
	fmt.Printf("Limit to %u, exclude: %s\n", limit, exclude)
}

func getArguments() (thread, exclude string, start, end, limit uint) {
	const sREQUIRED  = ""
	const iREQUIRED = 0

	flag.StringVar(&thread, "thread", sREQUIRED, "[REQUIRED] URL to thread that is meant to be indexed")
	flag.UintVar(&start, "start", 1, "[OPTIONAL] Page number on which to start indexing")
	flag.UintVar(&end, "end", iREQUIRED, "[REQUIRED] Page number on which to end indexing")

	flag.StringVar(&exclude, "exclude", "", "[OPTIONAL] Path to file that contains phrases to exclude from output")
	flag.UintVar(&limit, "limit", math.MaxUint32, "[OPTIONAL] Limit output to top #{value} entries")


	flag.Parse()
	if end == 0 || thread == "" {
		fmt.Fprintf(os.Stderr, "Missing arguments; --help for more information\n")
		os.Exit(1)
	}

	return
}