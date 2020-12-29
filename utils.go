package main

import (
	"fmt"
	"os"
)

func cleanUpPartitionFiles(partitionFiles []string) {
	for _, f := range partitionFiles {
		err := os.Remove(f)
		if err != nil {
			fmt.Printf("Error removing file %s: %v", f, err)
		}
	}
}
