package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/natefinch/lumberjack"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"

	"google.golang.org/grpc/reflection"

	"github.com/Denis/project_auth/internal/interceptor"
	"github.com/Denis/project_auth/internal/logger"

	desc "github.com/Denis/project_auth/pkg/user_v1"
)

var logLevel = flag.String("l", "info", "log level")

const (
	grpcPort   = 50051
	authPrefix = "Bearer "

	refreshTokenSecretKey = "W4/X+LLjehdxptt4YgGFCvMpq5ewptpZZYRHY6A72g0="
	accessTokenSecretKey  = "VqvguGiffXILza1f44TWXowDT4zwf03dtXmqWW4SYyE="

	refreshTokenExpiration = 60 * time.Minute
	accessTokenExpiration  = 2 * time.Minute
)

var accessibleRoles map[string]string

//	type serverAuth struct {
//		descAuth.UnimplementedAuthV1Server
//	}
type server struct {
	desc.UnimplementedUserV1Server
}

// func (s *serverAuth) Login(ctx context.Context, req *descAuth.LoginRequest) (*descAuth.LoginResponse, error) {
// 	// Лезем в базу или кэш за данными пользователя
// 	// Сверяем хэши пароля

// 	refreshToken, err := utils.GenerateToken(model.UserInfo{
// 		Username: req.GetUsername(),
// 		// Это пример, в реальности роль должна браться из базы или кэша
// 		Role: "admin",
// 	},
// 		[]byte(refreshTokenSecretKey),
// 		refreshTokenExpiration,
// 	)
// 	if err != nil {
// 		return nil, errors.New("failed to generate token")
// 	}

// 	return &descAuth.LoginResponse{RefreshToken: refreshToken}, nil
// }

// func (s *serverAuth) GetRefreshToken(ctx context.Context, req *descAuth.GetRefreshTokenRequest) (*descAuth.GetRefreshTokenResponse, error) {
// 	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
// 	if err != nil {
// 		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
// 	}

// 	// Можем слазать в базу или в кэш за доп данными пользователя

// 	refreshToken, err := utils.GenerateToken(model.UserInfo{
// 		Username: claims.Username,
// 		// Это пример, в реальности роль должна браться из базы или кэша
// 		Role: "admin",
// 	},
// 		[]byte(refreshTokenSecretKey),
// 		refreshTokenExpiration,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &descAuth.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
// }

// func (s *serverAuth) GetAccessToken(ctx context.Context, req *descAuth.GetAccessTokenRequest) (*descAuth.GetAccessTokenResponse, error) {
// 	claims, err := utils.VerifyToken(req.GetRefreshToken(), []byte(refreshTokenSecretKey))
// 	if err != nil {
// 		return nil, status.Errorf(codes.Aborted, "invalid refresh token")
// 	}

// 	// Можем слазать в базу или в кэш за доп данными пользователя

// 	accessToken, err := utils.GenerateToken(model.UserInfo{
// 		Username: claims.Username,
// 		// Это пример, в реальности роль должна браться из базы или кэша
// 		Role: "admin",
// 	},
// 		[]byte(accessTokenSecretKey),
// 		accessTokenExpiration,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &descAuth.GetAccessTokenResponse{AccessToken: accessToken}, nil
// }

// type serverAccess struct {
// 	descAccess.UnimplementedAccessV1Server
// }

// func (s *serverAccess) Check(ctx context.Context, req *descAccess.CheckRequest) (*emptypb.Empty, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, errors.New("metadata is not provided")
// 	}

// 	authHeader, ok := md["authorization"]
// 	if !ok || len(authHeader) == 0 {
// 		return nil, errors.New("authorization header is not provided")
// 	}

// 	if !strings.HasPrefix(authHeader[0], authPrefix) {
// 		return nil, errors.New("invalid authorization header format")
// 	}

// 	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

// 	claims, err := utils.VerifyToken(accessToken, []byte(accessTokenSecretKey))
// 	if err != nil {
// 		return nil, errors.New("access token is invalid")
// 	}

// 	accessibleMap, err := s.accessibleRoles(ctx)
// 	if err != nil {
// 		return nil, errors.New("failed to get accessible roles")
// 	}

// 	role, ok := accessibleMap[req.GetEndpointAddress()]
// 	if !ok {
// 		return &emptypb.Empty{}, nil
// 	}

// 	if role == claims.Role {
// 		return &emptypb.Empty{}, nil
// 	}

// 	return nil, errors.New("access denied")
// }

// // Возвращает мапу с адресом эндпоинта и ролью, которая имеет доступ к нему
// func (s *serverAccess) accessibleRoles(ctx context.Context) (map[string]string, error) {
// 	if accessibleRoles == nil {
// 		accessibleRoles = make(map[string]string)

// 		// Лезем в базу за данными о доступных ролях для каждого эндпоинта
// 		// Можно кэшировать данные, чтобы не лезть в базу каждый раз

// 		// Например, для эндпоинта /note_v1.NoteV1/Get доступна только роль admin
// 		accessibleRoles[model.ExamplePath] = "admin"
// 	}

// 	return accessibleRoles, nil
// }

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger.Init(getCore(getAtomicLevel()))

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
			),
		),
	)
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel() zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(*logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}
