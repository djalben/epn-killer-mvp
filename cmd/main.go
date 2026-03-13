package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/djalben/epn-killer-mvp/internal/app"
	"github.com/prometheus/prometheus/config"
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

	cfg := config.Load()
	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()
}
