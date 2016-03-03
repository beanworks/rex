package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/beanworks/rex/rabbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Version bool
var CfgFile string
var Config rabbit.Config

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

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/rex/")
	viper.AddConfigPath("$HOME/.rex")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error: config file: %s \n", err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("Unable to decode config into struct, %s \n", err)
	}
}
