package rabbit

import (
	"fmt"
	"net/url"

	"github.com/streadway/amqp"
)

type Worker struct {
	Config     *Config
	Logger     *Logger
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func NewWorker(c *Config, l *Logger) (*Worker, error) {
	w := &Worker{Config: c, Logger: l}
	if err := w.InitConnection(); err != nil {
		return nil, err
	}
	if err := w.InitWorker(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Worker) InitConnection() error {
	w.Logger.Info("Connecting RabbitMQ...")
	conn, err := amqp.Dial(w.Url())
	if nil != err {
		return fmt.Errorf("Failed connecting RabbitMQ: %s", err)
	}
	w.Logger.Info("Connected.")
	w.Connection = conn
	return nil
}

func (w *Worker) Url() string {
	c := w.Config.Connection
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Vhost,
	)
}

func (w *Worker) InitWorker() error {
	return nil
}

func (w *Worker) Consume() {
}
