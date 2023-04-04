package generate

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// HashPath takes a path as input and returns the SHA-256 hash of the path.
// If the path is a file, it computes the hash of the file.
// If the path is a folder, it computes the hash of all files within the folder.
// It accepts an additional argument 'ignoreNames' which is a slice of file names to ignore.
func HashPath(path string, ignoreNames []string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return hashFile(path)
	}

	return hashDirectory(path, ignoreNames)
}

// hashFile takes a file path as input and returns the SHA-256 hash of the file.
func hashFile(filePath string) (string, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(fileContent)
	return hex.EncodeToString(hash[:]), nil
}

// hashDirectory takes a directory path and an ignoreNames slice as input and returns the SHA-256 hash
// of all files within the directory, excluding files with the specified names in the ignoreNames slice.
// It ensures a consistent hash across different operating systems (Windows, Linux, and macOS)
// by normalizing file paths, using case-insensitive sorting, and using relative paths when computing the hash.
func hashDirectory(dirPath string, ignoreNames []string) (string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && !contains(ignoreNames, info.Name()) {
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			normalizedPath := filepath.ToSlash(relPath)
			files = append(files, normalizedPath)
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	sort.SliceStable(files, func(i, j int) bool {
		return strings.ToLower(files[i]) < strings.ToLower(files[j])
	})

	hasher := sha256.New()
	for _, relPath := range files {
		absPath := filepath.Join(dirPath, filepath.FromSlash(relPath))
		fileHash, err := hashFile(absPath)
		if err != nil {
			return "", err
		}
		_, err = hasher.Write([]byte(relPath + fileHash))
		if err != nil {
			return "", err
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// contains checks if a slice of strings contains a specific string.
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
