package cmd

import (
	"github.com/beanworks/rex/rabbit"
	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start hopping a rex rabbit consumer",
	Long: `Tell rex rabbit to start hopping, and consume messages from RabbitMQ.
A config file will need to be provided, and passed into this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger, _ := rabbit.NewLogger(&Config)
		worker, _ := rabbit.NewWorker(&Config, logger)
		worker.Consume()
	},
}

func init() {
	RexCmd.AddCommand(upCmd)
}
