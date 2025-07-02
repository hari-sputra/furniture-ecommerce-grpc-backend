package handler

import (
	"context"

	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/service"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/utils"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pb/auth"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	auth.UnimplementedAuthServiceServer
	authService service.IAuthService
}

func (ah *authHandler) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {

	validation, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validation != nil {
		return &auth.RegisterResponse{
			Meta: utils.ValidationErrorResponse(validation),
		}, nil
	}

	// register process
	res, err := ah.authService.Register(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ah *authHandler) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {

	validation, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validation != nil {
		return &auth.LoginResponse{
			Meta: utils.ValidationErrorResponse(validation),
		}, nil
	}

	res, err := ah.authService.Login(ctx, req)

	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return &auth.LoginResponse{
				Meta: utils.BadRequestResponse(st.Message()),
			}, nil
		}

		return &auth.LoginResponse{
			Meta: utils.BadRequestResponse(err.Error()),
		}, nil
	}

	return res, nil
}

func (ah *authHandler) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {

	validation, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validation != nil {
		return &auth.LogoutResponse{
			Meta: utils.ValidationErrorResponse(validation),
		}, nil
	}

	res, err := ah.authService.Logout(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ah *authHandler) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	validation, err := utils.CheckValidation(req)

	if err != nil {
		return nil, err
	}

	if validation != nil {
		return &auth.ChangePasswordResponse{
			Meta: utils.ValidationErrorResponse(validation),
		}, nil
	}

	res, err := ah.authService.ChangePassword(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ah *authHandler) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {

	res, err := ah.authService.GetProfile(ctx, req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewAuthHandler(authService service.IAuthService) *authHandler {
	return &authHandler{
		authService: authService,
	}
}
