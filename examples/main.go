package main

import (
	"context"
	"log"
	"time"

	ratelimiterv1 "github.com/sunnyjain123/go-ratelimiter/client/proto/ratelimiter/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverAddr = "localhost:50051"

func main() {
	// Create client connection
	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := ratelimiterv1.NewRateLimiterServiceClient(conn)

	// Context with timeout (VERY important)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &ratelimiterv1.PingRequest{
		Message: "hello from client",
	})
	if err != nil {
		log.Fatalf("Ping failed: %v", err)
	}

	log.Printf("Response from server: %s\n", resp.Reply)
}
