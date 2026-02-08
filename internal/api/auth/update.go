package auth

import (
	"context"
	"log"

	"github.com/Denis/project_auth/internal/converter"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update - обновление пользователя
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	// Валидация
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	// Проверяем, что есть что обновлять
	if req.GetName() == nil && req.GetEmail() == nil {
		return nil, status.Error(codes.InvalidArgument, "nothing to update")
	}

	// Конвертируем protobuf в модель
	authInfo := converter.ToAuthInfoFromUpdate(req)

	// Вызываем сервисный слой
	err := i.authService.Update(ctx, req.GetId(), authInfo)
	if err != nil {
		log.Printf("failed to update user %d: %v", req.GetId(), err)
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	log.Printf("updated user with id: %d", req.GetId())
	return &emptypb.Empty{}, nil
}
