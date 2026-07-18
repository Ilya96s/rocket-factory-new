package main

import (
	"context"
	"log/slog"
	"net"
	"os/signal"
	"syscall"
	"time"

	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	// Адрес сервера
	grpcAddress = "localhost:50052"

	// gRPC keepalive параметры
	grpcMaxConnectionIdle     = 15 * time.Minute // Закрыть idle-соединения (нет активных RPC)
	grpcMaxConnectionAge      = 30 * time.Minute // Принудительная ротация для балансировки
	grpcMaxConnectionAgeGrace = 5 * time.Second  // Время на завершение активных RPC
	grpcKeepaliveTime         = 5 * time.Minute  // Интервал ping'ов для обнаружения мёртвых соединений
	grpcKeepaliveTimeout      = 1 * time.Second  // Таймаут ожидания pong
	grpcMinPingInterval       = 5 * time.Minute  // Минимальный интервал ping'ов от клиента (защита от DoS)
)

type PaymentServer struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func NewServer() *PaymentServer {
	return &PaymentServer{}
}

// PayOrder Обработка оплаты заказа
func (s *PaymentServer) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	if req.OrderUuid == "" {
		return nil, status.Error(codes.InvalidArgument, "не указан uuid заказа")
	}
	if req.GetPaymentMethod() == paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "не указан метод оплаты")
	}
	if !isValidUuid(req.OrderUuid) {
		return nil, status.Error(codes.InvalidArgument, "невалидный uuid заказа")
	}
	transactionUuid := uuid.New().String()
	slog.Info("оплата прошла успешно", "order_uuid", req.OrderUuid, "transaction_uuid", transactionUuid)

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUuid,
	}, nil
}

// Проверка валидности переданного UUID
func isValidUuid(value string) bool {
	_, err := uuid.Parse(value)
	if err != nil {
		return false
	}
	return true
}

func main() {
	listener, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		slog.Error("ошибка создания TCP-сокета", "error", err)
		return
	}

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

	paymentServer := NewServer()
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentServer)

	reflection.Register(grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		slog.Info("gRPC сервер запущен", "address", grpcAddress)
		if serveErr := grpcServer.Serve(listener); serveErr != nil {
			slog.Error("ошибка запуска сервера", "error", serveErr)
			cancel()
		}
	}()

	<-ctx.Done()
	slog.Info("остановка gRPC сервера")
	grpcServer.GracefulStop()
	slog.Info("сервер остановлен")
}
