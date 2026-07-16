package service

import (
	"context"
	"slices"

	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type InventoryGrpcServer struct {
	inventoryv1.UnimplementedInventoryServiceServer
	parts map[uuid.UUID]Part
}

func NewInventoryGrpcServer() *InventoryGrpcServer {
	now := timestamppb.Now()

	return &InventoryGrpcServer{
		parts: map[uuid.UUID]Part{
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440001",
				Name:          "Алюминиевый корпус",
				Description:   "Лёгкий корпус для небольших кораблей",
				Price:         500000, // 5000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 10,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440002",
				Name:          "Титановый корпус",
				Description:   "Прочный корпус для средних кораблей",
				Price:         1500000, // 15000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 5,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440003",
				Name:          "Ионный двигатель C",
				Description:   "Базовый ионный двигатель класса C",
				Price:         300000, // 3000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 8,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440004",
				Name:          "Ионный двигатель B",
				Description:   "Улучшенный ионный двигатель класса B",
				Price:         800000, // 8000₽
				PartType:      inventoryv1.PartType_PART_TYPE_ENGINE,
				StockQuantity: 3,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440005",
				Name:          "Энергетический щит",
				Description:   "Стандартный энергетический щит",
				Price:         400000, // 4000₽
				PartType:      inventoryv1.PartType_PART_TYPE_SHIELD,
				StockQuantity: 6,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440006"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440006",
				Name:          "Лазерная пушка",
				Description:   "Точная лазерная пушка",
				Price:         250000, // 2500₽
				PartType:      inventoryv1.PartType_PART_TYPE_WEAPON,
				StockQuantity: 7,
				CreatedAt:     now,
			},
			uuid.MustParse("550e8400-e29b-41d4-a716-446655440007"): {
				UUID:          "550e8400-e29b-41d4-a716-446655440007",
				Name:          "Плазменный корпус",
				Description:   "Экспериментальный корпус (нет на складе)",
				Price:         2000000, // 20000₽
				PartType:      inventoryv1.PartType_PART_TYPE_HULL,
				StockQuantity: 0,
				CreatedAt:     now,
			},
		},
	}
}

// GetPart Возвращает деталь по UUID
func (s *InventoryGrpcServer) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	if req.GetUuid() == "" {
		return nil, status.Error(codes.InvalidArgument, "UUID детали обязателен")
	}
	partUuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "невалидный UUID детали: %s", partUuid.String())
	}

	part, ok := s.parts[partUuid]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "деталь по переданному UUID не найдена: %s", partUuid.String())
	}

	return &inventoryv1.GetPartResponse{
		Part: toProtoPart(part),
	}, nil
}

// ListParts Возвращает список деталей с возможностью фильтрации по типу или конкретным UUID
func (s *InventoryGrpcServer) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	uuids := req.GetUuids()

	// Если передан список uuids
	if len(uuids) > 0 {
		var parts []*inventoryv1.Part
		for _, val := range uuids {
			partUuid, err := uuid.Parse(val)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "невалидный UUID детали: %s", partUuid)
			}

			part, ok := s.parts[partUuid]
			if !ok {
				return nil, status.Errorf(codes.NotFound, "деталь по переданному UUID не найдена: %s", partUuid)
			}
			parts = append(parts, toProtoPart(part))
		}

		return &inventoryv1.ListPartsResponse{
			Parts: parts,
		}, nil
	}

	// Если тип не указан (PART_TYPE_UNSPECIFIED)
	if req.GetPartType() == inventoryv1.PartType_PART_TYPE_UNSPECIFIED {
		var parts []*inventoryv1.Part

		for _, part := range s.parts {
			parts = append(parts, toProtoPart(part))
		}
		return &inventoryv1.ListPartsResponse{
			Parts: parts,
		}, nil
	}

	// Фильтрация по типу
	var parts []*inventoryv1.Part
	for _, part := range s.parts {
		if part.PartType == req.GetPartType() {
			parts = append(parts, toProtoPart(part))
		}
	}

	// Сортировка по имени
	slices.SortFunc(parts, func(a, b *inventoryv1.Part) int {
		if a.Name < b.Name {
			return -1
		} else if a.Name > b.Name {
			return 1
		}
		return 0
	})

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      inventoryv1.PartType
	StockQuantity int64
	CreatedAt     *timestamppb.Timestamp
}

func toProtoPart(part Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      part.PartType,
		StockQuantity: part.StockQuantity,
		CreatedAt:     part.CreatedAt,
	}
}
