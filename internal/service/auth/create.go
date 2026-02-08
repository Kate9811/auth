package auth

import (
	"context"
	"log"

	"github.com/Denis/project_auth/internal/converter"
	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	log.Printf("[SERVICE] Create called: Name=%q, Email=%q, Role=%v",
		req.GetName(), req.GetEmail(), req.GetRole())

	// Валидация паролей
	if req.GetPassword() != req.GetPasswordConfirm() {
		log.Printf("[SERVICE] Password mismatch")
		return nil, status.Error(codes.InvalidArgument, "passwords don't match")
	}

	// Конвертируем protobuf в модель
	authInfo, err := converter.ToAuthInfoFromCreate(req)
	if err != nil {
		log.Printf("[SERVICE] Converter error: %v", err)
		return nil, status.Error(codes.Internal, "failed to process password")
	}

	log.Printf("[SERVICE] AuthInfo: Name=%q, Email=%q, Role=%q, PasswordHashLen=%d",
		authInfo.Name, authInfo.Email, authInfo.Role, len(authInfo.PasswordHash))

	// Сохраняем в БД
	id, err := s.authRepository.Create(ctx, authInfo)
	if err != nil {
		log.Printf("[SERVICE] Repository.Create error: %v", err)
		return nil, status.Error(codes.Internal, "failed to create user")
	}

	log.Printf("[SERVICE] User created with ID: %d", id)
	return &desc.CreateResponse{Id: id}, nil
}
