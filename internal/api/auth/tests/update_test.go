// internal/service/auth/tests/update_test.go
package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/Denis/project_auth/internal/model"
	repoMocks "github.com/Denis/project_auth/internal/repository/mocks"
	"github.com/Denis/project_auth/internal/service/auth"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		id   int64
		info *model.AuthInfo
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id       = gofakeit.Int64()
		name     = gofakeit.Name()
		email    = gofakeit.Email()
		password = gofakeit.Password(true, true, true, true, false, 10)
		role     = "user"

		info = &model.AuthInfo{
			Name:         name,
			Email:        email,
			PasswordHash: password,
			Role:         role,
		}

		repoErr = fmt.Errorf("repository error")
	)

	t.Cleanup(mc.Finish)

	tests := []struct {
		name          string
		args          args
		err           error
		setupRepoMock func(*repoMocks.AuthRepositoryMock)
	}{
		{
			name: "success case",
			args: args{ctx: ctx, id: id, info: info},
			err:  nil,
			setupRepoMock: func(repoMock *repoMocks.AuthRepositoryMock) {
				repoMock.UpdateMock.Expect(ctx, id, info).Return(nil)
			},
		},
		{
			name: "error case",
			args: args{ctx: ctx, id: id, info: info},
			err:  repoErr,
			setupRepoMock: func(repoMock *repoMocks.AuthRepositoryMock) {
				repoMock.UpdateMock.Expect(ctx, id, info).Return(repoErr)
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

			err := service.Update(tt.args.ctx, tt.args.id, tt.args.info)

			require.Equal(t, tt.err, err)
		})
	}
}
