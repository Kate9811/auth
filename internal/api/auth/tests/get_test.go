// internal/service/auth/tests/get_test.go
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

	"github.com/Denis/project_auth/internal/model"
	repoMocks "github.com/Denis/project_auth/internal/repository/mocks"
	"github.com/Denis/project_auth/internal/service/auth"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
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

		user = &model.Auth{
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
		name          string
		args          args
		want          *model.Auth
		err           error
		setupRepoMock func(*repoMocks.AuthRepositoryMock)
	}{
		{
			name: "success case",
			args: args{ctx: ctx, id: id},
			want: user,
			err:  nil,
			setupRepoMock: func(repoMock *repoMocks.AuthRepositoryMock) {
				repoMock.GetMock.Expect(ctx, id).Return(user, nil)
			},
		},
		{
			name: "not found error",
			args: args{ctx: ctx, id: id},
			want: nil,
			err:  repoErr,
			setupRepoMock: func(repoMock *repoMocks.AuthRepositoryMock) {
				repoMock.GetMock.Expect(ctx, id).Return(nil, repoErr)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := repoMocks.NewAuthRepositoryMock(mc)
			tt.setupRepoMock(repoMock)

			service := auth.NewMockService(repoMock)

			got, err := service.Get(tt.args.ctx, tt.args.id)

			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, got)
		})
	}
}
