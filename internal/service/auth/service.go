// internal/service/auth/service.go
package auth

import (
	"github.com/Denis/project_auth/internal/repository"
	desc "github.com/Denis/project_auth/pkg/user_v1"
)

type serv struct {
	desc.UnimplementedUserV1Server // ← Встраиваем для совместимости
	authRepository                 repository.AuthRepository
}

func NewService(
	authRepository repository.AuthRepository,

) desc.UserV1Server {
	return &serv{
		authRepository: authRepository,
	}
}
