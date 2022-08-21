// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mobiledatabooks/go-fetch/fetcher"
)

//!+

func main() { // Fetch prints the content found at each specified URL.
	fmt.Println("fetcher.FetchWithBuffer: Fetching URLs...") // print message to stdout
	start := time.Now()                                      // start a timer to measure the time it takes to fetch the URLs
	for _, url := range os.Args[1:] {                        // for each URL in the command line arguments
		fetcher.FetchWithBuffer(url) // fetch the URL and print the content
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds()) // print the time elapsed since the start of the timer

	fmt.Println()

	fmt.Println("fetcher.Fetch: Fetching URLs...") // print message to stdout
	start = time.Now()                             // start a timer to measure the time it takes to fetch the URLs
	for _, url := range os.Args[1:] {              // for each URL in the command line arguments
		fetcher.Fetch(url) // fetch the URL and print the content
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds()) // print the time elapsed since the start of the timer

	fmt.Println()

	fmt.Println("fetcher.FetchConcurrent: Fetching URLs...") // print message to stdout
	start = time.Now()                                       // start a timer to measure the time it takes to fetch the URLs
	ch := make(chan string)                                  // make a channel to receive the results of the fetching of the URLs in parallel (concurrent) and return the results to the channel
	for _, url := range os.Args[1:] {                        // for each URL in the command line arguments (concurrent)
		go fetcher.FetchConcurrent(url, ch) // start a goroutine to fetch the URL and return the result to the channel
	}
	for range os.Args[1:] { // for each URL in the command line arguments (concurrent) (wait for the results of the goroutines)
		fmt.Println(<-ch) // receive from channel ch and print the result
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds()) // print the time elapsed since the start of the timer
}

//!-
