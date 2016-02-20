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

type Rex struct {
	Config     *Config
	Logger     *Logger
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func NewRex(c *Config, l *Logger) (r *Rex, err error) {
	r = &Rex{Config: c, Logger: l}
	if err = r.connect(); err != nil {
		return
	}
	if err = r.createQueueAndExchange(); err != nil {
		return
	}
	return
}

func (r *Rex) connect() (err error) {
	c := r.Config.Connection
	p := r.Config.Consumer.Prefetch

	r.Logger.Infof("Connecting to RabbitMQ server...")
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

	r.Logger.Infof("Opening channel...")
	chn, err := conn.Channel()
	if err != nil {
		return
	}

	r.Logger.Infof("Setting QoS... ")
	if p.Count == 0 {
		p.Count = 3
	}
	// Args: prefetchCount, prefetchSize int, global bool
	if err = chn.Qos(p.Count, 0, p.Global); err != nil {
		return
	}

	r.Connection = conn
	r.Channel = chn

	return
}

func (r *Rex) createQueueAndExchange() (err error) {
	q := r.Config.Consumer.Queue
	e := r.Config.Consumer.Exchange

	// Create queue
	r.Logger.Infof("Declaring queue [%s]...", q.Name)
	// Args: name string, durable, autoDelete, exclusive, noWait bool, args Table
	_, err = r.Channel.QueueDeclare(q.Name, q.Durable, q.AutoDelete, false, false, nil)
	if err != nil {
		return
	}

	// Create exchange
	if e.Name == "" {
		r.Logger.Infof("Empty Exchange name - use default exchange.")
		return
	}
	r.Logger.Infof("Declaring exchange [%s]...", e.Name)
	if e.Type == "" {
		e.Type = "direct"
	}
	// Args: name, kind string, durable, autoDelete, internal, noWait bool, args Table
	err = r.Channel.ExchangeDeclare(e.Name, e.Type, e.Durable, e.AutoDelete, false, false, nil)
	if err != nil {
		return
	}

	// Bind queue and exchange
	r.Logger.Infof("Binding queue [%s] to exchange [%s]...", q.Name, e.Name)
	// Args: name, key, exchange string, noWait bool, args Table
	err = r.Channel.QueueBind(q.Name, q.RoutingKey, e.Name, false, nil)
	if err != nil {
		return
	}

	return
}

func (r *Rex) Consume() (err error) {
	r.handleConnectionCloseError()

	msgs, err := r.listenToQueue()
	if err != nil {
		return
	}
	r.forwardMessages(msgs)

	return
}

func (r *Rex) handleConnectionCloseError() {
	closeErr := make(chan *amqp.Error)
	r.Connection.NotifyClose(closeErr)
	go func() {
		r.Logger.Errorf("Connection closed: %v", <-closeErr)
		r.Close()
		os.Exit(1)
	}()
}

func (r *Rex) listenToQueue() (<-chan amqp.Delivery, error) {
	r.Logger.Infof("Starting a new consumer...")
	msgs, err := r.Channel.Consume(
		r.Config.Consumer.Queue.Name, // queue string
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

func (r *Rex) forwardMessages(msgs <-chan amqp.Delivery) {
	forever := make(chan bool)

	go func() {
		for m := range msgs {
			r.Logger.Infof("New message came in. Start processing...")
			if out, err := r.cmd(m.Body).CombinedOutput(); err != nil {
				r.Logger.Errorf("Failed to process message: %s \n Output: %s", err, out)
				// Sleep for 5 seconds and then retry
				time.Sleep(5 * time.Second)
				m.Nack(true, true)
			} else {
				r.Logger.Infof("[Message Processed]")
				m.Ack(true)
			}
		}
	}()

	r.Logger.Infof("Waiting for messages...")
	<-forever
}

func (r *Rex) cmd(msg []byte) *exec.Cmd {
	var name string = r.Config.Consumer.Script
	var args []string

	if subs := strings.Split(name, " "); len(subs) > 1 {
		name, args = subs[0], subs[1:]
	}

	args = append(args, base64.StdEncoding.EncodeToString(msg))
	return exec.Command(name, args...)
}

func (r *Rex) Close() {
	r.Connection.Close()
	r.Channel.Close()
	r.Logger.Close()
}
