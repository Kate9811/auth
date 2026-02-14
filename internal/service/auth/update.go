package auth

import (
	"context"

	"github.com/Denis/project_auth/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) Update(ctx context.Context, id int64, info *model.AuthInfo) error {
	// Update изменяет данные - нужна транзакция!
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		err := s.authRepository.Update(ctx, id, info)
		if err != nil {
			return err
		}

		// Можно добавить дополнительные проверки/операции
		// Например, проверка что пользователь существует
		_, err = s.authRepository.Get(ctx, id)
		return err
	})
	if err != nil {
		return status.Error(codes.Internal, "failed to update user: "+err.Error())
	}

	return nil
}
