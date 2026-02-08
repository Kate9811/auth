package main

import (
	"context"
	"log"
	"time"

	desc "github.com/Denis/project_auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	c := desc.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 1. Create
	createResp, err := c.Create(ctx, &desc.CreateRequest{
		Name:            "Test User",
		Email:           "test@test.com",
		Password:        "12345678",
		PasswordConfirm: "12345678",
		Role:            desc.Role_USER,
	})
	if err != nil {
		log.Println("Create error:", err)
	} else {
		log.Printf("Created user with ID: %d", createResp.GetId())
	}

	// 2. Get
	getresp, err := c.Get(ctx, &desc.GetRequest{
		Id: 1,
	})
	if err != nil {
		log.Println("Get error:", err)
	} else {
		log.Printf("Got user: ID=%d, Name=%s, Email=%s, Role=%v",
			getresp.GetId(),
			getresp.GetName(),
			getresp.GetEmail(),
			getresp.GetRole())
	}

	// 3. Update
	_, err = c.Update(ctx, &desc.UpdateRequest{ // <-- убрал :=, используй =
		Id:    1,
		Name:  &wrapperspb.StringValue{Value: "Петр Петров"}, // <-- убрал кавычки вокруг всей структуры
		Email: &wrapperspb.StringValue{Value: "petr@example.com"},
	})
	if err != nil {
		log.Println("Update error:", err)
	} else {
		log.Println("✅ Update worked")
	}

	// 4. Delete
	_, err = c.Delete(ctx, &desc.DeleteRequest{Id: 1})
	if err != nil {
		log.Println("Delete error:", err)
	} else {
		log.Println("✅ User deleted")
	}
}
