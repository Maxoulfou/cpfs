package files

import (
	"fmt"
	"github.com/spf13/cobra"
	"time"
)

// GenerateFilename generates a filename based on the current date time.
func GenerateFilename() string {
	now := time.Now()
	return fmt.Sprintf("structure-%04d-%02d-%02d--%02d_%02d_%02d.json",
		now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// CheckFilename verify if the filename flag is set, and return the value if it is, or generate a filename if it isn't.
func CheckFilename(cmd *cobra.Command) string {
	if filename, _ := cmd.Flags().GetString("filename"); filename != "" {
		return fmt.Sprintf("%s.json", filename)
	}
	return GenerateFilename()
}

func CheckDirectoryName(cmd *cobra.Command) string {
	if dir, _ := cmd.Flags().GetString("directory"); dir != "" {
		return dir
	}
	return "result"
}
