package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/logger"
	desc "github.com/Denis/project_auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete - удаление пользователя
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	// Логируем ВХОДЯЩИЙ запрос
	logger.Info("deleting user",
		zap.Int64("user_id", req.GetId()),
	)

	// Валидация
	if req.GetId() <= 0 {
		logger.Warn("invalid user id for deletion",
			zap.Int64("user_id", req.GetId()),
		)
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	// Вызываем сервисный слой
	err := i.authService.Delete(ctx, req.GetId())
	if err != nil {
		// Логируем ошибку с полным контекстом
		logger.Error("failed to delete user",
			zap.Int64("user_id", req.GetId()),
			zap.Error(err),
		)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	// 4. Логируем успешное удаление (важное бизнес-событие!)
	logger.Info("user successfully deleted",
		zap.Int64("user_id", req.GetId()),
	)
	return &emptypb.Empty{}, nil
}
