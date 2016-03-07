package cmd

import (
	"fmt"

	"github.com/beanworks/rex/rabbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start hopping a rex rabbit consumer",
	Long: `Tell rex rabbit to start hopping, and consume messages from RabbitMQ.
A config file will need to be provided, and passed into this command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger, err := rabbit.NewLogger(&Config)
		if err != nil {
			return fmt.Errorf("Unabled to create logger: %s \n", err)
		}
		defer logger.Close()

		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
		rex, err := rabbit.NewRex(&Config, logger)
		if err != nil {
			return fmt.Errorf("Rex encountered some trouble to start up: %s \n", err)
		}
		defer rex.Close()
		if err := rex.Consume(); err != nil {
			return fmt.Errorf("Life is hard! Rex couldn't consume any messages. \nSee the reason: %s \n", err)
		}

		return nil
	},
}

func init() {
	RexCmd.AddCommand(upCmd)
}
