package deploykf

import (
	"fmt"
	"io"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/deployKF/cli/internal/require"
	"github.com/deployKF/cli/internal/version"
)

const versionHelp = `
This will print a representation the version of deployKF.

The output will look something like this:
version.BuildInfo{Version:"v1.0.0", GitCommit:"47a35c6b53cd5535ab9308ee29ab170866a2857d", GitTreeState:"clean", GoVersion:"go1.19.7"}
- Version is the semantic version of the release.
- GitCommit is the SHA for the commit that this version was built from.
- GitTreeState is "clean" if there are no local code changes when this binary was built, 
and "dirty" if the binary was built from locally modified code.
- GoVersion is the version of Go that was used to compile deployKF.

When using the --template flag, the following properties are available to use in the template:
- .Version contains the semantic version of deployKF
- .GitCommit is the git commit
- .GitTreeState is the state of the git tree when deployKF was built
- .GoVersion contains the version of Go that deployKF was compiled with

For example, --template='Version: {{.Version}}' outputs 'Version: v1.0.0'.
`

type versionOptions struct {
	short    bool
	template string
}

func newVersionCmd(out io.Writer) *cobra.Command {
	o := &versionOptions{}

	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print CLI version information",
		Long:  versionHelp,
		Args:  require.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.run(out)
		},
	}

	// add local flags
	cmd.Flags().BoolVar(&o.short, "short", false, "print the version number")
	cmd.Flags().StringVar(&o.template, "template", "", "template for version string format")

	return cmd
}

func (o *versionOptions) run(out io.Writer) error {
	if o.template != "" {
		tt, err := template.New("_").Parse(o.template)
		if err != nil {
			return err
		}
		return tt.Execute(out, version.Get())
	}
	_, err := fmt.Fprintln(out, formatVersion(o.short))
	if err != nil {
		return err
	}
	return nil
}

func formatVersion(short bool) string {
	v := version.Get()
	if short {
		if len(v.GitCommit) >= 7 {
			return fmt.Sprintf("%s+g%s", v.Version, v.GitCommit[:7])
		}
		return version.GetVersion()
	}
	return fmt.Sprintf("%#v", v)
}
