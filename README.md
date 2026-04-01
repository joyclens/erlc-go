# erlc-go

Go client for the Police Roleplay Community API. Build server management tools for ER:LC with type safety, automatic retry logic, rate limiting, and optional caching.

## Features

- Complete PRC API v1 support
- Type-safe request/response structures
- Automatic retry with exponential backoff
- Built-in rate limiting (configurable)
- Optional in-memory caching with TTL
- Full context.Context support
- Zero external dependencies (stdlib only)
- Production-ready with comprehensive error handling

## Installation

```bash
go get github.com/joyclens/erlc-go
```

Requires Go 1.21 or later.

## Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/joyclens/erlc-go"
)

func main() {
	client, err := erlc.NewClient("your-server-key")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	server, err := client.GetServer(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	println(server.Name)
}
```

## API Reference

### Server Management

```go
server, err := client.GetServer(ctx)
err := client.ExecuteCommand(ctx, "announce", "message")
err := client.ResetAPIKey(ctx)
```

### Players

```go
players, err := client.GetPlayers(ctx)
queue, err := client.GetQueue(ctx)
logs, err := client.GetJoinLogs(ctx)

player, err := client.FindPlayerByUsername(ctx, "username")
staff, err := client.GetStaffMembers(ctx)
count, err := client.GetPlayerCount(ctx)
```

### Activity Logs

```go
commands, err := client.GetCommandLogs(ctx)
kills, err := client.GetKillLogs(ctx)
calls, err := client.GetModCalls(ctx)

filtered, err := client.FilterKillLogsByKiller(ctx, "name")
headshots, err := client.GetHeadshotCount(ctx)
```

### Records

```go
bans, err := client.GetBans(ctx)
vehicles, err := client.GetVehicles(ctx)
staff, err := client.GetStaff(ctx)

active, err := client.GetActiveBans(ctx)
permanent, err := client.GetPermanentBans(ctx)
stolen, err := client.GetStolenVehicles(ctx)
```

## Configuration

### Basic Options

```go
client, err := erlc.NewClient(
	"server-key",
	erlc.WithTimeout(45*time.Second),
	erlc.WithMaxRetries(5),
)
```

### Rate Limiting

Enabled by default at 100 requests/second:

```go
erlc.WithRateLimiting(true, 50)
```

### Caching

Disabled by default. Enable with TTL and size limit:

```go
erlc.WithCaching(true, 5*time.Minute, 10*1024*1024)
```

Per-request control:

```go
players, err := client.GetPlayersWithCache(ctx, &erlc.CacheOptions{
	Enabled: true,
	TTL:     2 * time.Minute,
})

client.ClearCache()
stats := client.CacheStats()
```

### Logging

```go
client, err := erlc.NewClient(
	"server-key",
	erlc.WithLogging(true, &CustomLogger{}),
)
```

Implement the Logger interface:

```go
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
```

## Error Handling

```go
resp, err := client.GetPlayers(ctx)
if err != nil {
	if apiErr, ok := err.(*erlc.APIError); ok {
		if apiErr.IsRateLimit() {
			// Handle rate limiting
		} else if apiErr.IsAuthError() {
			// Handle auth error
		} else if apiErr.IsServerError() {
			// Handle server error
		}
	} else if netErr, ok := err.(*erlc.NetworkError); ok {
		// Handle network error
	}
}
```

## Context Support

All API calls support context cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

players, err := client.GetPlayers(ctx)
```

## Examples

See the [examples](examples/) directory for complete implementations.

```bash
go run examples/basic_usage.go
go run examples/advanced_usage.go
```

## Testing

```bash
go test -v ./...
go test -cover ./...
```

## Performance

- Connection pooling enabled
- Minimal memory overhead
- Efficient cache eviction
- Token bucket rate limiting

## Architecture

```
erlc-go/
├── client.go        # HTTP client
├── types.go         # Data structures
├── errors.go        # Error types
├── options.go       # Configuration
├── constants.go     # Endpoints and constants
├── server.go        # Server endpoints
├── player.go        # Player endpoints
├── activity.go      # Log endpoints
├── records.go       # Record endpoints
├── ratelimit.go     # Rate limiting
├── cache.go         # Caching layer
├── internal/        # Utilities
├── examples/        # Example programs
└── go.mod           # Module definition
```

## Requirements

- Go 1.21+
- Standard library only

## License

Licensed under Apache License 2.0. See [LICENSE](LICENSE).

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

- [PRC API Docs](https://apidocs.policeroleplay.community/)
- GitHub Issues for bugs
- GitHub Discussions for questions

## Disclaimer

This library is not affiliated with Roblox or the Police Roleplay Community. Use in accordance with their terms of service.
