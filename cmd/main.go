package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/djalben/epn-killer-mvp/internal/app"
	"github.com/djalben/epn-killer-mvp/internal/config"
	logHandler "github.com/djalben/epn-killer-mvp/internal/infrastructure/logger/handler"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.Parse()
	if err != nil {
		// на этом этапе ещё нет логгера — используем обычный вывод
		panic("failed to parse config: " + err.Error())
	}

	// Создаём логгер (чисто, без Prometheus)
	handler := logHandler.Create(cfg.LogPlain, cfg.LogLevel)
	logger := slog.New(handler)

	logger.Info("🚀 Starting XPLR...")

	// Создаём контейнер
	container, err := app.NewContainer(&cfg)
	if err != nil {
		logger.Error("failed to create container", "error", err)
		os.Exit(1)
	}

	defer func() {
		err := container.Close()
		if err != nil {
			logger.Error("failed to close container", "error", err)
		}
	}()

	logger.Info("XPLR started successfully",
		"host", cfg.ServerHost,
		"port", cfg.ServerPort,
	)

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	<-ctx.Done()
	logger.Info("Shutting down...")
}
