package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.Auth, error) {
	authModel, err := s.authRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return authModel, nil
}
