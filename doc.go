/*
Rex is a high performance and concurrent RabbitMQ consumer that forwards messages to another worker script.

Rex connects to a RabbitMQ server, and listens to a queue for messages. When a new message arrives,
Rex forwards the message to another CLI software specified in config, and waits for execution result.
If a zero value returned, Rex will acknowledges the delivery, or negatively acknowledges the delivery
if a non-zero value returned.

Rex is designed as a long running service suited with multithread, concurrency and cross-platform
support. It's a good enhancement for your current RabbitMQ consumer/worker/subscriber written in
languages like PHP, JavaScript/Node.js, Ruby, Python, etc.

Example Usage

	$ rex up -c /path/to/config.yml

Print Help Menu

	$ rex help

Check Rex Version

	$ rex version

Example Config

	connection:
	  host: localhost
	  username: username
	  password: password
	  vhost: rex
	  port: 5672
	consumer:
	  exchange:
	    name: exchange.name
	    type: direct
	    durable: true
	    auto_delete: false
	  prefetch:
	    count: 10
	    global: false
	  queue:
	    name: message.queue.name
	    routing_key: message.queue.routing.key
	    durable: true
	    auto_delete: false
	  worker:
	    count: 10
	    script: /path/to/your/worker/script
	    retry_interval: 30
	logger:
	  output: both
	  formatter: text
	  level: debug
	  log_file: ./rex.log
*/
package main

// blank imports help docs.
import (
	// cmd package
	_ "github.com/beanworks/rex/cmd"
	// rabbit package
	_ "github.com/beanworks/rex/rabbit"
)
