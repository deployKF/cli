package generate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// DeployKFOutputMarker is the name of the marker file that is created by deployKF in output directories,
	// the presence of this file indicates that the directory is safe to clean.
	DeployKFOutputMarker = ".deploykf_output"
)

type RunInfo struct {
	GeneratedAt   string `json:"generated_at"`
	SourceVersion string `json:"source_version,omitempty"`
	SourcePath    string `json:"source_path,omitempty"`
	SourceHash    string `json:"source_hash,omitempty"`
	CLIVersion    string `json:"cli_version"`
}

// CleanOutputDirectory cleans the output directory if it's safe to do so.
func CleanOutputDirectory(outputDir string) error {
	// Check if the output directory exists, and return if not.
	dirExists, err := DirectoryExists(outputDir)
	if err != nil {
		return err
	}
	if !dirExists {
		return nil
	}

	// Check if the output folder is empty, and return if so.
	files, err := os.ReadDir(outputDir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}

	// Output directory is non-empty, check if it contains a marker file, and fail if not.
	markerFile := filepath.Join(outputDir, DeployKFOutputMarker)
	markerFileExists, err := FileExists(markerFile)
	if err != nil {
		return err
	}
	if !markerFileExists {
		return fmt.Errorf("output directory '%s' is not safe to clean: no '%s' marker found", outputDir, DeployKFOutputMarker)
	}

	// Output directory is safe to clean, remove all files.
	for _, file := range files {
		err = os.RemoveAll(filepath.Join(outputDir, file.Name()))
		if err != nil {
			return err
		}
	}

	return nil
}

// CreateMarkerFile creates a marker file with RunInfo JSON in the output directory.
func CreateMarkerFile(outputDir string, sourceVersion string, sourcePath string, sourceHash string, cliVersion string) error {
	// Check if the output folder exists, and create it if not.
	outputDirExists, err := DirectoryExists(outputDir)
	if err != nil {
		return err
	}
	if !outputDirExists {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return err
		}
	}

	// Create the RunInfo struct.
	runInfo := RunInfo{
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		SourceVersion: sourceVersion,
		SourcePath:    sourcePath,
		SourceHash:    sourceHash,
		CLIVersion:    cliVersion,
	}

	// Serialize the struct to JSON.
	data, err := json.MarshalIndent(runInfo, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the marker file.
	filePath := filepath.Join(outputDir, DeployKFOutputMarker)
	err = os.WriteFile(filePath, data, 0644)
	return err
}
