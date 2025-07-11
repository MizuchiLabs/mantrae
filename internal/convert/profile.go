package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func ProfileToProto(p *db.Profile) *mantraev1.Profile {
	return &mantraev1.Profile{
		Id:          p.ID,
		Name:        p.Name,
		Description: SafeString(p.Description),
		Token:       p.Token,
		CreatedAt:   SafeTimestamp(p.CreatedAt),
		UpdatedAt:   SafeTimestamp(p.UpdatedAt),
	}
}

func ProfilesToProto(profiles []db.Profile) []*mantraev1.Profile {
	var profilesProto []*mantraev1.Profile
	for _, p := range profiles {
		profilesProto = append(profilesProto, ProfileToProto(&p))
	}
	return profilesProto
}
