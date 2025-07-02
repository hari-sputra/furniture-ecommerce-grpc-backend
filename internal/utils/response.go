package utils

import (
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pb/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func SuccessResponse(messages ...string) *common.BaseResponse {
	message := "Data retrieved successfully"

	if len(messages) > 0 && messages[0] != "" {
		message = messages[0]
	}

	return &common.BaseResponse{
		StatusCode: 200,
		Message:    message,
	}
}

func BadRequestResponse(message string) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode: 400,
		Message:    message,
		IsError:    true,
	}
}

func ValidationErrorResponse(validationErrors []*common.ValidationError) *common.BaseResponse {
	return &common.BaseResponse{
		StatusCode:       400,
		Message:          "Validation Error",
		IsError:          true,
		ValidationErrors: validationErrors,
	}
}

func UnauthenticatedResponse() error {
	return status.Errorf(codes.Unauthenticated, "Unauthorized")
}
