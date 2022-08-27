
# Go (Golang) Multi-Module Workspaces: The Easy Way to Build and Run Code in Multiple Modules

If you're writing code in multiple modules at the same time, you can use multi-module workspaces to easily build and run code in those modules.

In Go, a module is a collection of related Go source files located in a single directory. A workspace is a directory containing multiple modules. When you build a Go project, the build system reads the go.mod file in the project's root directory to determine which other modules the project depends on. The build system then downloads and installs any missing dependencies. 

