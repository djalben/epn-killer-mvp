package logger

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
)

const levelLabelName = "level"

type prometheusHandler struct {
	slog.Handler

	handleCount *prometheus.CounterVec
}

// Prometheus - прометеус в логгере.
//
//	slogHandler := handler.NewHandler(os.Stdout, &slog.HandlerOptions{Level: level.GetLevel(lvl)})
//	prometheusHandler, err := Prometheus(slogHandler)
//	logger := slog.New(prometheusHandler)
func Prometheus(handler slog.Handler) (slog.Handler, error) {
	pHandler := prometheusHandler{
		Handler: handler,
		handleCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "logger_writes_total",
				Help: "Число записанных логгером событий",
			},
			[]string{levelLabelName},
		),
	}

	err := prometheus.Register(pHandler.handleCount)
	if err != nil {
		return nil, fmt.Errorf("ошибка регистрации метрик: error = %w", err)
	}

	return pHandler, nil
}

func (ph prometheusHandler) Handle(ctx context.Context, record slog.Record) error {
	labels := map[string]string{levelLabelName: record.Level.String()}
	ph.handleCount.With(labels).Inc()

	err := ph.Handler.Handle(ctx, record)
	if err != nil {
		return fmt.Errorf("ошибка при вызове ph.handler.Handle: error = %w", err)
	}

	return nil
}

func (ph prometheusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return prometheusHandler{
		Handler:     ph.Handler.WithAttrs(attrs),
		handleCount: ph.handleCount,
	}
}

func (ph prometheusHandler) WithGroup(name string) slog.Handler {
	return prometheusHandler{
		Handler:     ph.Handler.WithGroup(name),
		handleCount: ph.handleCount,
	}
}
