# go-ratelimiter

A production-oriented **gRPC Rate Limiter service** written in Go, designed to demonstrate clean service architecture, middleware-based interception, and pluggable rate-limiting strategies.

This repository was built step-by-step with an interview-style machine-coding mindset, focusing on **clarity, extensibility, and correctness** rather than shortcuts.

---

## ✨ Features

* gRPC-based Rate Limiter service
* Unary interceptor–based middleware
* Pluggable rate limiter strategies
* Per-RPC rate limiting configuration
* Environment-aware configuration (dev / prod)
* YAML-based config using Viper (typed + untyped access)
* Clean separation of concerns (API / middleware / business logic)

---

## 🧱 Architecture Overview

```
config/              # Config loading (Viper + typed structs)
  config.yaml
  config.go
internal/
  limiter/           # Rate limiter interfaces + implementations
pkg/
  middleware/        # gRPC interceptors
client/
  proto/             # Generated protobuf code
proto/
  ratelimiter/v1/    # Protobuf definitions
service/
  ratelimiter/       # gRPC service implementation
main.go              # entry point
```

### High-level flow

```
[gRPC Client]
     ↓
[gRPC Unary Interceptor]
     ↓
[Rate Limiter Middleware]
     ↓
[Limiter Implementation]
     ↓
[gRPC Service Handler]
```

---

## 🚀 Getting Started

### Prerequisites

* Go 1.21+
* Protocol Buffers
* `protoc-gen-go`
* `protoc-gen-go-grpc`

---

## 🧬 Protobuf

### Generate protos

```bash
make proto
```

Generated files live under:

```
client/proto/ratelimiter/v1/
```

---

## ▶️ Running the Server

From repo root:

```bash
CONFIG_PATH=. go run main.go
```

The server starts a gRPC service exposing:

```
ratelimiter.v1.RateLimiterService
```

---

## 📡 Calling the Service (Go client example)

Reflection is intentionally **disabled**.
Service-to-service communication is expected to use generated stubs.

```go
conn, _ := grpc.NewClient(
    "localhost:50051",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
client := pb.NewRateLimiterServiceClient(conn)

resp, err := client.Ping(ctx, &pb.PingRequest{})
```

---

## 🧠 Rate Limiter Design

### Interface

```go
type RateLimiter interface {
    Allow(ctx context.Context, key string) (bool, error)
}
```

### Implementations

* **Fixed Window Limiter**
* **Leaky Bucket Limiter**

Limiters are:

* Stateless from the middleware perspective
* Config-driven
* Easily replaceable

---

## ⚙️ Configuration

Configuration is loaded from `config.yaml` using **Viper**.

### Example

```yaml
app:
  env: dev

rate_limiter:
  default:
    type: leaky_bucket
    capacity: 10
    refill_rate: 1

  rpc_limits:
    /ratelimiter.v1.RateLimiterService/Ping:
      type: fixed_window
      requests: 5
      window_seconds: 10
```

### Access patterns

* Typed structs for core config
* Dynamic Viper access for RPC-level overrides

---

## 🧩 Middleware

Rate limiting is enforced via a **gRPC Unary Interceptor**:

* Extracts RPC full method name
* Resolves limiter config (RPC → default)
* Applies limiter before handler execution
* Returns `ResourceExhausted` on limit breach

---

## 🧪 Environments

Environment is controlled via config:

```yaml
app:
  env: prod
```

Behavior such as logging or strictness can be toggled per environment.

---

## 🧠 Design Decisions

* ❌ No gRPC reflection (explicit contracts only)
* ❌ No global variables for limiters
* ✅ Middleware-first enforcement
* ✅ Internal packages for business logic
* ✅ Clean separation of API vs core logic

---

## 📌 What This Repo Is Meant To Show

* gRPC service design
* Middleware-based cross-cutting concerns
* Config-driven behavior
* Interview-ready Go code structure
* Real-world operational considerations

---

## 📜 License

MIT

---

## 👤 Author

Sunny Jain
GitHub: [https://github.com/sunnyjain123](https://github.com/sunnyjain123)
