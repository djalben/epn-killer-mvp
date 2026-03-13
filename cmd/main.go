package main

import (
	"context"
	"os/signal"
	"syscall"
)

func main() {
	_, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGUSR1,
		syscall.SIGTERM,
	)
	defer cancel()
}
