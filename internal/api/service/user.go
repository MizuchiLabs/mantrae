package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/mizuchilabs/mantrae/internal/api/middlewares"
	"github.com/mizuchilabs/mantrae/internal/config"
	"github.com/mizuchilabs/mantrae/internal/convert"
	"github.com/mizuchilabs/mantrae/internal/mail"
	"github.com/mizuchilabs/mantrae/internal/settings"
	"github.com/mizuchilabs/mantrae/internal/store/db"
	"github.com/mizuchilabs/mantrae/internal/util"
	"github.com/mizuchilabs/mantrae/pkg/meta"
	mantraev1 "github.com/mizuchilabs/mantrae/proto/gen/mantrae/v1"
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
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("one of username or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Msg.Password)); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
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
	res := connect.NewResponse(&mantraev1.LoginUserResponse{Token: token})
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

func (s *UserService) VerifyJWT(
	ctx context.Context,
	req *connect.Request[mantraev1.VerifyJWTRequest],
) (*connect.Response[mantraev1.VerifyJWTResponse], error) {
	val := ctx.Value(middlewares.AuthUserIDKey)
	userID, ok := val.(string)
	if !ok || userID == "" {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("unauthenticated"))
	}
	result, err := s.app.Conn.GetQuery().GetUserByID(ctx, userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.VerifyJWTResponse{
		User: convert.UserToProto(&result),
	}), nil
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
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("one of username or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if user.Otp == nil || user.OtpExpiry == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("no OTP token found"))
	}

	if *user.Otp != util.HashOTP(req.Msg.Otp) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("invalid OTP token"))
	}

	if time.Now().After(*user.OtpExpiry) {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("OTP token expired"))
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	token, err := meta.EncodeUserToken(user.ID, s.app.Secret, expirationTime)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := s.app.Conn.GetQuery().UpdateUserLastLogin(ctx, user.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.VerifyOTPResponse{Token: token}), nil
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
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("one of username or email must be set"))
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
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("one of id, username, or email must be set"))
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.GetUserResponse{
		User: convert.UserToProto(&user),
	}), nil
}

func (s *UserService) CreateUser(
	ctx context.Context,
	req *connect.Request[mantraev1.CreateUserRequest],
) (*connect.Response[mantraev1.CreateUserResponse], error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	params := db.CreateUserParams{
		ID:       id.String(),
		Username: req.Msg.Username,
		Password: req.Msg.Password,
		IsAdmin:  req.Msg.IsAdmin,
	}
	if req.Msg.Email != "" {
		params.Email = &req.Msg.Email
	}

	result, err := s.app.Conn.GetQuery().CreateUser(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.CreateUserResponse{
		User: convert.UserToProto(&result),
	}), nil
}

func (s *UserService) UpdateUser(
	ctx context.Context,
	req *connect.Request[mantraev1.UpdateUserRequest],
) (*connect.Response[mantraev1.UpdateUserResponse], error) {
	params := db.UpdateUserParams{
		ID:       req.Msg.Id,
		Username: req.Msg.Username,
		IsAdmin:  req.Msg.IsAdmin,
	}
	if req.Msg.Email != "" {
		params.Email = &req.Msg.Email
	}
	user, err := s.app.Conn.GetQuery().UpdateUser(ctx, params)
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
			ID:       user.ID,
			Password: hash,
		}); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	return connect.NewResponse(&mantraev1.UpdateUserResponse{
		User: convert.UserToProto(&user),
	}), nil
}

func (s *UserService) DeleteUser(
	ctx context.Context,
	req *connect.Request[mantraev1.DeleteUserRequest],
) (*connect.Response[mantraev1.DeleteUserResponse], error) {
	if err := s.app.Conn.GetQuery().DeleteUser(ctx, req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(&mantraev1.DeleteUserResponse{}), nil
}

func (s *UserService) ListUsers(
	ctx context.Context,
	req *connect.Request[mantraev1.ListUsersRequest],
) (*connect.Response[mantraev1.ListUsersResponse], error) {
	var params db.ListUsersParams
	if req.Msg.Limit == nil {
		params.Limit = 100
	} else {
		params.Limit = *req.Msg.Limit
	}
	if req.Msg.Offset == nil {
		params.Offset = 0
	} else {
		params.Offset = *req.Msg.Offset
	}

	result, err := s.app.Conn.GetQuery().ListUsers(ctx, params)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	totalCount, err := s.app.Conn.GetQuery().CountUsers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&mantraev1.ListUsersResponse{
		Users:      convert.UsersToProto(result),
		TotalCount: totalCount,
	}), nil
}

func (s *UserService) GetOIDCStatus(
	ctx context.Context,
	req *connect.Request[mantraev1.GetOIDCStatusRequest],
) (*connect.Response[mantraev1.GetOIDCStatusResponse], error) {
	oidcEnabled, ok := s.app.SM.Get(settings.KeyOIDCEnabled)
	if !ok {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("failed to get oidc enabled setting"),
		)
	}
	loginEnabled, ok := s.app.SM.Get(settings.KeyPasswordLoginEnabled)
	if !ok {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("failed to get login disabled setting"),
		)
	}
	provider, ok := s.app.SM.Get(settings.KeyOIDCProviderName)
	if !ok {
		return nil, connect.NewError(
			connect.CodeInternal,
			errors.New("failed to get oidc provider name setting"),
		)
	}
	return connect.NewResponse(&mantraev1.GetOIDCStatusResponse{
		OidcEnabled:  oidcEnabled == "true",
		LoginEnabled: loginEnabled == "true",
		Provider:     provider,
	}), nil
}
