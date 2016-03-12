package main

import (
	"runtime"

	"github.com/beanworks/rex/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
