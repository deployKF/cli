package deploykf

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

const rootHelp = `deployKF is your open-source helper for deploying MLOps tools on Kubernetes.

Common actions for deployKF:

- deploykf generate:  Generate Kubernetes manifests from deployKF templates and config values

The default directories depend on the Operating System. The defaults are listed below:

| Operating System | Assets Cache Path              |
|------------------|--------------------------------|
| Linux            | $HOME/.deploykf/assets         |
| macOS            | $HOME/.deploykf/assets         |
| Windows          | %userprofile%\.deploykf\assets |
`

func newRootCmd(out io.Writer) *cobra.Command {
	var cmd = &cobra.Command{
		Use:          "deploykf",
		Short:        "deployKF is your open-source helper for deploying MLOps tools on Kubernetes",
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
