// completion clicmd recycled from:
// https://raw.githubusercontent.com/cli/cli/master/command/completion.go

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/schnoddelbotz/cloud-task-zip-zap/cmd/cobrafish"
)

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().StringP("shell", "s", "bash", "The type of shell (bash, zsh or fish)")
}

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "generates shell completion scripts",
	Long: `To enable tab-completion for gotsfn in bash:
  eval "$(cloud-task-zip-zap completion)"
In fish shell:
  cloud-task-zip-zap completion -s fish | .

You can add that to your '~/.bash_profile', '~/.config/fish/config.fish' etc.
to enable completion whenever you start a new shell.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		shellType, err := cmd.Flags().GetString("shell")
		if err != nil {
			return err
		}

		switch shellType {
		case "bash":
			return rootCmd.GenBashCompletion(cmd.OutOrStdout())
		case "zsh":
			return rootCmd.GenZshCompletion(cmd.OutOrStdout())
		case "fish":
			return cobrafish.GenCompletion(rootCmd, cmd.OutOrStdout())
		default:
			return fmt.Errorf("unsupported shell type %q", shellType)
		}
	},
}
