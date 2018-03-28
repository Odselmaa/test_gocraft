package main

import (
	"os"
	"os/signal"
	"github.com/test/lib"
)
// worker
func main() {

	var pool = lib.Pool
	pool.Middleware((*lib.Context).Log)
	pool.Middleware((*lib.Context).FindMessage)
	// Map the name of jobs to handler functions
	pool.Job("test", (*lib.Context).SendMessage)

	// Customize options:
	// Start processing jobs
	pool.Start()
	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	// Stop the pool
	pool.Stop()


}

