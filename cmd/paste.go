/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"cpfs/structures"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// pasteCmd represents the paste command
var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Paste("result/test.json", "result")
		if err != nil {
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(pasteCmd)
}

func Paste(jsonFile string, destPath string) error {
	// Open the JSON file containing the directory structure.
	file, err := os.Open(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to open JSON file: %s", err)
	}
	defer file.Close()

	// Decode the JSON data into a directory structure.
	var rootDir structures.Directory
	err = json.NewDecoder(file).Decode(&rootDir)
	if err != nil {
		return fmt.Errorf("failed to decode JSON data: %s", err)
	}

	// Recursively create directories and copy files.
	err = copyDirectory(rootDir, destPath)
	if err != nil {
		return fmt.Errorf("failed to copy directory: %s", err)
	}

	return nil
}

// copyDirectory recursively creates directories and copies files.
func copyDirectory(dir structures.Directory, parentPath string) error {
	// Create the new directory.
	if dir.Path == "." {
		dir.Path = "source"
	}

	newPath := filepath.Join(parentPath, dir.Path)
	err := os.Mkdir(newPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %s: %s", newPath, err)
	}

	// Recursively create subdirectories.
	for _, subdir := range dir.Subdirs {
		fmt.Printf("Creating directory %s\n", *subdir)
		err = copyDirectory(*subdir, newPath)
		if err != nil {
			return err
		}
	}

	// Copy files to the new directory.
	for _, file := range dir.Files {
		err = copyFile(filepath.Join(parentPath, file), filepath.Join(newPath, file))
		if err != nil {
			return err
		}
	}

	return nil
}

// copyFile copies a file from src to dest.
func copyFile(src string, dest string) error {
	// Open the source file.
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %s", src, err)
	}
	defer srcFile.Close()

	// Create the destination file.
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %s", dest, err)
	}
	defer destFile.Close()

	// Copy the data from source to destination.
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy data from %s to %s: %s", src, dest, err)
	}

	return nil
}
