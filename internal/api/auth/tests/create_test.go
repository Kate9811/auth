package test

import (
	"context" // Пакет для работы с контекстом (таймауты, отмена запросов)
	"fmt"     // Для форматирования строк и создания ошибок
	"testing" // Стандартный пакет Go для написания тестов

	"github.com/brianvoe/gofakeit/v6"     // Генерирует случайные данные (имя, email, пароль)
	"github.com/gojuno/minimock/v3"       // Создает моки для интерфейсов
	"github.com/stretchr/testify/require" // Удобные функции для проверок (require.Equal и др.)
	"google.golang.org/grpc/codes"        // Содержит коды gRPC ошибок (Internal, InvalidArgument)
	"google.golang.org/grpc/status"       // Позволяет создавать gRPC ошибки (status.Error)

	"github.com/Denis/project_auth/internal/api/auth"                   // Тестируемый API метод Create
	"github.com/Denis/project_auth/internal/model"                      // Бизнес-модель AuthInfo
	"github.com/Denis/project_auth/internal/service"                    // Интерфейс сервисного слоя
	serviceMocks "github.com/Denis/project_auth/internal/service/mocks" // Сгенерированные моки для сервиса
	desc "github.com/Denis/project_auth/pkg/user_v1"                    // Сгенерированный код из proto-файла
)

// TestCreate тестирует метод Create API слоя
func TestCreate(t *testing.T) {
	t.Parallel() // Раскомментировать для параллельного запуска тестов

	// authServiceMockFunc - тип функции, которая создает мок сервиса
	type authServiceMockFunc func(mc *minimock.Controller) service.AuthService

	// args - структура для хранения аргументов тестируемого метода
	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	// Объявляем переменные, общие для всех тест-кейсов
	var (
		// Базовый контекст для вызовов
		ctx = context.Background()
		// Контроллер minimock - следит за вызовами моков и проверяет ожидания
		mc = minimock.NewController(t)

		// Генерируем случайные данные для тестов
		id              = gofakeit.Int64()                                     // ID пользователя
		name            = gofakeit.Name()                                      // Имя
		email           = gofakeit.Email()                                     // Email
		password        = gofakeit.Password(true, true, true, true, false, 10) // Пароль
		passwordConfirm = password                                             // Для успешного кейса пароли совпадают
		role            = desc.Role_USER                                       // Роль USER

		// Ошибка, которую будет возвращать мок сервиса
		serviceErr = fmt.Errorf("service error")

		// Создаем запрос, который приходит в API (строго по proto)
		req = &desc.CreateRequest{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}

		// Ожидаемый ответ от API
		res = &desc.CreateResponse{
			Id: id,
		}
	)

	// Откладываем проверку, что все ожидаемые вызовы моков произошли
	t.Cleanup(mc.Finish)

	// Таблица тестовых кейсов
	tests := []struct {
		name            string               // Название тест-кейса
		args            args                 // Аргументы для вызова
		want            *desc.CreateResponse // Ожидаемый ответ
		err             error                // Ожидаемая ошибка
		authServiceMock authServiceMockFunc  // Функция создания мока для этого кейса
	}{
		// ТЕСТ-КЕЙС 1: Успешное создание пользователя
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res, // Ожидаем успешный ответ
			err:  nil, // Ошибки не ожидается
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				// Создаем мок сервиса
				mock := serviceMocks.NewAuthServiceMock(mc)

				// Set вместо Expect - мы сами контролируем проверки
				mock.CreateMock.Set(func(ctx context.Context, info *model.AuthInfo) (int64, error) {
					// Проверяем, что в сервис пришли правильные данные
					require.Equal(t, name, info.Name)                // Имя должно совпадать
					require.Equal(t, email, info.Email)              // Email должен совпадать
					require.Equal(t, "user", info.Role)              // Роль должна быть в нижнем регистре
					require.NotEmpty(t, info.PasswordHash)           // Пароль должен быть захеширован (не пустой)
					require.NotEqual(t, password, info.PasswordHash) // Хеш не равен исходному паролю

					// Возвращаем успешный результат от сервиса
					return id, nil
				})
				return mock
			},
		},
		// ТЕСТ-КЕЙС 2: Ошибка сервисного слоя
		{
			name: "service error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,                                                   // Ожидаем nil вместо ответа
			err:  status.Error(codes.Internal, "failed to create user"), // Ожидаем gRPC ошибку
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				mock := serviceMocks.NewAuthServiceMock(mc)
				mock.CreateMock.Set(func(ctx context.Context, info *model.AuthInfo) (int64, error) {
					// Проверяем те же данные, что и в успешном кейсе
					require.Equal(t, name, info.Name)
					require.Equal(t, email, info.Email)
					require.Equal(t, "user", info.Role)
					require.NotEmpty(t, info.PasswordHash)
					require.NotEqual(t, password, info.PasswordHash)

					// Возвращаем ошибку от сервиса
					return 0, serviceErr
				})
				return mock
			},
		},
		// ТЕСТ-КЕЙС 3: Ошибка валидации (пароли не совпадают)
		{
			name: "validation error - passwords do not match",
			args: args{
				ctx: ctx,
				// Создаем запрос с РАЗНЫМИ паролями
				req: &desc.CreateRequest{
					Name:            name,
					Email:           email,
					Password:        password,
					PasswordConfirm: "different_password", // Не совпадает с password
					Role:            role,
				},
			},
			want: nil,                                                           // Ожидаем nil вместо ответа
			err:  status.Error(codes.InvalidArgument, "passwords do not match"), // Ожидаем ошибку валидации
			authServiceMock: func(mc *minimock.Controller) service.AuthService {
				// Создаем мок НЕ настраивая CreateMock
				// Это означает, что сервис НЕ должен быть вызван!
				mock := serviceMocks.NewAuthServiceMock(mc)
				return mock
			},
		},
	}

	// Запускаем все тест-кейсы
	for _, tt := range tests {
		tt := tt // Фиксируем переменную для параллельного запуска
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Раскомментировать для параллельного запуска подтестов

			// Создаем мок сервиса для этого тест-кейса
			authServiceMock := tt.authServiceMock(mc)

			// Создаем API с моком вместо реального сервиса
			api := auth.NewImplementation(authServiceMock)

			// Вызываем тестируемый метод
			resp, err := api.Create(tt.args.ctx, tt.args.req)

			// Проверяем результаты
			require.Equal(t, tt.err, err)   // Ошибка должна совпадать с ожидаемой
			require.Equal(t, tt.want, resp) // Ответ должен совпадать с ожидаемым
		})
	}
}
