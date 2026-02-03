package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/internal/config"
	mantraev1 "github.com/mizuchilabs/mantrae/internal/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/internal/meta"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
)

type UserService struct {
	app *config.App
}

func NewUserService(app *config.App) *UserService {
	return &UserService{app: app}
}

func (s *UserService) LoginUser(
	ctx context.Context,
	req *mantraev1.LoginUserRequest,
) (*mantraev1.LoginUserResponse, error) {
	ci, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get call info"))
	}

	var user *db.User
	var err error
	switch id := req.GetIdentifier().(type) {
	case *mantraev1.LoginUserRequest_Username:
		user, err = s.app.Conn.Q.GetUserByUsername(ctx, id.Username)
	case *mantraev1.LoginUserRequest_Email:
		user, err = s.app.Conn.Q.GetUserByEmail(ctx, &id.Email)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if ok := util.VerifyPassword(req.Password, user.Password); !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid password"))
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	token, err := meta.EncodeUserToken(user.ID, s.app.Secret, expirationTime)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.Q.UpdateUserLastLogin(ctx, user.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	cookie := http.Cookie{
		Name:     meta.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(expirationTime.Unix() - time.Now().Unix()),
		Secure:   ci.RequestHeader().Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}
	ci.ResponseHeader().Set("Set-Cookie", cookie.String())
	return &mantraev1.LoginUserResponse{User: user.ToProto()}, nil
}

func (s *UserService) LogoutUser(
	ctx context.Context,
	req *mantraev1.LogoutUserRequest,
) (*mantraev1.LogoutUserResponse, error) {
	ci, ok := connect.CallInfoForHandlerContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get call info"))
	}

	cookie := http.Cookie{
		Name:     meta.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   ci.RequestHeader().Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}
	ci.ResponseHeader().Set("Set-Cookie", cookie.String())
	return &mantraev1.LogoutUserResponse{}, nil
}

func (s *UserService) GetUser(
	ctx context.Context,
	req *mantraev1.GetUserRequest,
) (*mantraev1.GetUserResponse, error) {
	var user *db.User
	var err error
	switch id := req.GetIdentifier().(type) {
	case *mantraev1.GetUserRequest_Id:
		user, err = s.app.Conn.Q.GetUserByID(ctx, id.Id)
	case *mantraev1.GetUserRequest_Username:
		user, err = s.app.Conn.Q.GetUserByUsername(ctx, id.Username)
	case *mantraev1.GetUserRequest_Email:
		user, err = s.app.Conn.Q.GetUserByEmail(ctx, &id.Email)
	default:
		userID := middlewares.GetUserIDFromContext(ctx)
		if userID == nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
		}
		user, err = s.app.Conn.Q.GetUserByID(ctx, *userID)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &mantraev1.GetUserResponse{User: user.ToProto()}, nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req *mantraev1.CreateUserRequest,
) (*mantraev1.CreateUserResponse, error) {
	params := &db.CreateUserParams{
		ID:       uuid.NewString(),
		Username: req.Username,
		Email:    req.Email,
	}

	var err error
	params.Password, err = util.HashPassword(req.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := s.app.Conn.Q.CreateUser(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.CreateUserResponse{User: result.ToProto()}, nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	req *mantraev1.UpdateUserRequest,
) (*mantraev1.UpdateUserResponse, error) {
	params := &db.UpdateUserParams{
		ID:       req.Id,
		Username: req.Username,
		Email:    req.Email,
	}
	result, err := s.app.Conn.Q.UpdateUser(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update password if provided
	if req.Password != nil {
		hash, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if err := s.app.Conn.Q.UpdateUserPassword(ctx, &db.UpdateUserPasswordParams{
			ID:       result.ID,
			Password: hash,
		}); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	return &mantraev1.UpdateUserResponse{User: result.ToProto()}, nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	req *mantraev1.DeleteUserRequest,
) (*mantraev1.DeleteUserResponse, error) {
	if err := s.app.Conn.Q.DeleteUser(ctx, req.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &mantraev1.DeleteUserResponse{}, nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
	req *mantraev1.ListUsersRequest,
) (*mantraev1.ListUsersResponse, error) {
	params := &db.ListUsersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	result, err := s.app.Conn.Q.ListUsers(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.Q.CountUsers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	users := make([]*mantraev1.User, 0, len(result))
	for _, u := range result {
		users = append(users, u.ToProto())
	}
	return &mantraev1.ListUsersResponse{
		Users:      users,
		TotalCount: totalCount,
	}, nil
}

func (s *UserService) GetOIDCStatus(
	ctx context.Context,
	req *mantraev1.GetOIDCStatusRequest,
) (*mantraev1.GetOIDCStatusResponse, error) {
	sets := s.app.SM.GetMany(
		ctx,
		[]string{
			settings.KeyOIDCEnabled,
			settings.KeyPasswordLoginEnabled,
			settings.KeyOIDCProviderName,
		},
	)
	return &mantraev1.GetOIDCStatusResponse{
		OidcEnabled:  sets[settings.KeyOIDCEnabled] == "true",
		LoginEnabled: sets[settings.KeyPasswordLoginEnabled] == "true",
		Provider:     sets[settings.KeyOIDCProviderName],
	}, nil
}
