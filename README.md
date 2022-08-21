
If you're writing code in multiple modules at the same time, you can use multi-module workspaces to easily build and run code in those modules.

In Go, a module is a collection of related Go source files located in a single directory. A workspace is a directory containing multiple modules. When you build a Go project, the build system reads the go.mod file in the project's root directory to determine which other modules the project depends on. The build system then downloads and installs any missing dependencies. 

This article is a part of the series of books by [Constantine Vassil](https://www.amazon.com/Constantine-Vassil/e/B09Z9S1Y77/ref=aufs_dp_fta_dsk).


In these hands on projects, you'll learn how to use Golang to develop applications quickly and effectively, both locally with multi-module workspaces and on Google Cloud. You'll get hands-on experience with the language, learning how to write code, debug applications, and deploy to the cloud. These quests will help you get started with using Golang on Google Cloud, and you'll be able to apply what you've learned to your own projects.

## Create the workspace folder:
```go
    mkdir go-workspace
    cd go-workspace
```

## Create the main module. Initialize the module:
```go
mkdir fetchall
cd fetchall
go mod init mobiledatabooks.com/fetchall
```

## Create go.mod. Add a dependency on the github.com/mobiledatabooks/go-fetch/fetcher module by using:
```go
go mod tidy
```

## Create fetchall.go in the fetchall directory with the following contents:

```go
touch fetchall.go
```

```go
package main

import (
	"os"

	"github.com/mobiledatabooks/go-fetch/fetcher"
)

func main() {
	for _, url := range os.Args[1:] {
		fetcher.Fetch(url)
	}
}
```

```go
go run main.go https://golang.org http://gopl.io https://godoc.org
```

- The program starts with a channel ch of type string to receive the results of the fetching of the URLs in parallel (concurrent).
- Then, for each URL in the command line arguments (concurrent), it starts a goroutine to fetch the URL and return the result to the channel.
- Finally, for each URL in the command line arguments (concurrent) (wait for the results of the goroutines), it receives from channel ch and prints the result.

```go
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetcher.FetchConcurrent(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

```

```go
go build -v *.go
```
## Now, run the main program:

```go
./fetchall https://golang.org http://gopl.io https://godoc.org

0.68s     4154  http://gopl.io
0.78s    17461  https://godoc.org
1.29s    59769  https://golang.org
1.29s elapsed


```

## Create the workspace

In this step, we’ll create a go.work file to specify a workspace with the module.

Change to the workspace directory
```go
cd go-workspace
```
In the workspace directory, run:
```go
go work init ./fetchall
```
The go work init command tells go to create a go.work file for a workspace containing the modules in the ./fetchall directory.

Created go.work
```go
go 1.19

use ./fetchall
```

The go.work file has similar syntax to go.mod.

The go directive tells Go which version of Go the file should be 
interpreted with. It’s similar to the go directive in the go.mod file.

The use directive tells Go that the module in the hello directory 
should be main modules when doing a build.

So in any subdirectory of workspace the module will be active.

## Run the program in the workspace directory

In the workspace directory, run:

```go
go run mobiledatabooks.com/fetchall  https://golang.org http://gopl.io https://godoc.org

0.68s     4154  http://gopl.io
0.78s    17461  https://godoc.org
1.29s    59769  https://golang.org
1.29s elapsed
```

The Go command includes all the modules in the workspace as main modules. 
This allows us to refer to a package in the module, even outside the module. 
Running the go run command outside the module or the workspace would result 
in an error because the go command wouldn’t know which modules to use.

Next, we’ll add a local copy of the go-fetch module to the workspace. 
We’ll then add a new function to the fetcher package that we can use instead of Fetch.

## Download and modify the github.com/mobiledatabooks/go-fetch module
 
In this step, we’ll download a copy of the Git repo containing the 
github.com/mobiledatabooks/go-fetch module, add it to the workspace, and then add a new function to it that we will use from the main program.


### Clone the repository

```go
cd go-workspace
```
From the workspace directory, run the git command to clone the repository:

```go
cd go-workspace
git clone https://github.com/mobiledatabooks/go-fetch.git
```
It creates go-fetch directory
```go

➜  go-workspace git:(main) ✗ 
.
├── LICENSE
├── README.md
├── fetchall
│   ├── fetchall.go
│   ├── go.mod
│   └── go.sum
├── go-fetch
│   ├── LICENSE
│   ├── README.md
│   ├── fetcher
│   │   ├── fetcher.go
│   │   └── go.mod
│   ├── go.mod
│   └── main.go
├── go.work

```

Add the module to the workspace
```go
go work use ./go-fetch
```

The go work use command adds a new module to the go.work file. 
It will now look like this:

go.work
```go
use (
	./fetchall
	./go-fetch
)
```

The module now includes both the mobiledatabooks.com/fetchall module and the github.com/mobiledatabooks/go-fetch/fetcher module.

## Working with local copy

This will allow us to use the new code we will write in our copy of 
the stringutil module instead of the version of the module in the 
module cache that we downloaded with the go get command.

### Add the new function.

We’ll add a new function to fetch concurrently to the github.com/mobiledatabooks/go-fetch/fetcher package.

Create a new file named fetchConcurrent.go in the go-workspace/go-fetch/fetcher directory containing the following contents:

```go
// Fetch prints the content found at each specified URL.
package fetcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// !+
func FetchConcurrent(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err) // send to channel ch
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close() // don't leak resources
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}

//!-

```

### Modify the main program to use the function.

Modify the contents of go-workspace/fetchall/fetchall.go to contain the following contents:

```go
//!+

// Fetchall fetches URLs in parallel and reports their times and sizes.
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mobiledatabooks/go-fetch/fetcher"
)

//!-

func main() {
	fmt.Println("fetcher.Fetch: Fetching URLs...")
	start := time.Now()
	for _, url := range os.Args[1:] {
		fetcher.Fetch(url)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	fmt.Println()
	fmt.Println("fetcher.FetchConcurrent: Fetching URLs...")
	start = time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetcher.FetchConcurrent(url, ch) // start a goroutine
	}
	for range os.Args[1:] {
		fmt.Println(<-ch) // receive from channel ch
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

//!-

```

Run the code in the workspace

From the workspace directory, run

```go
go run main.go https://golang.org http://gopl.io https://godoc.org
```

The Go command finds the mobiledatabooks.com/fetchall module specified in the command line in the fetchall directory specified by the go.work file, and similarly resolves the github.com/mobiledatabooks/go-fetch/fetcher import using the go.work file.

go.work can be used instead of adding replace directives to work across multiple modules.


### Make a change in one module and use it in another.

Since the two modules are in the same workspace it’s easy to make a change in one module and use it in another.

Future step

Now, to properly release these modules we’d need to make a release of 
the github.com/mobiledatabooks/go-fetch/fetcher module, for example at v0.1.0. 
This is usually done by tagging a commit on the module’s version control repository. See the module release workflow documentation for more details. 
Once the release is done, we can increase the requirement on the 
github.com/mobiledatabooks/go-fetch/fetcher module in fetchall/go.mod:
by deleting everhything after "go 1.19" and performing go mod tidy

```go
go mod tidy
```

# Module release and versioning workflow

https://go.dev/doc/modules/release-workflow

When you develop modules for use by other developers, you can follow a workflow that helps ensure a reliable, consistent experience for developers using the module. 
This topic describes the high-level steps in that workflow.

See also

If you’re merely wanting to use external packages in your code, be sure to see 
Managing dependencies.
With each new version, you signal the changes to your module with its version number. 
For more, see Module version numbering.

Common workflow steps

The following sequence illustrates release and versioning workflow steps for an example 
new module. For more about each step, see the sections in this topic.

1. Begin a module and organize its sources to make it easier for developers to use and for you to maintain.

If you’re brand new to developing modules, check out Tutorial: Create a Go module.

In Go’s decentralized module publishing system, how you organize your code matters. 
For more, see Managing module source.

2. Set up to write local client code that calls functions in the unpublished module.

Before you publish a module, it’s unavailable for the typical dependency management workflow using commands such as go get. A good way to test your module code at this stage is to try it while it is in a directory local to your calling code.

See Coding against an unpublished module for more about local development.

3. When the module’s code is ready for other developers to try it out, 
begin publishing v0 pre-releases such as alphas and betas. See Publishing pre-release versions for more.

Note. Example from:
```
The Go Programming Language
Alan A. A. Donovan
Google Inc.
Brian W. Kernighan
Princeton University
Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan

SECTION 1.6. FETCHING URLS CONCURRENTLY
```
