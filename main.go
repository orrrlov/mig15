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

	c, err := startContainer()
	if err != nil {
		panic(err)
	}

	// Gracefull shutdown
	<-done
	err = c.shutdownContainer()
	if err != nil {
		panic(err)
	}
}
