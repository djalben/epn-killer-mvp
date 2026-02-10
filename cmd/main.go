package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	logHandler "github.com/djalben/epn-killer-mvp/internal/lib/logger/handler"
	libEnvironment "github.com/djalben/epn-killer-mvp/internal/lib/tools/environment"
	"github.com/djalben/epn-killer-mvp/internal/services"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGUSR1,
		syscall.SIGTERM,
	)
	defer cancel()

	logLevel, logPlain := libEnvironment.GetLog()
	handler := logHandler.Create(logPlain, logLevel)
	logger := slog.New(handler)

	err := services.Run(ctx, logger)
	if err != nil {
		logger.ErrorContext(ctx, "services.Run", slog.Any("error", err))

		return
	}
}
