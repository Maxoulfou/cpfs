/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cpfs/display"
	"cpfs/files"
	"cpfs/structures"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy all files and folders in a directory to a JSON file.",
	Long: `When you run the copy command, it will copy all files and folders in the
specified directory to a JSON file. The JSON file will be saved in the
result directory. The result directory will be created if it doesn't exist.`,
	Run: func(cmd *cobra.Command, args []string) {
		Copying(cmd, args)
	},
}

func init() {
	copyCmd.Flags().String("directory", "result", "The folder in which the result will be copied.")
	copyCmd.Flags().String("dir", ".", "The starting folder to read the folder structure.")
	copyCmd.Flags().String("filename", "", "The file name used for the final JSON file.")
	copyCmd.Flags().BoolP("files", "f", false, "Include files in the result.")

	rootCmd.AddCommand(copyCmd)

	viper.SetConfigType("json")
}

// Copying is the function that is called when the copy command is run.
func Copying(cmd *cobra.Command, args []string) {
	timeStart := time.Now()
	IncludeFiles, _ := cmd.Flags().GetBool("files")

	// Get the directory to copy.
	dir, _ := cmd.Flags().GetString("dir")

	// Create the result directory if it doesn't exist.
	resultDir := files.CheckDirectoryName(cmd)
	if _, err := os.Stat(resultDir); os.IsNotExist(err) {
		os.Mkdir(resultDir, 0755)
	}

	// Get the current date to use in the filename.
	filename := files.CheckFilename(cmd)

	// Walk the directory and create a Directory object representing the structure.
	root := &structures.Directory{Path: dir}
	Stat := &structures.Stats{}
	walk(root, dir, Stat, IncludeFiles)

	// Save the Directory object to a JSON file.
	filepathJoin := filepath.Join(resultDir, filename)
	data, _ := json.MarshalIndent(root, "", "  ")
	if err := os.WriteFile(filepathJoin, data, 0644); err != nil {
		fmt.Printf("Failed to write to file: %s\n", err)
		return
	}

	display.DisplayResults(filepathJoin, Stat, timeStart)
}

// Recursively walk the directory and create a Directory object representing the structure.
func walk(dir *structures.Directory, path string, stat *structures.Stats, IncludeFiles bool) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Failed to open directory %s: %s\n", path, err)
		return
	}

	wg := sync.WaitGroup{}
	for _, entry := range entries {
		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()
			if entry.IsDir() {
				subdivide := &structures.Directory{Path: entry.Name()}
				dir.Subdirs = append(dir.Subdirs, subdivide)
				walk(subdivide, filepath.Join(path, entry.Name()), stat, IncludeFiles)
				stat.Directories++
			} else if IncludeFiles {
				dir.Files = append(dir.Files, entry.Name())
				stat.Files++
				fileInfo, err := entry.Info()
				if err != nil {
					fmt.Printf("Failed to get file info for %s: %s\n", entry.Name(), err)
					return
				}
				stat.TotalSize += fileInfo.Size()
			}
		}(entry)
	}
	wg.Wait()
}
