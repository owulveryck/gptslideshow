package main

import "fmt"

// printProgress displays a progress bar in the terminal.
// `current` is the current progress, and `total` is the maximum value of the progress.
func printProgress(current, total int) {
	// Calculate the percentage of progress
	percentage := float64(current) / float64(total) * 100

	// Calculate the number of filled blocks in the progress bar
	barLength := 20 // Length of the progress bar
	filledLength := int(float64(barLength) * percentage / 100)

	// Create the progress bar string
	bar := ""
	for i := 0; i < filledLength; i++ {
		bar += "#"
	}
	for i := filledLength; i < barLength; i++ {
		bar += " "
	}

	// Print the progress bar and overwrite the current line
	fmt.Printf("\r|%s| %.0f%%", bar, percentage)
}
