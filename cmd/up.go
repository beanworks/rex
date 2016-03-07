package cmd

import (
	"log"

	"github.com/beanworks/rex/rabbit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start hopping a rex rabbit consumer",
	Long: `Tell rex rabbit to start hopping, and consume messages from RabbitMQ.
A config file will need to be provided, and passed into this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := rabbit.NewLogger(&Config)
		if err != nil {
			log.Fatalf("Unabled to create logger: %s \n", err)
		}
		defer logger.Close()
		logger.Infof("Using config file: %s", viper.ConfigFileUsed())
		rex, err := rabbit.NewRex(&Config, logger)
		if err != nil {
			logger.Fatalf("Rex had some trouble starting to work: %s \n", err)
		}
		defer rex.Close()
		if err := rex.Consume(); err != nil {
			logger.Fatalf("Life is hard, and Rex said he couldn't consume any messages: %s \n", err)
		}
	},
}

func init() {
	RexCmd.AddCommand(upCmd)
}
