package cmd

import (
	"fmt"
	"os"

	"github.com/ravilock/dedupgo/dedup"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively search through all subdirectories for duplicated files")
	findCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Define sub-directory depth limit to search for duplicated files")

	findCmd.MarkFlagsMutuallyExclusive("recursive", "depth")
}

var recursive bool
var depth int

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find duplicate files in passed directory",
	Long:  "Iterate through directory searching for duplicated files",
	Args:  cobra.MaximumNArgs(1),
	RunE:  findCommand,
}

type visitedDirectory struct {
	path  string
	depth int
}

func findCommand(cmd *cobra.Command, args []string) error {
	dirPath, err := directoryPath(args)
	if err != nil {
		return err
	}

	findOptions := &dedup.FindOptions{}
	if recursive {
		findOptions.Mode = dedup.ModeRecursive
	} else if depth != 0 {
		findOptions.Mode = dedup.ModeDepth
		findOptions.DepthLimit = depth
	}

	filesByHash, err := dedup.Find(dirPath, findOptions)
	if err != nil {
		return err
	}

	if len(filesByHash) == 0 {
		fmt.Println("No duplicates found")
		return nil
	}

	for hash, value := range filesByHash {
		fmt.Printf("%s are duplicated files under the hash %q\n", value, hash)
	}
	return nil
}

func directoryPath(args []string) (string, error) {
	if len(args) == 0 {
		return os.Getwd()
	}
	return args[0], nil
}
