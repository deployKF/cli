package generate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type GeneratorMarker struct {
	GeneratorSchema string `json:"generator_schema"`
}

// GetGeneratorSchemaVersion returns the `generator_schema` version from the specified marker file.
func GetGeneratorSchemaVersion(markerPath string) (string, error) {
	bytes, err := os.ReadFile(markerPath)
	if err != nil {
		return "", fmt.Errorf("failed to read generator marker file: %v", err)
	}

	var marker GeneratorMarker
	err = json.Unmarshal(bytes, &marker)
	if err != nil {
		return "", fmt.Errorf("failed to parse generator marker file: %v", err)
	}

	if marker.GeneratorSchema == "" {
		return "", errors.New("generator marker file is missing 'generator_schema' field")
	}

	return marker.GeneratorSchema, nil
}

// VerifyGeneratorSource verifies that the specified paths make a valid generator source,
// and that this version of the CLI supports the generator schema version.
func VerifyGeneratorSource(templatesPath string, helpersPath string, defaultValuesPath string, markerPath string) error {

	// Verify that the marker file exists
	markerFileExists, err := FileExists(markerPath)
	if err != nil {
		return err
	}
	if !markerFileExists {
		return fmt.Errorf("invalid generator source: marker file is missing")
	}

	// Verify that we support the generator schema version
	generatorSchemaVersion, err := GetGeneratorSchemaVersion(markerPath)
	if err != nil {
		return err
	}
	if generatorSchemaVersion != "v1" {
		return fmt.Errorf("invalid generator source: unsupported schema version '%s'", generatorSchemaVersion)
	}

	// Verify that the templates directory exists
	templatesDirExists, err := DirectoryExists(templatesPath)
	if err != nil {
		return err
	}
	if !templatesDirExists {
		return fmt.Errorf("invalid generator source: templates directory is missing")
	}

	// Verify that the helpers directory exists
	helpersDirExists, err := DirectoryExists(helpersPath)
	if err != nil {
		return err
	}
	if !helpersDirExists {
		return fmt.Errorf("invalid generator source: helpers directory is missing")
	}

	// Verify that the default values file exists
	defaultValuesFileExists, err := FileExists(defaultValuesPath)
	if err != nil {
		return err
	}
	if !defaultValuesFileExists {
		return fmt.Errorf("invalid generator source: default values file is missing")
	}

	return nil
}
