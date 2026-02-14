package tests

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	//"github.com/Denis/project_auth/internal/client/db" // Импортируем пакет с TxManager
	"github.com/Denis/project_auth/internal/model"
	"github.com/Denis/project_auth/internal/repository"
	repoMocks "github.com/Denis/project_auth/internal/repository/mocks"
	"github.com/Denis/project_auth/internal/service/auth"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type authRepositoryMockFunc func(mc *minimock.Controller) repository.AuthRepository

	type args struct {
		ctx context.Context
		req int64
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = "user"
		now      = time.Now()

		repoErr = fmt.Errorf("repository error")

		res = &model.Auth{
			ID: id,
			Info: model.AuthInfo{
				Name:         name,
				Email:        email,
				PasswordHash: password,
				Role:         role,
			},
			CreatedAt: now,
			UpdatedAt: sql.NullTime{
				Time:  now,
				Valid: true,
			},
		}
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name               string
		args               args
		want               *model.Auth
		err                error
		authRepositoryMock authRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: res,
			err:  nil,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				// Настраиваем оба вызова: Create и Get
				mock.GetMock.Expect(ctx, id).Return(res, nil)

				return mock
			},
		},
		{
			name: "get user error",
			args: args{
				ctx: ctx,
				req: id,
			},
			want: nil,
			err:  repoErr,
			authRepositoryMock: func(mc *minimock.Controller) repository.AuthRepository {
				mock := repoMocks.NewAuthRepositoryMock(mc)
				// Только Create возвращает ошибку
				mock.GetMock.Expect(ctx, id).Return(nil, repoErr)

				// Get НЕ должен вызываться, поэтому не настраиваем его
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			authRepoMock := tt.authRepositoryMock(mc)
			service := auth.NewMockService(authRepoMock)

			newID, err := service.Get(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newID)
		})
	}
}
