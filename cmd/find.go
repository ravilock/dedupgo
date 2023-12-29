package cmd

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/ravilock/dedupgo/internal"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively search through all subdirectories for duplicated files")
	findCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Define directory depth limit to search for duplicated files")

	findCmd.MarkFlagsMutuallyExclusive("recursive", "depth")
}

var recursive bool
var depth int

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find duplicate files in passed directory",
	Long:  "Iterate through directory searching for duplicated files",
	Args:  cobra.MaximumNArgs(1),
	RunE:  find,
}

func find(cmd *cobra.Command, args []string) error {
	directoryQueue := &internal.Queue[string]{}

	dirPath, err := directoryPath(args)
	if err != nil {
		return err
	}

	directoryQueue.Enqueue(dirPath, 0)

	files := make([]string, 0)

	for directoryQueue.Length() != 0 {
		currentDir, currentDepth := directoryQueue.Dequeue()
		if currentDir == "" || (depth != 0 && currentDepth == depth) {
			break
		}

		dirEntries, err := os.ReadDir(currentDir)
		if err != nil {
			return err
		}

		for _, entry := range dirEntries {
			entryPath := path.Join(currentDir, entry.Name())
			if entry.IsDir() {
				directoryQueue.Enqueue(entryPath, currentDepth+1)
			} else {
				files = append(files, entryPath)
			}
		}

		if !recursive && depth == 0 {
			break
		}
	}

	filesByBash, err := groupFilesByContentHash(files, dirPath)
	if err != nil {
		return err
	}
	removeNonDuplicatedFiles(filesByBash)
	if len(filesByBash) == 0 {
		fmt.Println("No duplicates found")
		return nil
	}
	for _, value := range filesByBash {
		fmt.Printf("%s are duplicate files\n", value)
	}
	return nil
}

func directoryPath(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	}
	return args[0], nil
}

func getFiles(directory string) ([]os.DirEntry, error) {
	return nil, nil
}

func filterFiles(entries []os.DirEntry) []os.DirEntry {
	filtered := make([]os.DirEntry, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		filtered = append(filtered, entry)
	}
	return filtered
}

func groupFilesByContentHash(files []string, dirPath string) (map[string][]string, error) {
	contentHashToFileName := make(map[string][]string)
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

func removeNonDuplicatedFiles(filesByHash map[string][]string) {
	for key, value := range filesByHash {
		if len(value) <= 1 {
			delete(filesByHash, key)
		}
	}
}
