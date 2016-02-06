package cmd

import (
	"fmt"

	"github.com/beanworks/rex/rabbit"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start hopping a rex rabbit consumer",
	Long: `Tell rex rabbit to start hopping, and consume messages from RabbitMQ.
A config file will need to be provided, and passed into this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, err := rabbit.NewLogger(&Config)
		if err != nil {
			panic(fmt.Errorf("Unabled to create logger: %s \n", err))
		}
		worker, err := rabbit.NewWorker(&Config, logger)
		if err != nil {
			panic(fmt.Errorf("Rex had some trouble starting to work: %s \n", err))
		}
		worker.Consume()
	},
}

func init() {
	RexCmd.AddCommand(upCmd)
}
