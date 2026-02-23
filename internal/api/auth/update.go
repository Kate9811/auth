package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/converter"
	"github.com/Denis/project_auth/internal/logger"
	"go.uber.org/zap"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Update - обновление пользователя
func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	logger.Info("updating user",
		zap.Int64("user_id", req.Id),
		zap.String("name", req.Name.Value),
	)
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
		logger.Error("failed to update user",
			zap.Int64("user_id", req.Id),
		)

		return nil, status.Error(codes.Internal, "failed to update user")
	}
	logger.Info("updated user with id",
		zap.Int64("user_id", req.Id),
	)

	return &emptypb.Empty{}, nil
}
