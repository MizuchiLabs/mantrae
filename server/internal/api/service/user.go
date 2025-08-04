package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"connectrpc.com/connect"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	"github.com/mizuchilabs/mantrae/pkg/util"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
	"github.com/mizuchilabs/mantrae/server/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/server/internal/config"
	"github.com/mizuchilabs/mantrae/server/internal/mail"
	"github.com/mizuchilabs/mantrae/server/internal/settings"
	"github.com/mizuchilabs/mantrae/server/internal/store/db"
)

type UserService struct {
	app *config.App
}

func NewUserService(app *config.App) *UserService {
	return &UserService{app: app}
}

func (s *UserService) LoginUser(
	ctx context.Context,
	req *connect.Request[mantraev1.LoginUserRequest],
) (*connect.Response[mantraev1.LoginUserResponse], error) {
	var user db.User
	var err error
	switch id := req.Msg.GetIdentifier().(type) {
	case *mantraev1.LoginUserRequest_Username:
		user, err = s.app.Conn.GetQuery().GetUserByUsername(ctx, id.Username)
	case *mantraev1.LoginUserRequest_Email:
		user, err = s.app.Conn.GetQuery().GetUserByEmail(ctx, &id.Email)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if ok := util.VerifyPassword(req.Msg.Password, user.Password); !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid password"))
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	token, err := meta.EncodeUserToken(user.ID, s.app.Secret, expirationTime)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.GetQuery().UpdateUserLastLogin(ctx, user.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	cookie := http.Cookie{
		Name:     meta.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(expirationTime.Unix() - time.Now().Unix()),
		Secure:   req.Header().Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}
	res := connect.NewResponse(&mantraev1.LoginUserResponse{User: user.ToProto()})
	res.Header().Set("Set-Cookie", cookie.String())
	return res, nil
}

func (s *UserService) LogoutUser(
	ctx context.Context,
	req *connect.Request[mantraev1.LogoutUserRequest],
) (*connect.Response[mantraev1.LogoutUserResponse], error) {
	cookie := http.Cookie{
		Name:     meta.CookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		Secure:   req.Header().Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}
	res := connect.NewResponse(&mantraev1.LogoutUserResponse{})
	res.Header().Set("Set-Cookie", cookie.String())
	return res, nil
}

func (s *UserService) VerifyOTP(
	ctx context.Context,
	req *connect.Request[mantraev1.VerifyOTPRequest],
) (*connect.Response[mantraev1.VerifyOTPResponse], error) {
	var user db.User
	var err error
	switch id := req.Msg.GetIdentifier().(type) {
	case *mantraev1.VerifyOTPRequest_Username:
		user, err = s.app.Conn.GetQuery().GetUserByUsername(ctx, id.Username)
	case *mantraev1.VerifyOTPRequest_Email:
		user, err = s.app.Conn.GetQuery().GetUserByEmail(ctx, &id.Email)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if user.Otp == nil || user.OtpExpiry == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
	}
	if *user.Otp != util.HashOTP(req.Msg.Otp) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid token"))
	}
	if time.Now().After(*user.OtpExpiry) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("token expired"))
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	token, err := meta.EncodeUserToken(user.ID, s.app.Secret, expirationTime)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.GetQuery().UpdateUserLastLogin(ctx, user.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Delete OTP if it's set
	if user.Otp != nil && user.OtpExpiry != nil {
		if err := s.app.Conn.GetQuery().UpdateUserResetToken(ctx, db.UpdateUserResetTokenParams{
			ID:        user.ID,
			Otp:       nil,
			OtpExpiry: nil,
		}); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	cookie := http.Cookie{
		Name:     meta.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(expirationTime.Unix() - time.Now().Unix()),
		Secure:   req.Header().Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}
	res := connect.NewResponse(&mantraev1.VerifyOTPResponse{User: user.ToProto()})
	res.Header().Set("Set-Cookie", cookie.String())
	return res, nil
}

func (s *UserService) SendOTP(
	ctx context.Context,
	req *connect.Request[mantraev1.SendOTPRequest],
) (*connect.Response[mantraev1.SendOTPResponse], error) {
	var user db.User
	var err error
	switch id := req.Msg.GetIdentifier().(type) {
	case *mantraev1.SendOTPRequest_Username:
		user, err = s.app.Conn.GetQuery().GetUserByUsername(ctx, id.Username)
	case *mantraev1.SendOTPRequest_Email:
		user, err = s.app.Conn.GetQuery().GetUserByEmail(ctx, &id.Email)
	default:
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Generate OTP
	expiresAt := time.Now().Add(15 * time.Minute)
	token, err := util.GenerateOTP()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	hash := util.HashOTP(token)

	if err := s.app.Conn.GetQuery().UpdateUserResetToken(ctx, db.UpdateUserResetTokenParams{
		ID:        user.ID,
		Otp:       &hash,
		OtpExpiry: &expiresAt,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	data := map[string]any{
		"Token": token,
		"Date":  expiresAt.Format("Jan 2, 2006 at 15:04"),
	}
	if err := mail.Send(s.app.SM, *user.Email, "reset-password", data); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.SendOTPResponse{}), nil
}

func (s *UserService) GetUser(
	ctx context.Context,
	req *connect.Request[mantraev1.GetUserRequest],
) (*connect.Response[mantraev1.GetUserResponse], error) {
	var user db.User
	var err error
	switch id := req.Msg.GetIdentifier().(type) {
	case *mantraev1.GetUserRequest_Id:
		user, err = s.app.Conn.GetQuery().GetUserByID(ctx, id.Id)
	case *mantraev1.GetUserRequest_Username:
		user, err = s.app.Conn.GetQuery().GetUserByUsername(ctx, id.Username)
	case *mantraev1.GetUserRequest_Email:
		user, err = s.app.Conn.GetQuery().GetUserByEmail(ctx, &id.Email)
	default:
		userID := middlewares.GetUserIDFromContext(ctx)
		if userID == nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
		}
		user, err = s.app.Conn.GetQuery().GetUserByID(ctx, *userID)
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.GetUserResponse{
		User: user.ToProto(),
	}), nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateUserRequest],
) (*connect.Response[mantraev1.CreateUserResponse], error) {
	params := db.CreateUserParams{
		ID:       uuid.NewString(),
		Username: req.Msg.Username,
		Email:    req.Msg.Email,
	}

	var err error
	params.Password, err = util.HashPassword(req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	result, err := s.app.Conn.GetQuery().CreateUser(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_CREATED,
		Data: &mantraev1.EventStreamResponse_User{
			User: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.CreateUserResponse{
		User: result.ToProto(),
	}), nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateUserRequest],
) (*connect.Response[mantraev1.UpdateUserResponse], error) {
	params := db.UpdateUserParams{
		ID:       req.Msg.Id,
		Username: req.Msg.Username,
		Email:    req.Msg.Email,
	}
	result, err := s.app.Conn.GetQuery().UpdateUser(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Update password if provided
	if req.Msg.Password != nil {
		hash, err := util.HashPassword(*req.Msg.Password)
		if err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		if err := s.app.Conn.GetQuery().UpdateUserPassword(ctx, db.UpdateUserPasswordParams{
			ID:       result.ID,
			Password: hash,
		}); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_UPDATED,
		Data: &mantraev1.EventStreamResponse_User{
			User: result.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.UpdateUserResponse{
		User: result.ToProto(),
	}), nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteUserRequest],
) (*connect.Response[mantraev1.DeleteUserResponse], error) {
	user, err := s.app.Conn.GetQuery().GetUserByID(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if err := s.app.Conn.GetQuery().DeleteUser(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.app.Event.Broadcast(&mantraev1.EventStreamResponse{
		Action: mantraev1.EventAction_EVENT_ACTION_DELETED,
		Data: &mantraev1.EventStreamResponse_User{
			User: user.ToProto(),
		},
	})
	return connect.NewResponse(&mantraev1.DeleteUserResponse{}), nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
	req *connect.Request[mantraev1.ListUsersRequest],
) (*connect.Response[mantraev1.ListUsersResponse], error) {
	params := db.ListUsersParams{
		Limit:  req.Msg.Limit,
		Offset: req.Msg.Offset,
	}

	result, err := s.app.Conn.GetQuery().ListUsers(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountUsers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	users := make([]*mantraev1.User, 0, len(result))
	for _, u := range result {
		users = append(users, u.ToProto())
	}
	return connect.NewResponse(&mantraev1.ListUsersResponse{
		Users:      users,
		TotalCount: totalCount,
	}), nil
}

func (s *UserService) GetOIDCStatus(
	ctx context.Context,
	req *connect.Request[mantraev1.GetOIDCStatusRequest],
) (*connect.Response[mantraev1.GetOIDCStatusResponse], error) {
	sets := s.app.SM.GetMany(
		ctx,
		[]string{
			settings.KeyOIDCEnabled,
			settings.KeyPasswordLoginEnabled,
			settings.KeyOIDCProviderName,
		},
	)
	return connect.NewResponse(&mantraev1.GetOIDCStatusResponse{
		OidcEnabled:  sets[settings.KeyOIDCEnabled] == "true",
		LoginEnabled: sets[settings.KeyPasswordLoginEnabled] == "true",
		Provider:     sets[settings.KeyOIDCProviderName],
	}), nil
}
