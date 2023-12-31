package dedup

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path"

	"github.com/ravilock/dedupgo/internal"
)

type FindMode int

const (
	ModeNormal = iota
	ModeRecursive
	ModeDepth
)

// FindOptions represent the options that can be passed to the Find function, it can be freely instatiated and passed by callers
type FindOptions struct {
	// Mode represents the mode that the Find function will execute, it has 3 different possiblities:
	//  - Normal: Search for duplicated files in current directory only, this option will not dive into any sub-directories
	//  - Recursive: Recrusively search through all sub-directories, this option will only stop when there are not sub-directories left
	//  - Depth: Define a sub-directories depth limit, this will only dive to the n-th directory depth (n being the limit passed)
	// Default value is Normal
	Mode FindMode

	// DepthLimit represents the limit that the Find function will use when "Depth" Mode is passed
	// If DepthLimit is passed without "Depth" Mode, this options will be ignored
	DepthLimit int
}

// DuplicatedFiles groups a slice of duplicated file paths by content-hash
// Files that are in the same slice (under the same map key) have the same content
type DuplicatedFiles map[string][]string

// Find will search inside directoryPath looking for files that are duplicated
func Find(directoryPath string, opts ...*FindOptions) (DuplicatedFiles, error) {
	options := mergeFindOptions(opts...)
	result, err := find(directoryPath, options)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func mergeFindOptions(opts ...*FindOptions) *FindOptions {
	options := &FindOptions{}
	for _, opt := range opts {
		if opt == nil {
			return options
		}
		return opt
	}
	return options
}

type visitedDirectory struct {
	path  string
	depth int
}

func find(directoryPath string, opts *FindOptions) (DuplicatedFiles, error) {
	directoryQueue := &internal.Queue[visitedDirectory]{}

	directoryQueue.Enqueue(visitedDirectory{directoryPath, 0})

	filePaths := make([]string, 0)

	for directoryQueue.Length() != 0 {
		currentDir, ok := directoryQueue.Dequeue()
		if !ok || (opts.DepthLimit != 0 && currentDir.depth == opts.DepthLimit) {
			break
		}

		dirEntries, err := os.ReadDir(currentDir.path)
		if err != nil {
			return nil, err
		}

		for _, entry := range dirEntries {
			entryPath := path.Join(currentDir.path, entry.Name())
			if entry.IsDir() {
				directoryQueue.Enqueue(visitedDirectory{entryPath, currentDir.depth + 1})
			} else {
				filePaths = append(filePaths, entryPath)
			}
		}

		if opts.Mode == ModeNormal {
			break
		}
	}

	filePathsByHash, err := groupFilesByContentHash(filePaths, directoryPath)
	if err != nil {
		return nil, err
	}
	removeNonDuplicatedPaths(filePathsByHash)
	return filePathsByHash, nil
}

func groupFilesByContentHash(files []string, dirPath string) (DuplicatedFiles, error) {
	contentHashToFileName := make(DuplicatedFiles)
	for _, file := range files {
		hash2, err := hashFile(file, dirPath)
		if err != nil {
			return nil, err
		}
		hashString := hex.EncodeToString(hash2)
		_, ok := contentHashToFileName[hashString]
		if !ok {
			contentHashToFileName[hashString] = []string{file}
		} else {
			contentHashToFileName[hashString] = append(contentHashToFileName[hashString], file)
		}
	}
	return contentHashToFileName, nil
}

func hashFile(filePath string, dirPath string) ([]byte, error) {
	buffer := make([]byte, 1024)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	hasher := sha256.New()
	for {
		n, err := f.Read(buffer)
		if err != nil {
			if err == io.EOF {
				return hasher.Sum(nil), nil
			}
			return []byte{}, nil
		}
		if n == 1024 {
			hasher.Write(buffer)
		} else {
			hasher.Write(buffer[:n])
		}
	}
}

func removeNonDuplicatedPaths(duplicatedFiles DuplicatedFiles) {
	for key, value := range duplicatedFiles {
		if len(value) <= 1 {
			delete(duplicatedFiles, key)
		}
	}
}
