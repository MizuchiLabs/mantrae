package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func EntryPointToProto(e *db.EntryPoint) *mantraev1.EntryPoint {
	return &mantraev1.EntryPoint{
		Id:        e.ID,
		ProfileId: e.ProfileID,
		Name:      e.Name,
		Address:   e.Address,
		IsDefault: e.IsDefault,
		CreatedAt: SafeTimestamp(e.CreatedAt),
		UpdatedAt: SafeTimestamp(e.UpdatedAt),
	}
}

func EntryPointsToProto(entryPoints []db.EntryPoint) []*mantraev1.EntryPoint {
	var entryPointsProto []*mantraev1.EntryPoint
	for _, e := range entryPoints {
		entryPointsProto = append(entryPointsProto, EntryPointToProto(&e))
	}
	return entryPointsProto
}
