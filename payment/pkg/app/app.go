package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	paymentAPI "github.com/Ilya96s/rocket-factory-new/payment/internal/api/paymnet/v1"
	paymentService "github.com/Ilya96s/rocket-factory-new/payment/internal/service/payment"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	// Адрес сервера
	grpcAddress = ":50052"

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

	// Сервисный слой
	service := paymentService.New()

	// gRPC API
	api := paymentAPI.New(service)

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

	// Связываем реализацию API со сгенерированный gRPC-сервером
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)

	reflection.Register(grpcServer)

	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("payment service запущен",
			"address", grpcAddress,
		)

		serverErrors <- grpcServer.Serve(listener)
	}()

	select {
	case <-ctx.Done():
		slog.Info("получен сигнал остановки payment service")
	case serverErr := <-serverErrors:
		if serverErr != nil {
			return fmt.Errorf("работа gRPC-сервера %w", serverErr)
		}
		return nil
	}

	// Перестаем принимать новые запросы и ждем пока текущие запросы закончат работу
	grpcServer.GracefulStop()

	slog.Info("payment service остановлен")
	return nil

}
