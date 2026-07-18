package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ilya96s/rocket-factory-new/order/pkg/handler"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	inventoryGrpcAddress = "localhost:50051"
	paymentGrpcAddress   = "localhost:50052"
	orderHttpAddress     = ":8080"

	ReadHeaderTimeout = 5 * time.Second
	ReadTimeout       = 15 * time.Second
	WriteTimeout      = 15 * time.Second
	IdleTimeout       = 60 * time.Second
)

func main() {
	inventoryConn, err := grpc.NewClient(inventoryGrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к InventoryService", "error", err)
		os.Exit(1)
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(paymentGrpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("не удалось подключиться к PaymentService", "error", err)
		os.Exit(1)
	}
	defer paymentConn.Close()

	store := handler.NewOrderStore()
	h := handler.NewHandler(
		inventoryv1.NewInventoryServiceClient(inventoryConn),
		paymentv1.NewPaymentServiceClient(paymentConn),
		store,
	)

	orderServer, err := handler.SetupServer(h)
	if err != nil {
		slog.Error("ошибка создания сервера OpenAPI", "error", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:              orderHttpAddress,
		Handler:           orderServer,
		ReadHeaderTimeout: ReadHeaderTimeout,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("openAPI order сервер запущен", "port", 8080)
		if err = http.ListenAndServe(orderHttpAddress, orderServer); err != nil {
			slog.Error("ошибка запуска openAPI order сервера", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("остановка openAPI order сервера")
	server.Shutdown(ctx)
	slog.Info("openAPI order сервер остановлен")
}
