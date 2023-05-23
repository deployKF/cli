package generate

import (
	"os"
	"path/filepath"
)

const (
	InputDirTemplateFile  = "input_dir"
	OutputDirTemplateFile = "output_dir"
)

// WriteRuntimeTemplates writes the runtime templates to the specified directory
func WriteRuntimeTemplates(runtimeTemplatePath string, inputDirConfig string, outputDirConfig string) error {
	// check if the runtime templates folder exists, and create it if not
	runtimeDirExists, err := DirectoryExists(runtimeTemplatePath)
	if err != nil {
		return err
	}
	if !runtimeDirExists {
		err = os.MkdirAll(runtimeTemplatePath, 0755)
		if err != nil {
			return err
		}
	}

	// write `--input-dir` config as a template file
	inputDirTemplatePath := filepath.Join(runtimeTemplatePath, InputDirTemplateFile)
	err = os.WriteFile(inputDirTemplatePath, []byte(inputDirConfig), 0644)
	if err != nil {
		return err
	}

	// write `--output-dir` config as a template file
	outputDirTemplatePath := filepath.Join(runtimeTemplatePath, OutputDirTemplateFile)
	err = os.WriteFile(outputDirTemplatePath, []byte(outputDirConfig), 0644)
	if err != nil {
		return err
	}

	return nil
}
