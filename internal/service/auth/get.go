package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/converter"
	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	authModel, err := s.authRepository.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return converter.AuthToGetResponse(authModel), nil
}
