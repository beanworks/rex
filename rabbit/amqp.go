package rabbit

import "github.com/streadway/amqp"

type AmqpCloser interface {
	Close() error
}

type AmqpConnection interface {
	Channel() (*amqp.Channel, error)
}

type AmqpConnectionCloser interface {
	AmqpConnection
	AmqpCloser
}

type AmqpChannel interface {
	Qos(prefetchCount, prefetchSize int, global bool) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error
	QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error
	NotifyClose(c chan *amqp.Error) chan *amqp.Error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

type AmqpChannelCloser interface {
	AmqpChannel
	AmqpCloser
}

type AmqpConsumer interface {
	Dial(url string) (*amqp.Connection, error)

	AmqpConnection
	AmqpChannel
	AmqpCloser
}

type Amqp struct {
	Chan AmqpChannelCloser
	Conn AmqpConnectionCloser
}

func (a *Amqp) Dial(url string) (conn *amqp.Connection, err error) {
	conn, err = amqp.Dial(url)
	if err != nil {
		return
	}
	a.Conn = AmqpConnectionCloser(conn)
	return
}

func (a *Amqp) Channel() (chn *amqp.Channel, err error) {
	chn, err = a.Conn.Channel()
	if err != nil {
		return
	}
	a.Chan = AmqpChannelCloser(chn)
	return
}

func (a *Amqp) Qos(prefetchCount, prefetchSize int, global bool) error {
	return a.Chan.Qos(prefetchCount, prefetchSize, global)
}

func (a *Amqp) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return a.Chan.QueueDeclare(name, durable, autoDelete, exclusive, noWait, args)
}

func (a *Amqp) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return a.Chan.ExchangeDeclare(name, kind, durable, autoDelete, internal, noWait, args)
}

func (a *Amqp) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return a.Chan.QueueBind(name, key, exchange, noWait, args)
}

func (a *Amqp) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	return a.Chan.NotifyClose(c)
}

func (a *Amqp) Close() (err error) {
	if err = a.Chan.Close(); err != nil {
		return
	}
	if err = a.Conn.Close(); err != nil {
		return
	}
	return
}

func (a *Amqp) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return a.Chan.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
}
