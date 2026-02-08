package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/converter"
	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	authInfo := converter.ToAuthInfoFromUpdate(req)

	err := s.authRepository.Update(ctx, req.GetId(), authInfo)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &emptypb.Empty{}, nil
}
