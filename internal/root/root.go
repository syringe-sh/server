package root

import "github.com/spf13/cobra"

func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "syringe",
		Short: "🔐 Distributed database-per-user encrypted secrets management over SSH protocol.",
		Long:  "🔐 Distributed database-per-user encrypted secrets management over SSH protocol.",

		Example: `  Register user:
    syringe user register

  Add a project:
    syringe project add my_cool_project

  Add an environment:
    syringe environment add -p my_cool_project dev

  Add a secret:
    syringe secret set -p my_cool_project -e dev SECRET_KEY secret_value

  List secrets:
    syringe secret list -p my_cool_project -e dev

  Inject into command:
    syringe inject -p my_cool_project -e dev ./startserver

  For more examples, go to https://syringe.sh/examples`,

		// SilenceErrors: true,
	}

	additionalHelp := `
Please note: some commands are only available (and listed above) once registered and authenticated. Trying to use one of these while not authenticated will result in an 'unknown command' error.

For more help on how to use Syringe, go to https://syringe.sh/help

`

	rootCmd.SetHelpTemplate(rootCmd.HelpTemplate() + additionalHelp)

	return rootCmd
}