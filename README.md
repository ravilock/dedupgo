# dedupgo
Simple duplicate file finder written in Go.
Can be used to find identical files under a directory programmatically or via CLI.

# Contributing

Feel free to open issues for bugs, feature requests, improvements and missing tests.
Following Pull Requests will be appreciated!

# Installation

# Usage

## CLI

```bash
$ go install github.com/ravilock/dedupgo@latset
```

Dedupgo can be used with CLI

### Basic Usage

```bash
$ dedupgo find directory-path
```

This will iterate over the files in the passed directory looking for files that have the same content

### Recursive

```bash
$ dedupgo find directory-path -r
```

This will recursively search through all sub-directories of the passed directory looking for files that have the same content,
this option will only stop when there are not sub-directories left

### In-Depth

```bash
$ dedupgo find directory-path -d 4
```

This works similarly to recursive but will define a sub-directory depth limit,
this will only dive to the 4-th directory depth

## Library

Use in your Go project:

```bash
$ go get github.com/ravilock/dedupgo/dedup
```


* Uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies.

```go
import (
	"fmt"

	"github.com/ravilock/dedupgo/dedup"
)

func main() {
	result, err := dedup.Find("/home/raylok/projects/dedupgo-test")
	if err != nil {
		fmt.Println(err)
		return
	}
	for hash, value := range result {
		fmt.Printf("%s are duplicated files under hash %q\n", value, hash)
	}
}
```

