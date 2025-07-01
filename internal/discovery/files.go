package discovery

import (
	"os"
	"path/filepath"
	"strings"
)

// FindMatchingFiles recursively searches for files matching the given extensions
// in the specified directory and returns their paths.
func FindMatchingFiles(dir string, extensions []string) ([]string, error) {
	var matchingFiles []string
	
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if MatchesExtension(path, extensions) {
				matchingFiles = append(matchingFiles, path)
			}
		}
		return nil
	})

	return matchingFiles, err
}

// MatchesExtension checks if a file path has an extension that matches
// any of the provided extensions (case-insensitive).
func MatchesExtension(filePath string, extensions []string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, configExt := range extensions {
		if ext == strings.ToLower(configExt) {
			return true
		}
	}
	return false
}