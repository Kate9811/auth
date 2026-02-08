package service

import (
	"context"

	"github.com/Denis/project_auth/internal/model"
)

type AuthService interface {
	Create(ctx context.Context, info *model.AuthInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.Auth, error)
	Update(ctx context.Context, id int64, info *model.AuthInfo) error
	Delete(ctx context.Context, id int64) error
}
