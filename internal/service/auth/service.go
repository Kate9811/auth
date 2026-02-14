// internal/service/auth/service.go
package auth

import (
	"github.com/Denis/project_auth/internal/repository"
	"github.com/Denis/project_auth/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
}

func NewService(authRepository repository.AuthRepository) service.AuthService {
	return &serv{
		authRepository: authRepository,
	}
}

// NewMockService для тестов - максимально простой
func NewMockService(authRepo repository.AuthRepository) service.AuthService {
	return &serv{
		authRepository: authRepo,
	}
}
