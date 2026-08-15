package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	orderAPI "github.com/Ilya96s/rocket-factory-new/order/internal/api/order/v1"
	inventoryClient "github.com/Ilya96s/rocket-factory-new/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Ilya96s/rocket-factory-new/order/internal/client/grpc/payment/v1"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/order"
	"github.com/Ilya96s/rocket-factory-new/order/internal/repository/order_item"
	orderService "github.com/Ilya96s/rocket-factory-new/order/internal/service/order"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

const (
	orderHTTPAddress        = ":8080"
	inventoryServiceAddress = "localhost:50051"
	paymentServiceAddress   = "localhost:50052"
	// Конфигурация БД пока задаётся в коде.
	orderDSN = "postgres://order-service-user:order-service-password@localhost:5432/order-service?sslmode=disable"

	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 10 * time.Second
	idleTimeout       = 60 * time.Second

	shutdownTimeout = 10 * time.Second
)

// Run - собирает зависимости, запускает HTTP-сервер
// и останавливает приложение после отмены контекста
func Run(ctx context.Context) error {
	inventoryConnection, paymentConnection, err := newServiceConnections()
	if err != nil {
		return err
	}
	defer closeGRPCConnection(inventoryConnection, "inventory service")
	defer closeGRPCConnection(paymentConnection, "payment service")

	pool, txManager, err := newStorage(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	// Конкретная реализация репозитория
	orderRepo := order.NewRepository(pool)
	orderItemRepo := order_item.NewRepository(pool)

	// Сгенерированный protobuf-клиент
	inventoryProtoClient := inventoryv1.NewInventoryServiceClient(inventoryConnection)

	// Наша обертка, реализующая интерфейс service.InventoryClient
	inventoryGRPCClient := inventoryClient.New(inventoryProtoClient)

	// Сгенерированный protobuf-клиент
	paymentProtoClient := paymentv1.NewPaymentServiceClient(paymentConnection)

	// Наша обертка, реализующая интерфейс service.PaymentClient
	paymentGRPCClient := paymentClient.New(paymentProtoClient)

	// Бизнес-слой
	service := orderService.NewService(orderRepo, orderItemRepo, inventoryGRPCClient, paymentGRPCClient, txManager)

	// HTTP/OpenAPI-адаптер
	apiHandler := orderAPI.New(service)

	orderServer, err := orderAPI.SetupServer(apiHandler)
	if err != nil {
		return fmt.Errorf("создать ogen server: %w", err)
	}

	httpServer := &http.Server{
		Addr:              orderHTTPAddress,
		Handler:           orderServer,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("order service запущен",
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
	case <-ctx.Done():
		slog.Info("получен сигнал остановки order service")
	case serverErr := <-serverErrors:
		if serverErr != nil {
			return fmt.Errorf("работа HTTP-сервера: %w", serverErr)
		}
		return nil
	}

	shutdownCtx, shutDownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutDownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"не удалось корректно остановить HTTP-сервер",
			"err", err,
		)
		if closeErr := httpServer.Close(); closeErr != nil {
			return fmt.Errorf(
				"принудительно остановить HTTP-сервер: %w",
				closeErr,
			)
		}

		return fmt.Errorf(
			"корректно остановить HTTP-сервер: %w",
			err,
		)
	}

	slog.Info("order service остановлен")

	return nil
}

func newServiceConnections() (*grpc.ClientConn, *grpc.ClientConn, error) {
	clientKeepAlive := keepalive.ClientParameters{
		Time:                30 * time.Second,
		Timeout:             10 * time.Second,
		PermitWithoutStream: true,
	}

	inventoryConnection, err := grpc.NewClient(
		inventoryServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientKeepAlive),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("создать соединение с inventory service: %w", err)
	}

	paymentConnection, err := grpc.NewClient(
		paymentServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(clientKeepAlive),
	)
	if err != nil {
		closeGRPCConnection(inventoryConnection, "inventory service")
		return nil, nil, fmt.Errorf("создать соединение с payment service: %w", err)
	}

	return inventoryConnection, paymentConnection, nil
}

func newStorage(ctx context.Context) (*pgxpool.Pool, *manager.Manager, error) {
	pool, err := pgxpool.New(ctx, orderDSN)
	if err != nil {
		return nil, nil, fmt.Errorf("создать пул соединений: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("проверить соединение с БД: %w", err)
	}

	txManager, err := manager.New(trmpgx.NewDefaultFactory(pool))
	if err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("создать менеджер транзакций: %w", err)
	}

	return pool, txManager, nil
}

func closeGRPCConnection(connection *grpc.ClientConn, serviceName string) {
	if err := connection.Close(); err != nil {
		slog.Error(
			"не удалось закрыть gRPC-соединение",
			"service", serviceName,
			"error", err,
		)
	}
}
