# erlc-go

Go client for the Police Roleplay Community API. 

## Features

- Complete PRC API v1 support
- Automatic retry with exponential backoff
- Built-in rate limiting and optional caching
- Full context.Context support
- Zero external dependencies

## Installation

```bash
go get github.com/joyclens/erlc-go
```

## Quick Start

```go
package main

import (
	"context"
	"log"
	"github.com/joyclens/erlc-go"
)

func main() {
	client, err := erlc.NewClient("server-key")
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

## API Methods

**Server**
```go
client.GetServer(ctx)
client.ExecuteCommand(ctx, "announce", "message")
client.ResetAPIKey(ctx)
```

**Players**
```go
client.GetPlayers(ctx)
client.GetQueue(ctx)
client.FindPlayerByUsername(ctx, "name")
client.GetStaffMembers(ctx)
```

**Logs**
```go
client.GetCommandLogs(ctx)
client.GetKillLogs(ctx)
client.GetModCalls(ctx)
```

**Records**
```go
client.GetBans(ctx)
client.GetVehicles(ctx)
client.GetStaff(ctx)
```

## Configuration

```go
client, err := erlc.NewClient(
	"server-key",
	erlc.WithTimeout(45*time.Second),
	erlc.WithMaxRetries(5),
	erlc.WithRateLimiting(true, 50),
	erlc.WithCaching(true, 5*time.Minute, 10*1024*1024),
)
```

## Error Handling

```go
resp, err := client.GetPlayers(ctx)
if err != nil {
	if apiErr, ok := err.(*erlc.APIError); ok {
		if apiErr.IsRateLimit() {
			// handle rate limit
		}
	}
}
```

## Caching

Enable caching:
```go
erlc.WithCaching(true, 5*time.Minute, 10*1024*1024)
```

Per-request control:
```go
players, err := client.GetPlayersWithCache(ctx, &erlc.CacheOptions{
	Enabled: true,
	TTL:     2*time.Minute,
})

client.ClearCache()
stats := client.CacheStats()
```

## Testing

```bash
go test -v ./...
go test -cover ./...
```

## Examples

See [examples/](examples/) directory.

```bash
go run examples/basic_usage.go
go run examples/advanced_usage.go
```

## License

Apache License 2.0. See [LICENSE](LICENSE).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Support

- [PRC API Documentation](https://apidocs.policeroleplay.community/)
- Report issues on GitHub
