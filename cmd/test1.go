package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trimmer-io/go-xmp/xmp"
)

var dirFlag string

var excludeDirs = []string{".git", ".vscode", ".idea", "node_modules", "vendor"}

var test1Cmd = &cobra.Command{
	Use:   "test1",
	Short: "Scan files for XMP metadata",
	Long:  `The test1 command scans files in a specified directory for XMP metadata and displays the file path and XMP data if found.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanFilesForXMP(dirFlag)
	},
}

func init() {
	rootCmd.AddCommand(test1Cmd)
	test1Cmd.Flags().StringVar(&dirFlag, "dir", ".", "Directory to scan for files")
}

func scanFilesForXMP(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if contains(excludeDirs, info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			slog.Error("Error opening file", "path", path, "error", err)
			return nil
		}
		defer file.Close()

		doc, err := xmp.Scan(file)
		if err != nil {
			if err.Error() == "xmp: no XMP found" {
				fmt.Printf("%s: No XMP metadata found\n", path)
			} else {
				slog.Error("Error scanning XMP metadata", "path", path, "error", err)
			}
			return nil
		}

		fmt.Printf("%s: XMP metadata found\n", path)
		fmt.Println("XMP Data:")
		jsonData, err := json.MarshalIndent(doc, "", "  ")
		if err != nil {
			slog.Error("Error marshalling XMP data to JSON", "path", path, "error", err)
			return nil
		}
		fmt.Println(string(jsonData))

		return nil
	})
	if err != nil {
		slog.Error("Error walking directory", "error", err)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
