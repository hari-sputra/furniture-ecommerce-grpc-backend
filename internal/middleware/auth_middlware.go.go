package middleware

import (
	"context"

	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/entity/jwt"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

func (am *authMiddleware) Middleware(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {

	if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" {
		return handler(ctx, req)

	}

	// get token from metadata
	tokenStr, err := jwt.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// check token from logout cache
	_, ok := am.cacheService.Get(tokenStr)
	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// parse jwt to entity
	claims, err := jwt.GetClaimsFromToken(tokenStr)
	if err != nil {
		return nil, err
	}

	// insert entity into context
	ctx = claims.SetToContext(ctx)

	resp, err = handler(ctx, req)

	return resp, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
