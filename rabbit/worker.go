package rabbit

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

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
	if err := w.connect(); err != nil {
		return nil, err
	}
	if err := w.createQueueAndExchange(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Worker) connect() (err error) {
	c := w.Config.Connection
	p := w.Config.Worker.Prefetch

	w.Logger.Info("Connecting to RabbitMQ server...")
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s", // amqp scheme
		url.QueryEscape(c.Username),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Vhost,
	))
	if err != nil {
		return
	}

	w.Logger.Info("Opening channel...")
	chn, err := conn.Channel()
	if err != nil {
		return
	}

	w.Logger.Info("Setting QoS... ")
	if p.Count == 0 {
		p.Count = 3
	}
	// Args: prefetchCount, prefetchSize int, global bool
	if err = chn.Qos(p.Count, 0, p.Global); err != nil {
		return
	}

	w.Connection = conn
	w.Channel = chn

	return
}

func (w *Worker) createQueueAndExchange() (err error) {
	q := w.Config.Worker.Queue
	e := w.Config.Worker.Exchange

	// Create queue
	w.Logger.Info("Declaring queue [%s]...", q.Name)
	// Args: name string, durable, autoDelete, exclusive, noWait bool, args Table
	_, err = w.Channel.QueueDeclare(
		q.Name, q.Durable, q.AutoDelete, false, false, nil)
	if err != nil {
		return
	}

	// Create exchange
	if e.Name == "" {
		w.Logger.Info("Empty Exchange name - use default exchange.")
		return
	}
	w.Logger.Info("Declaring exchange [%s]...", e.Name)
	if e.Type == "" {
		e.Type = "direct"
	}
	// Args: name, kind string, durable, autoDelete, internal, noWait bool, args Table
	err = w.Channel.ExchangeDeclare(
		e.Name, e.Type, e.Durable, e.AutoDelete, false, false, nil)
	if err != nil {
		return
	}

	// Bind queue and exchange
	w.Logger.Info("Binding queue [%s] to exchange [%s]...", q.Name, e.Name)
	// Args: name, key, exchange string, noWait bool, args Table
	err = w.Channel.QueueBind(q.Name, "", e.Name, false, nil)
	if err != nil {
		return
	}

	return
}

func (w *Worker) Consume() (err error) {
	defer func() {
		w.Connection.Close()
		w.Channel.Close()
		w.Logger.Close()
	}()

	w.handleConnectionCloseError()

	msgs, err := w.listenToQueue()
	if err != nil {
		return
	}
	w.forwardMessages(msgs)

	return
}

func (w *Worker) handleConnectionCloseError() {
	closeErr := make(chan *amqp.Error)
	w.Connection.NotifyClose(closeErr)
	go func() {
		w.Logger.Error("Connection closed: %v", <-closeErr)
		w.Logger.Close()
		os.Exit(1)
	}()
}

func (w *Worker) listenToQueue() (<-chan amqp.Delivery, error) {
	w.Logger.Info("Starting a new consumer...")
	msgs, err := w.Channel.Consume(
		w.Config.Worker.Queue.Name, // queue string
		"",    // consumer string
		false, // autoAck bool
		false, // exclusive bool
		false, // noLocal bool
		false, // noWait bool
		nil,   // args Table
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (w *Worker) forwardMessages(msgs <-chan amqp.Delivery) {
	forever := make(chan bool)

	go func() {
		for m := range msgs {
			if out, err := w.cmd(m.Body).CombinedOutput(); err != nil {
				w.Logger.Error("Failed to process message: %s \n Output: %s", err, out)
				// Sleep for 5 seconds and then retry
				time.Sleep(5 * time.Second)
				m.Nack(true, true)
			} else {
				w.Logger.Info("One message processed")
				m.Ack(true)
			}
		}
	}()

	w.Logger.Info("Waiting for messages...")
	<-forever
}

func (w *Worker) cmd(msg []byte) *exec.Cmd {
	var name string = w.Config.Worker.Script
	var args []string

	if subs := strings.Split(name, " "); len(subs) > 1 {
		name, args = subs[0], subs[1:]
	}

	args = append(args, base64.StdEncoding.EncodeToString(msg))
	return exec.Command(name, args...)
}
