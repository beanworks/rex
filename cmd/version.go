package cmd

import "github.com/spf13/cobra"

const currentVersion = "1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show rex version",
	Long:  "Show rex version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(GetVersionString())
	},
}

func init() {
	RexCmd.AddCommand(versionCmd)
}

func GetVersionString() string {
	return currentVersion
}
