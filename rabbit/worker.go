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
	if err := w.Connect(); err != nil {
		return nil, err
	}
	if err := w.OpenChannel(); err != nil {
		return nil, err
	}
	if err := w.SetupPrefetch(); err != nil {
		return nil, err
	}
	if err := w.CreateQueue(); err != nil {
		return nil, err
	}
	if err := w.CreateExchange(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Worker) Connect() error {
	c := w.Config.Connection
	w.Logger.Info("Connecting RabbitMQ...")
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s", // amqp scheme
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Vhost,
	))
	if err != nil {
		return fmt.Errorf("Failed connecting RabbitMQ: %s", err)
	}
	w.Logger.Info("[Done]")
	w.Connection = conn
	return nil
}

func (w *Worker) OpenChannel() error {
	w.Logger.Info("Opening channel...")
	ch, err := w.Connection.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %s", err)
	}
	w.Logger.Info("[Done]")
	w.Channel = ch
	return nil
}

func (w *Worker) SetupPrefetch() error {
	p := w.Config.Worker.Prefetch
	w.Logger.Info("Setting QoS... ")
	if p.Count == 0 {
		p.Count = 3
	}
	err := w.Channel.Qos(
		p.Count,  // prefetchCount int
		0,        // prefetchSize int
		p.Global, // global bool
	)
	if err != nil {
		return fmt.Errorf("Failed to set QoS: %s", err)
	}
	w.Logger.Info("[Done]")
	return nil
}

func (w *Worker) CreateQueue() error {
	q := w.Config.Worker.Queue
	w.Logger.Info("Declaring queue [%s]...", q.Name)
	_, err := w.Channel.QueueDeclare(
		q.Name,       // name string
		q.Durable,    // durable bool
		q.AutoDelete, // autoDelete bool
		false,        // exclusive bool
		false,        // noWait bool
		nil,          // args Table
	)
	if err != nil {
		return fmt.Errorf("Failed to declare queue: %s", err)
	}
	w.Logger.Info("[Done]")
	return nil
}

func (w *Worker) CreateExchange() error {
	var err error

	e := w.Config.Worker.Exchange
	if e.Name == "" {
		w.Logger.Info("Empty Exchange name - use default exchange.")
		return nil
	}
	w.Logger.Info("Declaring exchange [%s]...", e.Name)
	if e.Type == "" {
		e.Type = "direct"
	}
	err = w.Channel.ExchangeDeclare(
		e.Name,       // name string
		e.Type,       // kind string
		e.Durable,    // durable bool
		e.AutoDelete, // autoDelete bool
		false,        // internal bool
		false,        // noWait bool
		amqp.Table{}, // args Table
	)
	if err != nil {
		return fmt.Errorf("Failed to declare exchange: %s", err)
	}
	w.Logger.Info("[Done]")

	q := w.Config.Worker.Queue
	w.Logger.Info("Binding queue [%s] to exchange [%s]...", q.Name, e.Name)
	err = w.Channel.QueueBind(
		q.Name, // name string
		"",     // key string
		e.Name, // exchange string
		false,  // noWait bool
		nil,    // args Table
	)
	if err != nil {
		return fmt.Errorf("Failed to bind queue to exchange: %s", err)
	}
	w.Logger.Info("[Done]")

	return nil
}

func (w *Worker) Consume() {
}
