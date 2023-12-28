package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	dirPath, err := directoryPath()
	if err != nil {
		log.Fatalln(err)
	}
	dirEntries, err := os.ReadDir(dirPath)
	if err != nil {
		log.Fatalln(err)
	}
	files := filterFiles(dirEntries)
	filesByBash, err := groupFilesByContentHash(files, dirPath)
	if err != nil {
		log.Fatalln(err)
	}
	removeNonDuplicatedFiles(filesByBash)
	if len(filesByBash) == 0 {
		fmt.Println("No duplicates found")
		os.Exit(0)
	}
	for _, value := range filesByBash {
		fmt.Printf("%s are duplicate files\n", value)
	}
	os.Exit(0)
}

func directoryPath() (string, error) {
	if len(os.Args) == 1 {
		return os.Getwd()
	}
	return os.Args[1], nil
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

func groupFilesByContentHash(files []os.DirEntry, dirPath string) (map[string][]string, error) {
	contentHashToFileName := make(map[string][]string)
	for _, file := range files {
		hash2, err := hashFile(file, dirPath)
		if err != nil {
			return nil, err
		}
		hashString := hex.EncodeToString(hash2)
		_, ok := contentHashToFileName[hashString]
		if !ok {
			contentHashToFileName[hashString] = []string{file.Name()}
		} else {
			contentHashToFileName[hashString] = append(contentHashToFileName[hashString], file.Name())
		}
	}
	return contentHashToFileName, nil
}

func hashFile(file os.DirEntry, dirPath string) ([]byte, error) {
	buffer := make([]byte, 1024)
	f, err := os.Open(path.Join(dirPath, file.Name()))
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
