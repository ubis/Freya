package directory

import (
	"os"
	"path/filepath"
	"strings"
)

// Returns Current working directory
func Current() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "nil"
	}
	return dir + "/"
}

// Returns Root directory
func Root() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "bin", "", 1)
}
