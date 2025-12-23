package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"google.golang.org/grpc"

	"github.com/spf13/viper"
	ratelimiterv1 "github.com/sunnyjain123/go-ratelimiter/client/proto/ratelimiter/v1"
	"github.com/sunnyjain123/go-ratelimiter/config"
	"github.com/sunnyjain123/go-ratelimiter/internal/ratelimiter"
	"github.com/sunnyjain123/go-ratelimiter/pkg/middleware"
	grpcapi "github.com/sunnyjain123/go-ratelimiter/service/ratelimiter"
)

const grpcPort = ":50051"

func main() {
	config.LoadViperConfig()
	// defaultLimiter, err := ratelimiter.BuildLimiterFromViper(
	// 	"rate_limiter.default",
	// )
	// if err != nil {
	// 	log.Fatalf("failed to build default limiter: %v", err)
	// }

	cfg, _ := config.Load("config/config.yaml")
	defaultLimiter, _ := ratelimiter.BuildLimiter(cfg.RateLimiter.Default)

	rpcLimiters := map[string]ratelimiter.Limiter{}
	// for method, rule := range cfg.RateLimiter.RPCLimits {
	// 	limiter, _ := ratelimiter.BuildLimiter(rule)
	// 	rpcLimiters[method] = limiter
	// }

	fmt.Println(viper.AllKeys())
	for _, key := range viper.AllKeys() {
		if strings.HasPrefix(key, "rate_limiter.rpc_limits.") &&
			strings.HasSuffix(key, ".type") {

			method := strings.TrimSuffix(
				strings.TrimPrefix(key, "rate_limiter.rpc_limits."),
				".type",
			)

			limiter, _ := ratelimiter.BuildLimiterFromViper(
				"rate_limiter.rpc_limits." + method,
			)
			rpcLimiters["/"+method] = limiter
		}
	}

	fmt.Println(rpcLimiters)

	rlService := ratelimiter.NewService(defaultLimiter, rpcLimiters)

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.RateLimiterUnaryInterceptor(rlService)),
	}

	grpcServer := grpc.NewServer(opts...)
	registerServices(grpcServer)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go func() {
		log.Printf("gRPC server listening on %s\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	waitForShutdown(grpcServer)
}

func registerServices(grpcServer *grpc.Server) {
	ratelimiterv1.RegisterRateLimiterServiceServer(
		grpcServer,
		grpcapi.NewRateLimiterService(),
	)
}

func waitForShutdown(grpcServer *grpc.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	log.Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
}
