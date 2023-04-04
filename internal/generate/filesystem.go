package generate

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) (isDir bool, isFile bool, err error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return info.IsDir(), !info.IsDir(), nil
}

func DirectoryExists(dir string) (bool, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return info.IsDir(), nil
}

func FileExists(file string) (bool, error) {
	info, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !info.IsDir(), nil
}

// UnzipFile extracts the contents of a .zip file to a destination directory
// extractPath is the relative path inside the zip archive that should be extracted
// If extractPath does not match any files or directories in the zip archive, an error is returned
func UnzipFile(zipFilePath string, targetDir string, extractPath string) error {
	// Open the zip file
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Normalize the extractPath
	// NOTE: zip files always use forward slashes, so we need to convert the OS-specific path separator
	extractPath = filepath.Clean(filepath.FromSlash(extractPath))

	// Flag to track if any files or directories have been extracted
	extracted := false

	// Iterate through each file in the zip archive
	for _, file := range reader.File {
		// Normalize file.Name before processing
		// NOTE: zip files always use forward slashes, so we need to convert the OS-specific path separator
		normalizedFileName := filepath.Clean(filepath.FromSlash(file.Name))

		// Check if the file is within the extractPath
		if !strings.HasPrefix(normalizedFileName, extractPath) {
			continue
		}

		// Mark the extracted flag as true
		extracted = true

		// Create the target file path
		relPath := strings.TrimPrefix(normalizedFileName, extractPath)
		cleanTargetFilePath := filepath.Join(targetDir, relPath)
		finalTargetFilePath, err := filepath.Rel(targetDir, filepath.Clean(cleanTargetFilePath))

		// Check if the file path is valid
		// NOTE: We need to check if the file path is valid because it's possible to have a zip file with
		//       a file path that goes outside the target directory. For example, if the target directory
		//       is /tmp/foo and the file path is ../../bar, the final target file path would be /bar.
		//       This is a security issue, so we need to check for it.
		if err != nil || strings.HasPrefix(finalTargetFilePath, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path: %s", file.Name)
		}

		// Create the target file's directory if it doesn't exist
		targetDir := filepath.Dir(cleanTargetFilePath)
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			return err
		}

		// If it's a directory, move to the next file
		if file.FileInfo().IsDir() {
			continue
		}

		// Open the file inside the zip archive for reading
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		// Create the target file for writing
		targetFile, err := os.OpenFile(cleanTargetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// Copy the contents of the file inside the zip archive to the target file
		_, err = io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
	}

	// Check if any files or directories have been extracted
	if !extracted {
		return fmt.Errorf("the provided extractPath '%s' does not exist within the zip", extractPath)
	}

	return nil
}

// CopyFolder recursively copies the contents of the source folder to the destination folder
func CopyFolder(src, dest string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get the relative path of the current file/directory within the source folder
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create the destination path by joining the destination folder and the relative path
		destPath := filepath.Join(dest, relPath)

		// If the current item is a directory, create the corresponding directory in the destination folder
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// If the current item is a file, copy it to the destination folder
		return copyFile(path, destPath)
	})

	return err
}

// copyFile copies an individual file from the source path to the destination path
func copyFile(src, dest string) error {
	// Open the source file for reading
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get file info (including mode/permissions) of the source file
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create the destination file with the same mode/permissions as the source file
	destFile, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destFile, srcFile)
	return err
}
