// internal/service/auth/tests/delete_test.go
package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	repoMocks "github.com/Denis/project_auth/internal/repository/mocks"
	"github.com/Denis/project_auth/internal/service/auth"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx     = context.Background()
		mc      = minimock.NewController(t)
		id      = gofakeit.Int64()
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
			args: args{ctx: ctx, id: id},
			err:  nil,
			setupRepoMock: func(repoMock *repoMocks.AuthRepositoryMock) {
				repoMock.DeleteMock.Expect(ctx, id).Return(nil)
			},
		},
		{
			name: "error case",
			args: args{ctx: ctx, id: id},
			err:  repoErr,
			setupRepoMock: func(repoMock *repoMocks.AuthRepositoryMock) {
				repoMock.DeleteMock.Expect(ctx, id).Return(repoErr)
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

			err := service.Delete(tt.args.ctx, tt.args.id)

			require.Equal(t, tt.err, err)
		})
	}
}
