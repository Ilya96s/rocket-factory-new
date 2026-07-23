package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/payment/v1"
)

const grpcRequestTimeout = 5 * time.Second

var _ orderv1.Handler = (*Handler)(nil)

// Handler реализует HTTP API заказов, сгенерированный ogen.
type Handler struct {
	orderv1.UnimplementedHandler

	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *orderStore
}

// NewHandler создаёт обработчик заказов.
func NewHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *orderStore,
) *Handler {
	return &Handler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт ogen HTTP-сервер.
func SetupServer(handler *Handler) (*orderv1.Server, error) {
	return orderv1.NewServer(handler)
}

// CreateOrder реализует POST /api/v1/orders.
func (h *Handler) CreateOrder(
	ctx context.Context,
	req *orderv1.CreateOrderRequest,
) (orderv1.CreateOrderRes, error) {
	// Наличие required-полей уже проверяет ogen.
	// Дополнительно не разрешаем нулевые UUID.
	if req.HullUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid не может быть нулевым UUID",
		}, nil
	}

	if req.EngineUUID == uuid.Nil {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid не может быть нулевым UUID",
		}, nil
	}

	partUUIDs := make([]string, 0, 4)

	partUUIDs = append(
		partUUIDs,
		req.HullUUID.String(),
		req.EngineUUID.String(),
	)

	shieldUUIDValue, hasShield := req.ShieldUUID.Get()
	if hasShield {
		partUUIDs = append(partUUIDs, shieldUUIDValue.String())
	}

	weaponUUIDValue, hasWeapon := req.WeaponUUID.Get()
	if hasWeapon {
		partUUIDs = append(partUUIDs, weaponUUIDValue.String())
	}

	inventoryCtx, cancel := context.WithTimeout(ctx, grpcRequestTimeout)
	defer cancel()

	listPartsResponse, err := h.inventoryClient.ListParts(
		inventoryCtx,
		&inventoryv1.ListPartsRequest{
			Uuids: partUUIDs,
		},
	)
	if err != nil {
		return mapInventoryErrorToCreateOrderResponse(err), nil
	}

	parts := listPartsResponse.GetParts()

	// Защищаемся от ситуации, когда Inventory вернул не все детали.
	if len(parts) != len(partUUIDs) {
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "одна или несколько деталей не найдены",
		}, nil
	}

	partsByUUID := make(
		map[string]*inventoryv1.Part,
		len(parts),
	)

	for _, part := range parts {
		if part == nil {
			return &orderv1.CreateOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "inventory service вернул пустую деталь",
			}, nil
		}

		partsByUUID[part.GetUuid()] = part
	}

	for _, requestedUUID := range partUUIDs {
		if _, ok := partsByUUID[requestedUUID]; !ok {
			return &orderv1.CreateOrderNotFound{
				Code: http.StatusNotFound,
				Message: fmt.Sprintf(
					"деталь с UUID %s не найдена",
					requestedUUID,
				),
			}, nil
		}
	}

	hull := partsByUUID[req.HullUUID.String()]
	if hull.GetPartType() != inventoryv1.PartType_PART_TYPE_HULL {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "hull_uuid указывает не на корпус",
		}, nil
	}

	engine := partsByUUID[req.EngineUUID.String()]
	if engine.GetPartType() != inventoryv1.PartType_PART_TYPE_ENGINE {
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: "engine_uuid указывает не на двигатель",
		}, nil
	}

	if hasShield {
		shield := partsByUUID[shieldUUIDValue.String()]

		if shield.GetPartType() != inventoryv1.PartType_PART_TYPE_SHIELD {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "shield_uuid указывает не на щит",
			}, nil
		}
	}

	if hasWeapon {
		weapon := partsByUUID[weaponUUIDValue.String()]

		if weapon.GetPartType() != inventoryv1.PartType_PART_TYPE_WEAPON {
			return &orderv1.CreateOrderBadRequest{
				Code:    http.StatusBadRequest,
				Message: "weapon_uuid указывает не на оружие",
			}, nil
		}
	}

	var totalPrice int64

	for _, part := range parts {
		if part.GetStockQuantity() <= 0 {
			return &orderv1.CreateOrderConflict{
				Code: http.StatusConflict,
				Message: fmt.Sprintf(
					"деталь %s отсутствует на складе",
					part.GetUuid(),
				),
			}, nil
		}

		totalPrice += part.GetPrice()
	}

	var shieldUUID *uuid.UUID
	if hasShield {
		value := shieldUUIDValue
		shieldUUID = &value
	}

	var weaponUUID *uuid.UUID
	if hasWeapon {
		value := weaponUUIDValue
		weaponUUID = &value
	}

	order := Order{
		OrderUUID:       uuid.New(),
		HullUUID:        req.HullUUID,
		EngineUUID:      req.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      totalPrice,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		Status:          OrderStatusPendingPayment,
		CreatedAt:       time.Now().UTC(),
	}

	h.store.Save(order)

	return &orderv1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}

// GetOrder реализует GET /api/v1/orders/{order_uuid}.
func (h *Handler) GetOrder(
	_ context.Context,
	params orderv1.GetOrderParams,
) (orderv1.GetOrderRes, error) {
	order, ok := h.store.Get(params.OrderUUID)
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	return orderToDTO(order), nil
}

// PayOrder реализует POST /api/v1/orders/{order_uuid}/pay.
func (h *Handler) PayOrder(
	ctx context.Context,
	req *orderv1.PayOrderRequest,
	params orderv1.PayOrderParams,
) (orderv1.PayOrderRes, error) {
	order, ok := h.store.Get(params.OrderUUID)
	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if order.Status != OrderStatusPendingPayment {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен или отменён",
		}, nil
	}

	grpcPaymentMethod, err := mapHTTPPaymentMethodToGRPC(
		req.PaymentMethod,
	)
	if err != nil {
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	paymentCtx, cancel := context.WithTimeout(ctx, grpcRequestTimeout)
	defer cancel()

	paymentResponse, err := h.paymentClient.PayOrder(
		paymentCtx,
		&paymentv1.PayOrderRequest{
			OrderUuid:     params.OrderUUID.String(),
			PaymentMethod: grpcPaymentMethod,
		},
	)
	if err != nil {
		return mapPaymentErrorToPayOrderResponse(err), nil
	}

	transactionUUID, err := uuid.Parse(
		paymentResponse.GetTransactionUuid(),
	)
	if err != nil {
		slog.Error(
			"payment service вернул невалидный transaction_uuid",
			"transaction_uuid", paymentResponse.GetTransactionUuid(),
			"error", err,
		)

		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "payment service вернул некорректный transaction_uuid",
		}, nil
	}

	paymentMethod := PaymentMethod(req.PaymentMethod)

	updatedOrder, found, updateErr := h.store.Update(
		params.OrderUUID,
		func(currentOrder *Order) error {
			// Между чтением заказа и ответом Payment Service
			// его состояние могло измениться.
			if currentOrder.Status != OrderStatusPendingPayment {
				return errOrderStateChanged
			}

			currentOrder.Status = OrderStatusPaid
			currentOrder.TransactionUUID = &transactionUUID
			currentOrder.PaymentMethod = &paymentMethod

			return nil
		},
	)
	if !found {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if errors.Is(updateErr, errOrderStateChanged) {
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "состояние заказа изменилось во время оплаты",
		}, nil
	}

	if updateErr != nil {
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "не удалось обновить заказ после оплаты",
		}, nil
	}

	return &orderv1.PayOrderResponse{
		TransactionUUID: *updatedOrder.TransactionUUID,
	}, nil
}

// CancelOrder реализует POST /api/v1/orders/{order_uuid}/cancel.
func (h *Handler) CancelOrder(
	_ context.Context,
	params orderv1.CancelOrderParams,
) (orderv1.CancelOrderRes, error) {
	_, found, err := h.store.Update(
		params.OrderUUID,
		func(order *Order) error {
			if order.Status != OrderStatusPendingPayment {
				return errOrderCannotBeCancelled
			}

			order.Status = OrderStatusCancelled

			return nil
		},
	)
	if !found {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	if errors.Is(err, errOrderCannotBeCancelled) {
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "оплаченный или уже отменённый заказ нельзя отменить",
		}, nil
	}

	if err != nil {
		return &orderv1.CancelOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "не удалось отменить заказ",
		}, nil
	}

	return &orderv1.CancelOrderResponse{}, nil
}

var (
	errOrderStateChanged      = errors.New("статус заказа изменен")
	errOrderCannotBeCancelled = errors.New("заказ невозможно отменить")
)

// orderToDTO преобразует внутреннюю модель заказа в HTTP DTO.
func orderToDTO(order Order) *orderv1.OrderDto {
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(
			*order.TransactionUUID,
		)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(
			orderv1.PaymentMethod(*order.PaymentMethod),
		)
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}
}

// mapHTTPPaymentMethodToGRPC преобразует HTTP enum в protobuf enum.
func mapHTTPPaymentMethodToGRPC(
	method orderv1.PaymentMethod,
) (paymentv1.PaymentMethod, error) {
	switch method {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD, nil

	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP, nil

	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil

	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil

	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED,
			fmt.Errorf("неподдерживаемый способ оплаты: %q", method)
	}
}

// mapInventoryErrorToCreateOrderResponse преобразует gRPC-ошибку
// Inventory Service в HTTP-ответ CreateOrder.
func mapInventoryErrorToCreateOrderResponse(
	err error,
) orderv1.CreateOrderRes {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		slog.Error(
			"inventory service вернул неизвестную ошибку",
			"error", err,
		)

		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "неизвестная ошибка inventory service",
		}
	}

	slog.Error(
		"ошибка inventory service",
		"code", grpcStatus.Code(),
		"message", grpcStatus.Message(),
	)

	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		return &orderv1.CreateOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: grpcStatus.Message(),
		}

	case codes.NotFound:
		return &orderv1.CreateOrderNotFound{
			Code:    http.StatusNotFound,
			Message: grpcStatus.Message(),
		}

	case codes.DeadlineExceeded:
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "превышено время ожидания inventory service",
		}

	case codes.Unavailable:
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "inventory service временно недоступен",
		}

	default:
		return &orderv1.CreateOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка inventory service",
		}
	}
}

// mapPaymentErrorToPayOrderResponse преобразует gRPC-ошибку
// Payment Service в HTTP-ответ PayOrder.
func mapPaymentErrorToPayOrderResponse(
	err error,
) orderv1.PayOrderRes {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		slog.Error(
			"payment service вернул неизвестную ошибку",
			"error", err,
		)

		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "неизвестная ошибка payment service",
		}
	}

	slog.Error(
		"ошибка payment service",
		"code", grpcStatus.Code(),
		"message", grpcStatus.Message(),
	)

	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		return &orderv1.PayOrderBadRequest{
			Code:    http.StatusBadRequest,
			Message: grpcStatus.Message(),
		}

	case codes.DeadlineExceeded:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "превышено время ожидания payment service",
		}

	case codes.Unavailable:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "payment service временно недоступен",
		}

	default:
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка payment service",
		}
	}
}
