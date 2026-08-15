package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	partAPI "github.com/Ilya96s/rocket-factory-new/inventory/internal/api/inventory/v1"
	partRepository "github.com/Ilya96s/rocket-factory-new/inventory/internal/repository/part"
	partService "github.com/Ilya96s/rocket-factory-new/inventory/internal/service/part"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	grpcAddress = ":50051"

	// gRPC keepalive параметры
	grpcMaxConnectionIdle     = 5 * time.Minute
	grpcMaxConnectionAge      = 30 * time.Minute
	grpcMaxConnectionAgeGrace = 5 * time.Minute

	grpcKeepaliveTime    = 2 * time.Hour
	grpcKeepaliveTimeout = 20 * time.Second

	grpcMinPingInterval = 30 * time.Second

	inventoryDSN = "postgres://inventory-service-user:inventory-service-password@localhost:5433/inventory-service?sslmode=disable"
)

func Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("создать TCP listener на %q: %w", grpcAddress, err)
	}
	defer listener.Close()

	pool, err := newStorage(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	repository := partRepository.New(pool)
	service := partService.New(repository)
	api := partAPI.New(service)

	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     grpcMaxConnectionIdle,
			MaxConnectionAge:      grpcMaxConnectionAge,
			MaxConnectionAgeGrace: grpcMaxConnectionAgeGrace,
			Time:                  grpcKeepaliveTime,
			Timeout:               grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcMinPingInterval,
			PermitWithoutStream: true,
		}),
	)

	inventoryv1.RegisterInventoryServiceServer(grpcServer, api)

	reflection.Register(grpcServer)

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("inventory service запущен", "address", grpcAddress)
		serverErrors <- grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		slog.Info("получен сигнал остановки inventory service")
	case serverErr := <-serverErrors:
		if serverErr != nil {
			return fmt.Errorf("работа gRPC-сервера %w", serverErr)
		}
		return nil
	}

	grpcServer.GracefulStop()

	slog.Info("inventory service остановлен")
	return nil
}

func newStorage(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, inventoryDSN)
	if err != nil {
		return nil, fmt.Errorf("создать пул соединений: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("проверить соединение с БД: %w", err)
	}

	return pool, nil
}
