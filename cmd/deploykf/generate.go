package deploykf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hairyhenderson/gomplate/v3"
	"github.com/spf13/cobra"

	"github.com/deployKF/cli/internal/generate"
	"github.com/deployKF/cli/internal/version"
)

const generateHelp = `This command will generate an output folder containing Kubernetes manifests.

ARGUMENTS:
----------------

You must provide either '--source-version' OR '--source-path' to specify the source of the generator:
 - If '--source-version' is provided, the provided version tag will be downloaded from the 'deployKF/deployKF' GitHub.
 - If '--source-path' is provided, the source will be read from the provided local directory or '.zip' file.

You may provide one or more '--values' files that contain your configuration values:
 - For more information on how to structure your values files, see the 'deployKF/deployKF' GitHub repository.

You must provide '--output-dir' to specify the output directory for the generated manifests:
 - If the directory does not exist, it will be created.
 - If the directory is non-empty, it will be cleaned before generating the manifests.
   However, it must contain a '.deploykf_output' marker file, otherwise the command will fail.

OUTPUT:
----------------

The '.deploykf_output' marker file contains the following information:
 - generated_at: the time the generator was run
 - source_version: the source version that was used (if '--source-version' was provided)
 - source_path: the path of the source artifact that was used 
 - source_hash: the SHA256 hash of the source artifact that was used
 - cli_version: the version of the deployKF CLI that was used

EXAMPLES:
----------------

To generate manifests from a GitHub source version:

    $ deploykf generate --source-version v0.1.0 --values ./values.yaml --output-dir ./GENERATOR_OUTPUT

To generate manifests from a local source zip file:

    $ deploykf generate --source-path ./deploykf.zip --values ./values.yaml --output-dir ./GENERATOR_OUTPUT

To generate manifests from a local source directory:

    $ deploykf generate --source-path ./deploykf --values ./values.yaml --output-dir ./GENERATOR_OUTPUT
`

type generateOptions struct {
	sourceVersion string
	sourcePath    string
	values        []string
	outputDir     string
}

func newGenerateCmd(out io.Writer) *cobra.Command {
	o := &generateOptions{}

	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate Kubernetes manifests from deployKF templates and config values",
		Long:  generateHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}

	// add local flags
	cmd.Flags().StringVarP(&o.sourceVersion, "source-version", "V", "", "a version tag from the 'deployKF/deployKF' GitHub repository")
	cmd.Flags().StringVar(&o.sourcePath, "source-path", "", "a local path to a directory or '.zip' file containing a generator source")
	cmd.Flags().StringSliceVarP(&o.values, "values", "f", []string{}, "a YAML file containing configuration values")
	cmd.Flags().StringVarP(&o.outputDir, "output-dir", "O", "", "the output directory in which to generate the manifests")

	// mark local flags
	cmd.MarkFlagsMutuallyExclusive("source-version", "source-path")
	cmd.MarkFlagRequired("output-dir")

	return cmd
}

func (o *generateOptions) run(out io.Writer) error {
	// TODO: verify the provided `--values`:
	//  - check the YAML schema against a spec that is defined in the generator source
	//  - check that all provided file paths exist (before gomplate fails)

	// initialise the source helper
	// TODO: let users provide their own repo/owner for the source
	sourceHelper := generate.NewSourceHelper()

	// create a temporary directory to store our generator source,
	// and defer a function to clean it up after this function returns
	tempSourcePath, err := os.MkdirTemp("", "deploykf-generator-source-*")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		return nil
	}
	defer func() {
		err := os.RemoveAll(tempSourcePath)
		if err != nil {
			fmt.Printf("Error removing temporary directory: %v\n", err)
		}
	}()

	// populate the temporary directory with the generator source
	//  - CASE 1: if `--source-version` is provided, download that version's `.zip` file and unzip it into the temp folder
	//  - CASE 2: if `--source-path` points to a `.zip` file, unzip it into the temp folder
	//  - CASE 3: if `--source-path` points to a folder, copy the contents of that folder into the temp folder
	var sourcePath string
	if o.sourceVersion != "" {
		// CASE 1: download the source from GitHub
		sourcePath, err = sourceHelper.DownloadAndUnpackSource(o.sourceVersion, tempSourcePath, out)
		if err != nil {
			return err
		}
	} else if o.sourcePath != "" {
		sourcePath, err = filepath.EvalSymlinks(o.sourcePath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("the provided --source-path '%s' does not exist", o.sourcePath)
			}
			return err
		}
		sourceIsDir, sourceIsFile, err := generate.PathExists(sourcePath)
		if err != nil {
			return err
		}
		if sourceIsFile && strings.HasSuffix(sourcePath, ".zip") {
			// CASE 2: source is a .zip file
			fmt.Fprintf(out, "Using custom source file: %s\n", o.sourcePath)
			err := generate.UnzipFile(sourcePath, tempSourcePath, "generator")
			if err != nil {
				return err
			}
		} else if sourceIsDir {
			// CASE 3: source is a folder
			fmt.Fprintf(out, "Using custom source folder: %s\n", o.sourcePath)
			err := generate.CopyFolder(sourcePath, tempSourcePath)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("the provided --source-path '%s' must be a folder or a .zip file", o.sourcePath)
		}
	} else {
		return fmt.Errorf("at least one of `--source-version` or `--source-path` must be provided")
	}

	// important paths from the generator source
	templatesPath := filepath.Join(tempSourcePath, "templates")
	helpersPath := filepath.Join(tempSourcePath, "helpers")
	defaultValuesPath := filepath.Join(tempSourcePath, "default_values.yaml")

	// GENERATOR PHASE 1: render `.gomplateignore_template` files
	//  - note, we are rendering the `.gomplateignore` files into the generator source
	//    templates folder, not the output folder
	phase1Config := &gomplate.Config{ //nolint:staticcheck
		InputDir:      templatesPath,
		OutputMap:     templatesPath + `/{{< .in | strings.ReplaceAll ".gomplateignore_template" ".gomplateignore" >}}`,
		ExcludeGlob:   []string{"*", "!*.gomplateignore_template"},
		LDelim:        "{{<",
		RDelim:        ">}}",
		DataSources:   o.gomplateDataSources(),
		Contexts:      o.gomplateContexts(defaultValuesPath),
		Templates:     o.gomplateTemplates(helpersPath),
		SuppressEmpty: true,
	}
	err = gomplate.RunTemplates(phase1Config) //nolint:staticcheck
	if err != nil {
		return err
	}

	// clean the `--output-dir` if it's safe to do so
	err = generate.CleanOutputDirectory(o.outputDir)
	if err != nil {
		return err
	}

	// calculate the hash of the generator source
	//  - if the source was a `.zip` file, we'll use the hash of the file
	//  - if the source was a folder, we'll use the hash of the folder
	//    note, we'll ignore the `.gomplateignore` files when calculating the hash
	//    see `generate.HashPath` for more details
	sourceArtifactHash, err := generate.HashPath(sourcePath, []string{".gomplateignore"})
	if err != nil {
		return err
	}

	// create marker file in the `--output-dir`
	//  - will create the directory if it doesn't already exist
	//  - the marker will contain JSON with information like run time and source version
	err = generate.CreateMarkerFile(o.outputDir, o.sourceVersion, sourcePath, sourceArtifactHash, version.GetVersion())
	if err != nil {
		return err
	}

	// GENERATOR PHASE 2: render to the output folder
	phase2Config := &gomplate.Config{ //nolint:staticcheck
		InputDir:      templatesPath,
		OutputDir:     o.outputDir,
		LDelim:        "{{<",
		RDelim:        ">}}",
		DataSources:   o.gomplateDataSources(),
		Contexts:      o.gomplateContexts(defaultValuesPath),
		Templates:     o.gomplateTemplates(helpersPath),
		SuppressEmpty: true,
	}
	err = gomplate.RunTemplates(phase2Config) //nolint:staticcheck
	if err != nil {
		return err
	}

	// log the output directory
	fmt.Fprintf(out, "Generated manifests at: %s\n", o.outputDir)

	return nil
}

// build the `DataSources` for our `gomplate.Config`
func (o *generateOptions) gomplateDataSources() []string {

	// create DataSource for each `--values` provided by the user
	values := o.values
	dataSources := make([]string, len(values))
	for i, v := range values {
		dataSources[i] = "Values_" + strconv.Itoa(i) + "=" + v
	}

	return dataSources
}

// build the `Contexts` for our `gomplate.Config`
func (o *generateOptions) gomplateContexts(defaultValuesPath string) []string {
	var sb strings.Builder

	// merge the user-provided `--values` DataSources by placing a "|" between them
	// NOTE: merges happen from right to left, so we add the `--values` in reverse
	//       order so the `--values` that were provided later take precedence
	values := o.values
	for i := len(values) - 1; i >= 0; i-- {
		if sb.Len() > 0 {
			sb.WriteString("|")
		}
		sb.WriteString("Values_")
		sb.WriteString(strconv.Itoa(i))
	}

	// add the `default_values.yaml` to the end of the merge
	if sb.Len() > 0 {
		sb.WriteString("|")
	}
	sb.WriteString(defaultValuesPath)

	// we only use `merge` if there is more than one DataSource
	valuesContext := "Values="
	if len(values) > 0 {
		valuesContext += "merge:"
	}
	valuesContext += sb.String()

	return []string{valuesContext}
}

// build the `Templates` for our `gomplate.Config`
func (o *generateOptions) gomplateTemplates(helpersPath string) []string {
	return []string{"helpers=" + helpersPath}
}
