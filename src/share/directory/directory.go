package directory

import (
	"os"
	"strings"
)

// Returns Current working directory
func Current() string {
	dir, err := os.Getwd()
	if err != nil {
		return "nil"
	}
	
	return dir
}

// Returns Root directory
func Root() string {
	dir, _ := os.Getwd()
	return strings.Replace(dir, "bin", "", 1)
}
