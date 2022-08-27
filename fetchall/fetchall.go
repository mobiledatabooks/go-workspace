// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

// MIT License

// Copyright (c) 2022 Mobile Data Books, LLC

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
