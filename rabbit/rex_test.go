package rabbit

import (
	"errors"
	"runtime"
	"sync"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type fakeAmqp struct {
	Chan AmqpChannelCloser
	Conn AmqpConnectionCloser

	AmqpConsumer
}

func (a *fakeAmqp) Dial(url string) (*amqp.Connection, error) {
	return &amqp.Connection{}, nil
}

func (a *fakeAmqp) Channel() (*amqp.Channel, error) {
	return &amqp.Channel{}, nil
}

func (a *fakeAmqp) Qos(prefetchCount, prefetchSize int, global bool) error {
	return nil
}

func (a *fakeAmqp) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	return amqp.Queue{}, nil
}

func (a *fakeAmqp) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return nil
}

func (a *fakeAmqp) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return nil
}

func (a *fakeAmqp) NotifyClose(c chan *amqp.Error) chan *amqp.Error {
	c <- &amqp.Error{Code: 100, Reason: "Fake AMQP Error"}
	return c
}

func (a *fakeAmqp) Close() (err error) {
	return nil
}

func (a *fakeAmqp) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	d := make(chan amqp.Delivery)
	return d, nil
}

var (
	errorBadAmqpDial            = errors.New("Error while dialing.")
	errorBadAmqpChannel         = errors.New("Error connecting to channel.")
	errorBadAmqpQos             = errors.New("Error setting up QoS.")
	errorBadAmqpQueueDeclare    = errors.New("Error declaring queue.")
	errorBadAmqpExchangeDeclare = errors.New("Error declaring exchange.")
	errorBadAmqpQueueBind       = errors.New("Error binding queue.")
	errorBadAmqpConsume         = errors.New("Error consuming messages.")
	errorBadScriptExecWith      = errors.New("Error executing script.")
)

type fakeAmqpBadDial struct{ *fakeAmqp }

func (a *fakeAmqpBadDial) Dial(url string) (*amqp.Connection, error) {
	return nil, errorBadAmqpDial
}

type fakeAmqpBadChannel struct{ *fakeAmqp }

func (a *fakeAmqpBadChannel) Channel() (*amqp.Channel, error) {
	return nil, errorBadAmqpChannel
}

type fakeAmqpBadQos struct{ *fakeAmqp }

func (a *fakeAmqpBadQos) Qos(prefetchCount, prefetchSize int, global bool) error {
	return errorBadAmqpQos
}

type fakeAmqpBadQueueDeclare struct{ *fakeAmqp }

func (a *fakeAmqpBadQueueDeclare) QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	var q amqp.Queue
	return q, errorBadAmqpQueueDeclare
}

type fakeAmqpBadExchangeDeclare struct{ *fakeAmqp }

func (a *fakeAmqpBadExchangeDeclare) ExchangeDeclare(name, kind string, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return errorBadAmqpExchangeDeclare
}

type fakeAmqpBadQueueBind struct{ *fakeAmqp }

func (a *fakeAmqpBadQueueBind) QueueBind(name, key, exchange string, noWait bool, args amqp.Table) error {
	return errorBadAmqpQueueBind
}

type fakeAmqpBadConsume struct{ *fakeAmqp }

func (a *fakeAmqpBadConsume) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	return nil, errorBadAmqpConsume
}

type fakeAcknowledger struct{ mock.Mock }

func (a *fakeAcknowledger) Ack(tag uint64, multiple bool) error {
	args := a.Called(tag, multiple)
	return args.Error(0)
}

func (a *fakeAcknowledger) Nack(tag uint64, multiple, requeue bool) error {
	args := a.Called(tag, multiple, requeue)
	return args.Error(0)
}

func (a *fakeAcknowledger) Reject(tag uint64, requeue bool) error {
	args := a.Called(tag, requeue)
	return args.Error(0)
}

type fakeScript struct{ ScriptCaller }

func (s fakeScript) ExecWith(msg []byte) ([]byte, error) {
	return []byte("Success"), nil
}

type fakeScriptBadExecWith struct{ fakeScript }

func (s fakeScriptBadExecWith) ExecWith(msg []byte) ([]byte, error) {
	return nil, errorBadScriptExecWith
}

func TestConnect(t *testing.T) {
	var (
		err error
		rex *Rex
	)

	cfg := &Config{}
	log := &Logger{}
	scr := &fakeScript{}

	done := make(chan bool)

	rex = &Rex{&fakeAmqpBadDial{}, cfg, log, scr, done}
	err = rex.connect()
	assert.Error(t, err, errorBadAmqpDial.Error())

	rex = &Rex{&fakeAmqpBadChannel{}, cfg, log, scr, done}
	err = rex.connect()
	assert.Error(t, err, errorBadAmqpChannel.Error())

	rex = &Rex{&fakeAmqpBadQos{}, cfg, log, scr, done}
	err = rex.connect()
	assert.Error(t, err, errorBadAmqpQos.Error())

	rex = &Rex{&fakeAmqp{}, cfg, log, scr, done}
	err = rex.connect()
	assert.NoError(t, err)
}

func TestCreateQueueAndExchange(t *testing.T) {
	var (
		err error
		rex *Rex
	)

	cfg := &Config{}
	log := &Logger{}
	scr := &fakeScript{}

	done := make(chan bool)

	rex = &Rex{&fakeAmqpBadQueueDeclare{}, cfg, log, scr, done}
	err = rex.createQueueAndExchange()
	assert.Error(t, err, errorBadAmqpQueueDeclare.Error())

	// Test using default exchange
	cfg.Consumer.Exchange.Name = ""
	rex = &Rex{&fakeAmqpBadExchangeDeclare{}, cfg, log, scr, done}
	err = rex.createQueueAndExchange()
	assert.NoError(t, err)

	// Test with a custom exchange name
	cfg.Consumer.Exchange.Name = "FakeExchange"
	rex = &Rex{&fakeAmqpBadExchangeDeclare{}, cfg, log, scr, done}
	err = rex.createQueueAndExchange()
	assert.Error(t, err, errorBadAmqpExchangeDeclare.Error())

	rex = &Rex{&fakeAmqpBadQueueBind{}, cfg, log, scr, done}
	err = rex.createQueueAndExchange()
	assert.Error(t, err, errorBadAmqpQueueBind.Error())

	rex = &Rex{&fakeAmqp{}, cfg, log, scr, done}
	err = rex.createQueueAndExchange()
	assert.NoError(t, err)
}

func TestNewRex(t *testing.T) {
	var err error

	cfg := &Config{}
	log := &Logger{}
	scr := &fakeScript{}

	_, err = NewRex(cfg, log, &fakeAmqpBadDial{}, scr)
	assert.Error(t, err, errorBadAmqpDial.Error())

	_, err = NewRex(cfg, log, &fakeAmqpBadQueueDeclare{}, scr)
	assert.Error(t, err, errorBadAmqpQueueDeclare.Error())

	_, err = NewRex(cfg, log, &fakeAmqp{}, scr)
	assert.NoError(t, err)
}

func TestNotifyClose(t *testing.T) {
	r, _ := NewRex(&Config{}, &Logger{}, &fakeAmqp{}, &fakeScript{})

	called := false
	done := make(chan bool)
	r.NotifyClose(func() {
		called = true
		done <- true
	})
	<-done

	assert.True(t, called)
}

func TestListenToQueue(t *testing.T) {
	var (
		err error
		rex *Rex
	)

	cfg := &Config{}
	log := &Logger{}
	scr := &fakeScript{}

	rex, _ = NewRex(cfg, log, &fakeAmqpBadConsume{}, scr)
	_, err = rex.listenToQueue()
	assert.Error(t, err, errorBadAmqpConsume.Error())

	rex, _ = NewRex(cfg, log, &fakeAmqp{}, scr)
	_, err = rex.listenToQueue()
	assert.NoError(t, err)
}

func TestHandleMessages(t *testing.T) {
	var (
		rex *Rex
		wg  sync.WaitGroup
		mu  sync.Mutex

		workerCount         int
		workerRetryInterval int
	)

	cfg := &Config{}
	log := &Logger{}
	scr := &fakeScript{}
	chn := make(chan amqp.Delivery)

	rex, _ = NewRex(cfg, log, &fakeAmqp{}, scr)
	close(rex.Forever)

	//
	// Test with default worker count and retry interval
	//
	workerCount = 0
	workerRetryInterval = 0

	wg.Add(runtime.NumCPU())
	rex.handleMessages(chn, func(msgs <-chan amqp.Delivery, retryInterval int) {
		mu.Lock()
		workerCount++
		workerRetryInterval = retryInterval
		mu.Unlock()
		wg.Done()
	})
	wg.Wait()

	assert.Equal(t, DefaultRetryInterval, workerRetryInterval)
	assert.Equal(t, runtime.NumCPU(), workerCount)

	//
	// Test with custom worker count and retry interval
	//
	cfg.Consumer.Worker.RetryInterval = 50
	cfg.Consumer.Worker.Count = 2

	workerCount = 0
	workerRetryInterval = 0

	wg.Add(2)
	rex.handleMessages(chn, func(msgs <-chan amqp.Delivery, retryInterval int) {
		mu.Lock()
		workerCount++
		workerRetryInterval = retryInterval
		mu.Unlock()
		wg.Done()
	})
	wg.Wait()

	assert.Equal(t, 50, workerRetryInterval)
	assert.Equal(t, 2, workerCount)
}

func TestForwardToWorker(t *testing.T) {
	var (
		chn chan amqp.Delivery
		rex *Rex
	)

	cfg := &Config{}
	log := &Logger{}
	ack := new(fakeAcknowledger)
	msg := amqp.Delivery{Acknowledger: ack}

	ack.On("Ack", uint64(0), false).Return(nil)
	ack.On("Nack", uint64(0), false, true).Return(nil)

	rex, _ = NewRex(cfg, log, &fakeAmqp{}, &fakeScriptBadExecWith{})
	chn = make(chan amqp.Delivery, 1)
	chn <- msg
	close(chn)
	rex.forwardToWorker(chn, 0)

	ack.AssertNotCalled(t, "Ack")
	ack.AssertCalled(t, "Nack", uint64(0), false, true)

	rex, _ = NewRex(cfg, log, &fakeAmqp{}, &fakeScript{})
	chn = make(chan amqp.Delivery, 1)
	chn <- msg
	close(chn)
	rex.forwardToWorker(chn, 0)

	ack.AssertCalled(t, "Ack", uint64(0), false)
	ack.AssertNotCalled(t, "Nack")
}
