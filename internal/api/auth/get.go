package auth

import (
	"context"
	"log"

	"github.com/Denis/project_auth/internal/converter"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Get - получение пользователя по ID
func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	// Валидация
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	// Вызываем сервисный слой
	auth, err := i.authService.Get(ctx, req.GetId())
	if err != nil {
		log.Printf("failed to get user %d: %v", req.GetId(), err)
		return nil, status.Error(codes.NotFound, "user not found")
	}

	// Конвертируем модель в protobuf ответ
	return converter.AuthToGetResponse(auth), nil
}
