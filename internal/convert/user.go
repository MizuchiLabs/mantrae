package convert

import (
	"github.com/mizuchilabs/mantrae/internal/store/db"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
)

func UserToProtoUnsafe(u *db.User) *mantraev1.User {
	return &mantraev1.User{
		Id:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     SafeString(u.Email),
		Otp:       SafeString(u.Otp),
		OtpExpiry: SafeTimestamp(u.OtpExpiry),
		LastLogin: SafeTimestamp(u.LastLogin),
		CreatedAt: SafeTimestamp(u.CreatedAt),
		UpdatedAt: SafeTimestamp(u.UpdatedAt),
	}
}

func UsersToProtoUnsafe(users []db.User) []*mantraev1.User {
	var usersProto []*mantraev1.User
	for _, u := range users {
		usersProto = append(usersProto, UserToProtoUnsafe(&u))
	}
	return usersProto
}

// Safe version omits password and otp fields
func UserToProto(u *db.User) *mantraev1.User {
	return &mantraev1.User{
		Id:        u.ID,
		Username:  u.Username,
		Email:     SafeString(u.Email),
		LastLogin: SafeTimestamp(u.LastLogin),
		CreatedAt: SafeTimestamp(u.CreatedAt),
		UpdatedAt: SafeTimestamp(u.UpdatedAt),
	}
}

func UsersToProto(users []db.User) []*mantraev1.User {
	var usersProto []*mantraev1.User
	for _, u := range users {
		usersProto = append(usersProto, UserToProto(&u))
	}
	return usersProto
}
