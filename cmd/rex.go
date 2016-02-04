package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version bool
var CfgFile string

var RexCmd = &cobra.Command{
	Use:   "rex",
	Short: "Rex rabbit likes consuming queued messages",
	Long: `Rex rabbit is a command line message queue consumer for RabbitMQ.
Rex pulls messages from a queue, takes a good care of the jobs,
redirects message bodies to other responsible parties.

When Rex is not busy, he also likes to hang out with Octocat.`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			cmd.Println(GetVersionString())
		} else {
			cmd.Help()
		}
	},
}

func Execute() {
	if err := RexCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RexCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file (default is $HOME/.rex.yml)")
	RexCmd.Flags().BoolVarP(&Version, "version", "v", false, "Show rex version")

	RexCmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})
}

func initConfig() {
	if CfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(CfgFile)
	}

	viper.SetConfigName(".rex")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
