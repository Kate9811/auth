package auth

import (
	"context"
	"log"

	"github.com/Denis/project_auth/internal/converter"
	"github.com/Denis/project_auth/internal/logger"
	desc "github.com/Denis/project_auth/pkg/user_v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	// err := req.Validate()
	// if err != nil {
	// 	return nil, err
	// }
	logger.Info("creating new user",
		zap.String("email", req.GetEmail()),
		zap.String("name", req.GetName()),
		zap.String("role", req.GetRole().String()),
	)
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
		logger.Info("failed to create user",
			zap.String("email", req.GetEmail()),
			zap.String("name", req.GetName()),
		)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	logger.Info("created user with id",
		zap.Int64("id", id),
	)
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
