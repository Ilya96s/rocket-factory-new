package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	orderHandler "github.com/Ilya96s/rocket-factory-new/order/pkg/handler"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
)

const (
	orderHTTPAddress        = ":8080"
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second

	shutdownTimeout = 10 * time.Second
)

func main() {
	rootCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	clientKeepalive := keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}

	inventoryConnection, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientKeepalive),
	)
	if err != nil {
		slog.Error(
			"не удалось создать соединение с inventory service",
			"error", err,
		)
		os.Exit(1)
	}
	defer func() {
		if closeErr := inventoryConnection.Close(); closeErr != nil {
			slog.Error(
				"не удалось закрыть соединение с inventory service",
				"error", closeErr,
			)
		}
	}()

	paymentConnection, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientKeepalive),
	)
	if err != nil {
		slog.Error(
			"не удалось создать соединение с payment service",
			"error", err,
		)
		os.Exit(1)
	}
	defer func() {
		if closeErr := paymentConnection.Close(); closeErr != nil {
			slog.Error(
				"не удалось закрыть соединение с payment service",
				"error", closeErr,
			)
		}
	}()

	store := orderHandler.NewOrderStore()

	handler := orderHandler.NewHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConnection),
		paymentv1.NewPaymentServiceClient(paymentConnection),
		store,
	)

	orderServer, err := orderHandler.SetupServer(handler)
	if err != nil {
		slog.Error(
			"не удалось создать ogen server",
			"error", err,
		)
		os.Exit(1)
	}

	httpServer := &http.Server{
		Addr:              orderHTTPAddress,
		Handler:           orderServer,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info(
			"order service запущен",
			"address", orderHTTPAddress,
		)

		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
			return
		}

		serverErrors <- nil
	}()

	select {
	case <-rootCtx.Done():
		slog.Info("получен сигнал остановки order service")

	case serverErr := <-serverErrors:
		if serverErr != nil {
			slog.Error(
				"ошибка HTTP-сервера",
				"error", serverErr,
			)
		}
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		shutdownTimeout,
	)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"не удалось корректно остановить HTTP-сервер",
			"error", err,
		)

		if closeErr := httpServer.Close(); closeErr != nil {
			slog.Error(
				"не удалось принудительно остановить HTTP-сервер",
				"error", closeErr,
			)
		}
	}

	slog.Info("order service остановлен")
}
