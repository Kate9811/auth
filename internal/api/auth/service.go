package auth

import (
	"github.com/Denis/project_auth/internal/service"
	desc "github.com/Denis/project_auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
