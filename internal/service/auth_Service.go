package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/entity"
	JwtEntity "github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/entity/jwt"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/repository"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/utils"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pb/auth"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pb/common"
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthService interface {
	Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error)
	Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error)
	Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error)
	ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error)
	GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error)
}

type authService struct {
	authRepository repository.IAuthRepository
	cacheService   *gocache.Cache
}

func (as *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	if req.Password != req.PasswordConfirmation {
		return &auth.RegisterResponse{
			Meta: utils.BadRequestResponse("Password confirmation does not match"),
		}, nil
	}

	// check email to database
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// if email already exist, return error
	if user != nil {
		return &auth.RegisterResponse{
			Meta: utils.BadRequestResponse("Email already exist"),
		}, nil
	}

	// hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)

	if err != nil {
		return nil, err
	}

	// insert to database
	newUser := entity.User{
		FullName:  req.FullName,
		Email:     req.Email,
		Password:  string(hashedPassword),
		RoleCode:  entity.UserRoleCustomer,
		CreatedBy: &req.FullName,
	}

	as.authRepository.InsertUSer(ctx, &newUser)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Meta: utils.SuccessResponse("Register success"),
	}, nil
}

// Login implements IAuthService.
func (as *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	// email and password checking
	user, err := as.authRepository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.LoginResponse{
			Meta: utils.BadRequestResponse("Email or password is incorrect"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, status.Errorf(codes.Unauthenticated, "Email or password is incorrect")
		}

		return nil, err
	}

	// generate jwt
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtEntity.JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Id,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.RoleCode,
	})

	secretKey := os.Getenv("JWT_SECRET_KEY")
	accesToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	// send response
	return &auth.LoginResponse{
		Meta:        utils.SuccessResponse("Login success"),
		AccessToken: accesToken,
	}, nil
}

// Logout implements IAuthService.
func (as *authService) Logout(ctx context.Context, req *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// get token from metadata grpc
	jwtToken, err := JwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// convert token to entity jwt
	tokenClaims, err := JwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// insert token to cache / memory db
	as.cacheService.Set(jwtToken, "", time.Duration(tokenClaims.ExpiresAt.Time.Unix()-time.Now().Unix())*time.Second)

	return &auth.LogoutResponse{
		Meta: utils.SuccessResponse("Logout successfully"),
	}, nil
}

// ChangePassword implements IAuthService.
func (as *authService) ChangePassword(ctx context.Context, req *auth.ChangePasswordRequest) (*auth.ChangePasswordResponse, error) {
	// check new password and new password confirmation matched
	if req.NewPassword != req.NewPasswordConfirmation {
		return &auth.ChangePasswordResponse{
			Meta: utils.BadRequestResponse("Password confirmation does not match"),
		}, nil
	}

	// check old password
	jwtToken, err := JwtEntity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	claims, err := JwtEntity.GetClaimsFromToken(jwtToken)
	if err != nil {
		return nil, err
	}

	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return &auth.ChangePasswordResponse{
			Meta: utils.BadRequestResponse("User not found"),
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return &auth.ChangePasswordResponse{
				Meta: utils.BadRequestResponse("Old password is not matched"),
			}, nil
		}

		return nil, err
	}

	// update new password ke
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		return nil, err
	}

	updatePassword := entity.UpdateUserPassword{
		UserId:            user.Id,
		HashedNewPassword: string(hashedNewPassword),
		UpdatedBy:         user.FullName,
	}

	err = as.authRepository.UpdateUserPassword(ctx, &updatePassword)
	if err != nil {
		return nil, err
	}

	// send response
	return &auth.ChangePasswordResponse{
		Meta: utils.SuccessResponse("Change password success"),
	}, nil
}

// GetProfile implements IAuthService.
func (as *authService) GetProfile(ctx context.Context, req *auth.GetProfileRequest) (*auth.GetProfileResponse, error) {
	// get data token
	claims, err := JwtEntity.GetClaimsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// retrive data from db
	user, err := as.authRepository.GetUserByEmail(ctx, claims.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return &auth.GetProfileResponse{
			Meta: utils.BadRequestResponse("User not found"),
		}, nil
	}

	// make response
	data := &common.UserResponse{
		UserId:      claims.Subject,
		FullName:    claims.FullName,
		Email:       claims.Email,
		RoleCode:    claims.Role,
		MemberSince: timestamppb.New(user.CreatedAt),
	}

	return &auth.GetProfileResponse{
		Meta: utils.SuccessResponse("Get profile success"),
		Data: data,
	}, nil

}

func NewAuthService(authRepository repository.IAuthRepository, cacheService *gocache.Cache) IAuthService {
	return &authService{
		authRepository: authRepository,
		cacheService:   cacheService,
	}
}
