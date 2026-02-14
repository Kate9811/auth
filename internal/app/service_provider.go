package app

import (
	"context"
	"log"

	"github.com/Denis/project_auth/internal/api/auth"
	"github.com/Denis/project_auth/internal/client/db"
	"github.com/Denis/project_auth/internal/client/db/pg"

	"github.com/Denis/project_auth/internal/client/db/transaction"
	"github.com/Denis/project_auth/internal/closer"
	"github.com/Denis/project_auth/internal/config"
	"github.com/Denis/project_auth/internal/repository"
	authRepository "github.com/Denis/project_auth/internal/repository/auth"
	"github.com/Denis/project_auth/internal/service"
	authService "github.com/Denis/project_auth/internal/service/auth"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Структура serviceProvider
// Создает пустой контейнер зависимостей.
//
//	Все зависимости будут инициализированы лениво (по требованию).
type serviceProvider struct {
	// Конфигурации
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	// Клиенты/подключения
	pgPool    *pgxpool.Pool
	dbClient  db.Client    // ← ДОБАВИТЬ если используешь
	txManager db.TxManager // ← ДОБАВИТЬ для транзакций

	// Репозитории
	authRepository repository.AuthRepository

	// Сервисы
	authService service.AuthService

	// gRPC реализации
	authImpl *auth.Implementation
}

// Конструктор
func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// Конфигурация для ПОДКЛЮЧЕНИЯ К БД
func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}
		s.pgConfig = cfg
	}
	return s.pgConfig
}

//Тут создается и возвращается конфигурация для gRPC сервера

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	// Если репозиторий ещё не создан → создаём
	if s.authRepository == nil {
		//Создание репозитория для работы с пользователями, который использует пул соединений с БД
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx)) // s.TxManager(ctx)

	}
	return s.authService
}
func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
