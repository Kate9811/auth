package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	// Delete изменяет данные - нужна транзакция!
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// Проверяем, существует ли пользователь перед удалением
		_, err := s.authRepository.Get(ctx, id)
		if err != nil {
			return status.Error(codes.NotFound, "user not found")
		}

		// Удаляем
		return s.authRepository.Delete(ctx, id)
	})

	if err != nil {
		// Проверяем тип ошибки
		if _, ok := status.FromError(err); ok {
			return err // Уже grpc ошибка
		}
		return status.Error(codes.Internal, "failed to delete user: "+err.Error())
	}

	return nil
}
