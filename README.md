# Rex

[![Build Status](https://travis-ci.org/beanworks/rex.svg)](https://travis-ci.org/beanworks/rex)
[![Go Report Card](https://goreportcard.com/badge/github.com/beanworks/rex)](https://goreportcard.com/report/github.com/beanworks/rex)

Rex is a command line based RabbitMQ message consuming and distribution app.

<img align="left" src="https://cloud.githubusercontent.com/assets/965430/13627896/af1cf03a-e582-11e5-8de9-ca4b62665a10.jpg">

Rex connects to a RabbitMQ server, and listens to a queue for messages. When a new message arrives,
Rex forwards the message to another CLI software specified in config, and waits for execution result.
If a zero value returned, Rex will acknowledges the delivery, or negatively acknowledges the delivery
if a non-zero value returned.

Rex is designed as a long running service suited with multithread, concurrency and cross-platform
support. It's a good enhancement for your current RabbitMQ consumer/worker/subscriber written in
languages like PHP, JavaScript/Node.js, Python, etc.
