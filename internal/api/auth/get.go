package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/converter"
	"github.com/Denis/project_auth/internal/logger"
	"go.uber.org/zap"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get - получение пользователя по ID
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// Логируем ВХОДЯЩИЙ запрос
	logger.Info("getting user by id",
		zap.Int64("user_id", req.GetId()),
	)
	// Валидация
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	// Вызываем сервисный слой
	auth, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		// Логируем ОШИБКУ с контекстом
		logger.Error("failed to get user from service",
			zap.Int64("user_id", req.GetId()),
			zap.Error(err),
		)
		return nil, status.Error(codes.NotFound, "user not found")
	}
	// Логируем успех (можно на DEBUG, если не критично)
	logger.Debug("user successfully retrieved",
		zap.Int64("user_id", req.GetId()),
	)
	// Конвертируем модель в protobuf ответ
	return converter.AuthToGetResponse(auth), nil
}
