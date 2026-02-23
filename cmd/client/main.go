// package main

// import (
// 	"context"
// 	"log"
// 	"time"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials"

// 	"google.golang.org/protobuf/types/known/wrapperspb"

// 	desc "github.com/Denis/project_auth/pkg/user_v1"
// )

// const (
// 	address = "localhost:50051"
// )

// func main() {
// 	// Настройка подключения (выберите один вариант)

// 	// Вариант 1: С TLS (если сервер с TLS)
// 	creds, err := credentials.NewClientTLSFromFile("../../service.pem", "")
// 	if err != nil {
// 		log.Fatalf("failed to load TLS credentials: %v", err)
// 	}
// 	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))

// 	// Вариант 2: Без TLS (для разработки)
// 	// conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	if err != nil {
// 		log.Fatalf("failed to connect: %v", err)
// 	}
// 	defer conn.Close()

// 	// Создаем клиент
// 	client := desc.NewUserV1Client(conn)

// 	// Контекст с таймаутом
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// 1. СОЗДАНИЕ ПОЛЬЗОВАТЕЛЯ
// 	log.Println("=== 1. Creating user ===")
// 	createResp, err := client.Create(ctx, &desc.CreateRequest{
// 		Name:            "JohnDoe",
// 		Email:           "john@example.com",
// 		Password:        "secret123",
// 		PasswordConfirm: "secret123",
// 		Role:            desc.Role_ADMIN,
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to create user: %v", err)
// 	}
// 	log.Printf("✅ User created with ID: %d\n", createResp.GetId())

// 	// Небольшая пауза
// 	time.Sleep(time.Second)

// 	// 2. ПОЛУЧЕНИЕ ПОЛЬЗОВАТЕЛЯ
// 	log.Println("=== 2. Getting user ===")
// 	getResp, err := client.Get(ctx, &desc.GetRequest{
// 		Id: createResp.GetId(),
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to get user: %v", err)
// 	}
// 	log.Printf("✅ Got user: ID=%d, Name=%s, Email=%s, Role=%v\n",
// 		getResp.GetId(),
// 		getResp.GetName(),
// 		getResp.GetEmail(),
// 		getResp.GetRole())

// 	// Небольшая пауза
// 	time.Sleep(time.Second)

// 	// 3. ОБНОВЛЕНИЕ ПОЛЬЗОВАТЕЛЯ
// 	log.Println("=== 3. Updating user ===")
// 	_, err = client.Update(ctx, &desc.UpdateRequest{
// 		Id:    createResp.GetId(),
// 		Name:  wrapperspb.String("JaneDoe"), // 👈 Правильно для StringValue
// 		Email: wrapperspb.String("jane@example.com"),
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to update user: %v", err)
// 	}
// 	log.Println("✅ User updated")

// 	// Небольшая пауза
// 	time.Sleep(time.Second)

// 	// 4. ПОЛУЧЕНИЕ ОБНОВЛЁННОГО ПОЛЬЗОВАТЕЛЯ
// 	log.Println("=== 4. Getting updated user ===")
// 	getResp2, err := client.Get(ctx, &desc.GetRequest{
// 		Id: createResp.GetId(),
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to get updated user: %v", err)
// 	}
// 	log.Printf("✅ Got updated user: Name=%s, Email=%s\n",
// 		getResp2.GetName(),
// 		getResp2.GetEmail())

// 	// Небольшая пауза
// 	time.Sleep(time.Second)

// 	// 5. УДАЛЕНИЕ ПОЛЬЗОВАТЕЛЯ (опционально)
// 	log.Println("=== 5. Deleting user ===")
// 	_, err = client.Delete(ctx, &desc.DeleteRequest{
// 		Id: createResp.GetId(),
// 	})
// 	if err != nil {
// 		log.Fatalf("failed to delete user: %v", err)
// 	}
// 	log.Println("✅ User deleted")

// 	// 6. ПРОВЕРКА УДАЛЕНИЯ
// 	log.Println("=== 6. Checking deletion ===")
// 	_, err = client.Get(ctx, &desc.GetRequest{
// 		Id: createResp.GetId(),
// 	})
// 	if err != nil {
// 		log.Printf("✅ Expected error after deletion: %v\n", err)
// 	} else {
// 		log.Println("❌ User still exists!")
// 	}

//		log.Println("🎉 All operations completed successfully!")
//	}
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/Denis/project_auth/internal/model"
	descAccess "github.com/Denis/project_auth/pkg/access_v1"
)

var accessToken = flag.String("a", "", "access token")

const servicePort = 50051

func main() {
	flag.Parse()

	ctx := context.Background()
	md := metadata.New(map[string]string{"Authorization": "Bearer " + *accessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	conn, err := grpc.Dial(
		fmt.Sprintf(":%d", servicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to dial GRPC client: %v", err)
	}

	cl := descAccess.NewAccessV1Client(conn)

	_, err = cl.Check(ctx, &descAccess.CheckRequest{
		EndpointAddress: model.ExamplePath,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Access granted")
}
