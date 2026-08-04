package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Ilya96s/rocket-factory-new/payment/pkg/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Run(ctx); err != nil {
		slog.Error(
			"payment service завершился с ошибкой",
			"error", err,
		)
		os.Exit(1)
	}
}
