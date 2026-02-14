package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/model"
)

func (s *serv) Update(ctx context.Context, id int64, authInfo *model.AuthInfo) error {
	err := s.authRepository.Update(ctx, id, authInfo)
	if err != nil {
		return err
	}

	return nil
}
