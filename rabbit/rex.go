package rabbit

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"runtime"
	"time"

	"github.com/streadway/amqp"
)

const (
	DefaultPrefetchCount = 10
	DefaultRetryInterval = 30
)

// Rex represents the RabbitMQ consumer. It has methods connect to RabbitMQ, create exchange,
// queue, bind queue to exchange with routing key, listen to channel for incoming messages,
// forward message body to worker scripts, and ACK/NACK deliveries based on the worker script's
// return code.
type Rex struct {
	Amqp    AmqpConsumer
	Config  *Config
	Logger  *Logger
	Script  ScriptCaller
	Forever chan bool
}

// NewRex returns a configured Rex instance. It will return a non nil error if anything goes
// wrong when connecting to RabbitMQ, or creating queue and exchange.
func NewRex(c *Config, l *Logger, a AmqpConsumer, s ScriptCaller) (r *Rex, err error) {
	r = &Rex{Amqp: a, Config: c, Logger: l, Script: s, Forever: make(chan bool)}
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

	r.Logger.Infof("Connecting to RabbitMQ server [%s:%d]...", c.Host, c.Port)
	_, err = r.Amqp.Dial(fmt.Sprintf(
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
	_, err = r.Amqp.Channel()
	if err != nil {
		return
	}

	if p.Count == 0 {
		p.Count = DefaultPrefetchCount
	}
	r.Logger.Infof("Setting QoS [prefetch: %d]...", p.Count)
	// Args: prefetchCount, prefetchSize int, global bool
	if err = r.Amqp.Qos(p.Count, 0, p.Global); err != nil {
		return
	}

	return
}

func (r *Rex) createQueueAndExchange() (err error) {
	q := r.Config.Consumer.Queue
	e := r.Config.Consumer.Exchange

	// Create queue
	r.Logger.Infof("Declaring queue [%s]...", q.Name)
	// Args: name string, durable, autoDelete, exclusive, noWait bool, args Table
	_, err = r.Amqp.QueueDeclare(q.Name, q.Durable, q.AutoDelete, false, false, nil)
	if err != nil {
		return
	}

	// Create exchange
	if e.Name == "" {
		r.Logger.Infof("Empty Exchange name - use default exchange.")
		// The default exchange is implicitly bound to every queue,
		// with a routing key equal to the queue name. It is not possible
		// to explicitly bind to, or unbind from the default exchange.
		// It also cannot be deleted.
		return
	}
	r.Logger.Infof("Declaring exchange [%s]...", e.Name)
	if e.Type == "" {
		e.Type = "direct"
	}
	// Args: name, kind string, durable, autoDelete, internal, noWait bool, args Table
	err = r.Amqp.ExchangeDeclare(e.Name, e.Type, e.Durable, e.AutoDelete, false, false, nil)
	if err != nil {
		return
	}

	// Bind queue and exchange
	r.Logger.Infof("Binding queue [%s] to exchange [%s]...", q.Name, e.Name)
	// Args: name, key, exchange string, noWait bool, args Table
	err = r.Amqp.QueueBind(q.Name, q.RoutingKey, e.Name, false, nil)
	if err != nil {
		return
	}

	return
}

type NotifyCloseCallback func()

// NotifyClose registers a listener for when the server sends a channel or connection
// exception in the form of a Connection.Close or Channel.Close method. Connection
// exceptions will be broadcast to all open channels and all channels will be closed,
// where channel exceptions will only be broadcast to listeners to this channel.
func (r *Rex) NotifyClose(fn NotifyCloseCallback) {
	err := make(chan *amqp.Error, 1)
	r.Amqp.NotifyClose(err)
	go func() {
		r.Logger.Errorf("Connection closed: %v", <-err)
		r.Close()
		fn()
	}()
}

// Close closes RabbitMQ connection and channel. It also closes Logger.
func (r *Rex) Close() {
	r.Amqp.Close()
	r.Logger.Close()
}

// Consume listens to RabbitMQ channel, forward each message's body to a worker script.
// Based on return code of the worker script, it acknowledge the delivery when the code
// is 0, and negatively acknowledge the delivery when the code is 1.
func (r *Rex) Consume() (err error) {
	msgs, err := r.listenToQueue()
	if err != nil {
		return
	}
	r.handleMessages(msgs, r.forwardToWorker)

	return
}

func (r *Rex) listenToQueue() (<-chan amqp.Delivery, error) {
	r.Logger.Infof("Starting a new consumer...")
	msgs, err := r.Amqp.Consume(
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

type MessageForwarder func(msgs <-chan amqp.Delivery, retryInterval int)

func (r *Rex) handleMessages(msgs <-chan amqp.Delivery, forwarder MessageForwarder) {
	w := r.Config.Consumer.Worker

	retryInterval := w.RetryInterval
	if retryInterval == 0 {
		retryInterval = DefaultRetryInterval
	}
	workerCount := w.Count
	if workerCount == 0 {
		workerCount = runtime.NumCPU()
	}

	r.Logger.Infof("Waiting for messages...")
	for i := 0; i < workerCount; i++ {
		go forwarder(msgs, retryInterval)
	}
	<-r.Forever
}

func (r Rex) forwardToWorker(msgs <-chan amqp.Delivery, retryInterval int) {
	for m := range msgs {
		md5sum := md5.Sum(m.Body)
		r.Logger.Infof("New message came in. md5sum: %x", md5sum)
		r.Logger.Debugf(
			"Message md5sum: %x; Redelivered: %v; Body: %v",
			md5sum, m.Redelivered, string(m.Body),
		)
		out, err := r.Script.ExecWith(m.Body)
		if err != nil {
			r.Logger.Errorf(
				"Failed to process message. md5sum: %x; Error: %s; Output: %s",
				md5sum, err, out,
			)
			// Sleep and then retry
			time.Sleep(time.Duration(retryInterval) * time.Second)
			m.Nack(false, true)
			r.Logger.Debugf("Message Nack'ed and will be redelivered. md5sum: %x", md5sum)
		} else {
			r.Logger.Debugf(
				"Command successfully executed with. md5sum: %x; Output: %s",
				md5sum, out,
			)
			m.Ack(false)
			r.Logger.Infof("Message Processed. md5sum: %x", md5sum)
		}
	}
}
