# Rex

[![Go Report Card](https://goreportcard.com/badge/github.com/beanworks/rex)](https://goreportcard.com/report/github.com/beanworks/rex)

Rex is a command line based RabbitMQ message consuming and distribution app.

Rex connects to a RabbitMQ server, and listens to a queue for messages. When a new message arrives,
Rex forwards the message to another command specified from config, and waits for the execution result.
If a zero value returned, Rex will acknowledges the delivery, or negatively acknowledges the delivery
if a non-zero value returned.
