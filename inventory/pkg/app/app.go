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
)

func Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		return fmt.Errorf("создать TCP listener на %q: %w", grpcAddress, err)
	}
	defer listener.Close()

	repository := partRepository.New()
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
