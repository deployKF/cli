package deploykf

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

const rootHelp = `
XXXXXXXXXX
XXXXXXXXXX
XXXXXXXXXX
XXXXXXXXXX
`

func newRootCmd(out io.Writer) *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "deploykf",
		Short:        "XXXXXXXXXX",
		Long:         rootHelp,
		SilenceUsage: true,
	}

	// add subcommands
	cmd.AddCommand(
		newGenerateCmd(out),
		newVersionCmd(out),
	)

	return cmd
}

// Execute instantiates and runs root command, this is called by main.main()
func Execute() {
	rootCmd := newRootCmd(os.Stdout)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
