package converter

import (
	"github.com/Ilya96s/rocket-factory-new/inventory/internal/model"
	inventoryv1 "github.com/Ilya96s/rocket-factory-new/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func InventoryResultToProto(parts []model.Part) *inventoryv1.ListPartsResponse {
	return &inventoryv1.ListPartsResponse{
		Parts: modelPartsToProto(parts),
	}
}

func modelPartsToProto(parts []model.Part) []*inventoryv1.Part {
	protoParts := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, modelPartToProto(part))
	}
	return protoParts
}

func modelPartToProto(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		PartType:      modelPartTypeToProto(part.PartType),
		StockQuantity: part.StockQuantity,
		CreatedAt:     timestamppb.New(part.CreatedAt),
	}
}

func ListPartsRequestToModel(req *inventoryv1.ListPartsRequest) model.PartFilter {
	return model.PartFilter{
		UUIDs:    req.GetUuids(),
		PartType: protoPartTypeToModel(req.GetPartType()),
	}
}

func protoPartTypeToModel(partType inventoryv1.PartType) model.PartType {
	switch partType {
	case inventoryv1.PartType_PART_TYPE_HULL:
		return model.PartTypeHull
	case inventoryv1.PartType_PART_TYPE_ENGINE:
		return model.PartTypeEngine
	case inventoryv1.PartType_PART_TYPE_SHIELD:
		return model.PartTypeShield
	case inventoryv1.PartType_PART_TYPE_WEAPON:
		return model.PartTypeWeapon
	default:
		return model.PartTypeUnspecified
	}
}

func GetPartResultToProto(part model.Part) *inventoryv1.GetPartResponse {
	return &inventoryv1.GetPartResponse{
		Part: &inventoryv1.Part{
			Uuid:          part.UUID,
			Name:          part.Name,
			Description:   part.Description,
			PartType:      modelPartTypeToProto(part.PartType),
			StockQuantity: part.StockQuantity,
			CreatedAt:     timestamppb.New(part.CreatedAt),
		},
	}
}

func modelPartTypeToProto(partType model.PartType) inventoryv1.PartType {
	switch partType {
	case model.PartTypeHull:
		return inventoryv1.PartType_PART_TYPE_HULL
	case model.PartTypeEngine:
		return inventoryv1.PartType_PART_TYPE_ENGINE
	case model.PartTypeShield:
		return inventoryv1.PartType_PART_TYPE_SHIELD
	case model.PartTypeWeapon:
		return inventoryv1.PartType_PART_TYPE_WEAPON
	default:
		return inventoryv1.PartType_PART_TYPE_UNSPECIFIED
	}
}
