package auth

import (
	"context"
	"log"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Delete - удаление пользователя
func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	// Валидация
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid user id")
	}

	// Вызываем сервисный слой
	err := i.authService.Delete(ctx, req.GetId())
	if err != nil {
		log.Printf("failed to delete user %d: %v", req.GetId(), err)
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	log.Printf("deleted user with id: %d", req.GetId())
	return &emptypb.Empty{}, nil
}
