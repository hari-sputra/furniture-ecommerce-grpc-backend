package handler

import (
	"context"

	"github.com/hari-sputra/furniture-ecommerce-grpc-be/internal/utils"
	"github.com/hari-sputra/furniture-ecommerce-grpc-be/pb/service"
)

type serviceHandler struct {
	service.UnimplementedHelloWorldServiceServer
}

func (sh *serviceHandler) SayHello(ctx context.Context, request *service.HelloWorldRequest) (*service.HelloWorldResponse, error) {

	validation, err := utils.CheckValidation(request)

	if err != nil {
		return nil, err
	}

	if validation != nil {
		return &service.HelloWorldResponse{
			Meta: utils.ValidationErrorResponse(validation),
		}, nil
	}

	return &service.HelloWorldResponse{
		Message: "Hello, " + request.Name,
		Meta:    utils.SuccessResponse(),
	}, nil
}

func NewServiceHandler() *serviceHandler {
	return &serviceHandler{}
}
