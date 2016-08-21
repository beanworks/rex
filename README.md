# Rex

[![Build Status](https://travis-ci.org/beanworks/rex.svg)](https://travis-ci.org/beanworks/rex)
[![Go Report Card](https://goreportcard.com/badge/github.com/beanworks/rex)](https://goreportcard.com/report/github.com/beanworks/rex)
[![GoDoc](https://godoc.org/github.com/beanworks/rex?status.svg)](https://godoc.org/github.com/beanworks/rex)

<img align="left" src="https://cloud.githubusercontent.com/assets/965430/13627896/af1cf03a-e582-11e5-8de9-ca4b62665a10.jpg">

***Not a T-Rex***

Rex is a high performance and concurrent RabbitMQ consumer that forwards messages to another worker script.

Rex connects to a RabbitMQ server, and listens to a queue for messages. When a new message arrives,
Rex forwards the message to another CLI software specified in config, and waits for execution result.
If a zero value returned, Rex will acknowledges the delivery, or negatively acknowledges the delivery
if a non-zero value returned.

Rex is designed as a long running service suited with multithread, concurrency and cross-platform
support. It's a good enhancement for your current RabbitMQ consumer/worker/subscriber written in
languages like PHP, JavaScript/Node.js, Ruby, Python, etc.

## Getting started

[Download](https://github.com/beanworks/rex/releases) precompiled binaries specific for your platform, and go from there:

```bash
wget https://github.com/beanworks/rex/releases/download/v0.1.0/rex_linux_amd64.tar.gz
mkdir rex
tar zxvf rex_linux_amd64.tar.gz -C rex
mv ./rex/rex /usr/local/bin

# create a config file from config.yml.dist and then:
rex up -c /path/to/config.yml

# or place the config to /etc/rex/config.yml and then:
rex up
```

Config file can be placed in `/etc/rex/config.yml`, `$HOME/.rex/config.yml` or specified with a `--config | -c` flag.

## Build rex from source

There are two methods to build rex from source.

**[Go 1.6+](https://golang.org/dl/) is needed**

1. Use `go get` tool

```bash
go get github.com/beanworks/rex
rex up -c /path/to/config.yml
```

2. Build with `make`

```bash
git clone git@github.com:beanworks/rex.git
cd rex
make
mv ./rex/rex /usr/local/bin
rex up -c /path/to/config.yml
```

## Usage

```
Rex rabbit is a command line message queue consumer for RabbitMQ.
Rex pulls messages from a queue, takes a good care of the jobs,
redirects message bodies to other responsible parties.

When Rex is not busy, he also likes to hang out with Octocat.

Usage:
  rex [flags]
  rex [command]

Available Commands:
  help        help for rex
  up          Start hopping a rex rabbit consumer
  version     Show rex version

Flags:
  -c, --config string   config file (default is $HOME/.rex.yml)
  -v, --version         Show rex version

Use "rex [command] --help" for more information about a command.
```

## Contribution Guide

Rex is written in golang with cobra (cli), viper (config), logrus (logging) and amqp (amqp protocol).
It's dependencies managed by godep, unit tested with testify, and release managed by gox (cross compile)
and make (Makefile).

Set up golang development environment locally:

```bash
wget https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz -P /tmp
sudo tar zxf /tmp/go1.7.1.linux-amd64.tar.gz -C /usr/local
mkdir -p $HOME/go
cd $HOME/go
mkdir bin pkg src
# add the following lines to your .bashrc or .zshrc
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
```

Git clone the repo to local and build from source:

```bash
git clone git@github.com:beanworks/rex.git $HOME/go/src/github.com/beanworks/rex
cd $HOME/go/src/github.com/beanworks/rex
# build
make
# test
make test
# test that has more feedback
make vtest
# make coverage report
make cover
# make testing rabbitmq cluster (docker required)
make cluster
# make new release binaries
make release
# run rex that you just built from source
cp config.yml.dist config.yml
./rex up -c config.yml
```
