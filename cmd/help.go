package cmd

import "github.com/spf13/cobra"

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "help for rex",
	Long:  "help for rex",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Parent().Help()
	},
}

func init() {
	RexCmd.AddCommand(helpCmd)
}
