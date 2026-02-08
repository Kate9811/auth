package auth

import (
	"context"
	"log"

	"github.com/Denis/project_auth/internal/converter"
	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// Конвертируем protobuf запрос в модель
	authInfo, err := converter.ToAuthInfoFromCreate(req)
	if err != nil {
		// Обработка ошибок валидации
		if err == converter.ErrPasswordsNotMatch {
			return nil, status.Error(codes.InvalidArgument, "passwords do not match")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Вызываем сервисный слой
	id, err := i.authService.Create(ctx, authInfo)
	if err != nil {
		log.Printf("failed to create user: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	log.Printf("created user with id: %d", id)
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
