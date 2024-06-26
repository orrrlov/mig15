package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	c, err := initContainer()
	if err != nil {
		panic(err)
	}

	// Graceful shutdown
	<-done
	err = c.shutdownContainer()
	if err != nil {
		panic(err)
	}
}
