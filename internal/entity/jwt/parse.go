package jwt

import (
	"context"
	"strings"

	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/utils"
	"google.golang.org/grpc/metadata"
)

func ParseTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}

	bearerToken, ok := md["authorization"]
	if !ok {
		return "", utils.UnauthenticatedResponse()
	}

	if len(bearerToken) == 0 {
		return "", utils.UnauthenticatedResponse()
	}

	splitToken := strings.Split(bearerToken[0], " ")

	if len(splitToken) != 2 {
		return "", utils.UnauthenticatedResponse()
	}

	if splitToken[0] != "Bearer" {
		return "", utils.UnauthenticatedResponse()
	}

	return splitToken[1], nil
}
