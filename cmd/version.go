package cmd

import (
	"runtime"

	"github.com/spf13/cobra"
)

const currentVersion = "0.5.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show rex version",
	Long:  "Show rex version",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf(
			"rex version %s %s/%s \n",
			GetVersionString(),
			runtime.GOOS,
			runtime.GOARCH,
		)
	},
}

func init() {
	RexCmd.AddCommand(versionCmd)
}

func GetVersionString() string {
	return currentVersion
}
