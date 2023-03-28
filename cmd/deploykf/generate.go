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

const generateHelp = `
XXXXXXXXXX
XXXXXXXXXX
XXXXXXXXXX
XXXXXXXXXX
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
		Short: "XXXXXXXXXX",
		Long:  generateHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}

	// add local flags
	cmd.Flags().StringVarP(&o.sourceVersion, "source-version", "V", "", "XXXXXXX")
	cmd.Flags().StringVar(&o.sourcePath, "source-path", "", "")
	cmd.Flags().StringSliceVarP(&o.values, "values", "f", []string{}, "XXXXXXX")
	cmd.Flags().StringVarP(&o.outputDir, "output-dir", "O", "", "XXXXXXX")

	// mark local flags
	cmd.MarkFlagsMutuallyExclusive("source-version", "source-path")
	cmd.MarkFlagRequired("output-dir")

	return cmd
}

func (o *generateOptions) run(out io.Writer) error {
	// TODO: verify the provided `--values`:
	//  - ensure we test with multiple provided values files (also try with none)
	//  - check the YAML schema against a spec that is defined in the generator
	//  - check that all listed files exist
	//  - consider more complex verification, like mutually-exclusive fields, etc.

	// initialise the source helper
	// TODO: let users provide their own repo/owner
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
		DataSources:   o.gomplateDataSources(defaultValuesPath),
		Contexts:      o.gomplateContexts(),
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
		DataSources:   o.gomplateDataSources(defaultValuesPath),
		Contexts:      o.gomplateContexts(),
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
func (o *generateOptions) gomplateDataSources(defaultValuesPath string) []string {

	// create DataSource for each `--values` provided by the user
	values := o.values
	dataSources := make([]string, len(values)+1)
	for i, v := range values {
		dataSources[i] = "Values_" + strconv.Itoa(i) + "=" + v
	}

	// create DataSource for the default values
	dataSources[len(dataSources)-1] = "Values_default=" + defaultValuesPath

	return dataSources
}

// build the `Contexts` for our `gomplate.Config`
func (o *generateOptions) gomplateContexts() []string {
	var sb strings.Builder
	sb.WriteString("Values=merge:")

	// add the DataSources for any user provided `--values` to the merge
	values := o.values
	for i := range values {
		if i > 0 {
			sb.WriteString("|")
		}
		sb.WriteString("Values_")
		sb.WriteString(strconv.Itoa(i))
	}

	// add the DataSource for the default values to the merge
	if len(values) > 0 {
		sb.WriteString("|")
	}
	sb.WriteString("Values_default")

	return []string{sb.String()}
}

// build the `Templates` for our `gomplate.Config`
func (o *generateOptions) gomplateTemplates(helpersPath string) []string {
	return []string{"helpers=" + helpersPath}
}
