package auth

import (
	"context"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *serv) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.authRepository.Delete(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &emptypb.Empty{}, nil
}
