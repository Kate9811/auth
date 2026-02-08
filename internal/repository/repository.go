package repository

import (
	"context"

	"github.com/Denis/project_auth/internal/model"
)

// UserRepository - интерфейс для работы с пользователями
type AuthRepository interface {
	// // Create создает нового пользователя и возвращает его ID
	// Create(ctx context.Context, info *model.AuthInfo) (int64, error)

	// // Get получает пользователя по ID
	// Get(ctx context.Context, id int64) (*model.Auth, error)

	// // Update обновляет информацию о пользователе
	// Update(ctx context.Context, id int64, info *model.AuthInfo) error

	// // Delete удаляет пользователя по ID
	// Delete(ctx context.Context, id int64) error
	Create(ctx context.Context, authInfo *model.AuthInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Auth, error)
	Update(ctx context.Context, id int64, authInfo *model.AuthInfo) error
	Delete(ctx context.Context, id int64) error
}
