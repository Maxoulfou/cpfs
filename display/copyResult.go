package display

import (
	"cpfs/structures"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"time"
)

// Purpose improved DisplayResults function
func DisplayResults(filepath string, stat *structures.Stats, time_start time.Time) {
	// Format the total size in MB with 2 decimal places
	totalSize := fmt.Sprintf("%.2f Mo", float64(stat.TotalSize)/1024/1024)

	// Create bulk data
	data := [][]string{
		{"Filepath", filepath},
		{"Total Files", strconv.Itoa(stat.Files)},
		{"Total Directories", strconv.Itoa(stat.Directories)},
		{"Total Size", totalSize},
		{"Time Elapsed", time.Since(time_start).String()},
	}

	// Create the table writer with merge cells option
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetBorder(true)
	table.SetTablePadding(" ")

	// Merge the cells of the header
	table.SetHeaderLine(true)
	table.SetHeader([]string{"RESULTS"})

	// Append the data to the table
	table.AppendBulk(data)

	// Render the table
	table.Render()
}
