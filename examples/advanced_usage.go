package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joyclens/erlc-go"
)

type CustomLogger struct{}

func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	fmt.Printf("[DEBUG] %s\n", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[INFO] %s\n", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Warnf(format string, args ...interface{}) {
	fmt.Printf("[WARN] %s\n", fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] %s\n", fmt.Sprintf(format, args...))
}

func main() {
	serverKey := "your-server-key-here"

	client, err := erlc.NewClient(
		serverKey,
		erlc.WithTimeout(45*time.Second),
		erlc.WithMaxRetries(5),
		erlc.WithRateLimiting(true, 50),
		erlc.WithCaching(true, 10*time.Minute, 5*1024*1024),
		erlc.WithLogging(true, &CustomLogger{}),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	fmt.Println("Caching Example")
	players1, err := client.GetPlayers(ctx)
	if err != nil {
		log.Fatalf("Failed to get players: %v", err)
	}
	fmt.Printf("First call: %d players (from API)\n", players1.Count)

	start := time.Now()
	players2, err := client.GetPlayers(ctx)
	if err != nil {
		log.Fatalf("Failed to get players: %v", err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Second call: %d players (cached in %v)\n", players2.Count, elapsed)

	stats := client.CacheStats()
	if stats != nil {
		fmt.Printf("Cache stats - Hits: %d, Misses: %d, Size: %d/%d bytes\n",
			stats.Hits, stats.Misses, stats.Size, stats.MaxSize)
	}

	fmt.Println("\nError Handling Example")
	_, err = client.ExecuteCommand(ctx, "announce", "Server maintenance!")
	if err != nil {
		if apiErr, ok := err.(*erlc.APIError); ok {
			if apiErr.IsRateLimit() {
				fmt.Printf("Rate limited! Retry after %v\n", apiErr.RetryAfter)
			} else if apiErr.IsAuthError() {
				fmt.Println("Authentication failed")
			} else if apiErr.IsBadRequest() {
				fmt.Printf("Bad request: %s\n", apiErr.Message)
			} else if apiErr.IsServerError() {
				fmt.Println("Server error")
			}
		}
	}

	fmt.Println("\nFiltering Examples")
	policeTeam, err := client.FilterPlayersByTeam(ctx, "police")
	if err != nil {
		log.Fatalf("Failed to filter: %v", err)
	}
	fmt.Printf("Police team members: %d\n", len(policeTeam))

	staff, err := client.GetStaffMembers(ctx)
	if err != nil {
		log.Fatalf("Failed to get staff: %v", err)
	}
	fmt.Printf("Staff members online: %d\n", len(staff))

	fmt.Println("\nBan Analysis")
	permanent, err := client.GetPermanentBans(ctx)
	if err != nil {
		log.Fatalf("Failed to get permanent bans: %v", err)
	}
	fmt.Printf("Permanent bans: %d\n", len(permanent))

	temporary, err := client.GetTemporaryBans(ctx)
	if err != nil {
		log.Fatalf("Failed to get temporary bans: %v", err)
	}
	fmt.Printf("Temporary bans: %d\n", len(temporary))

	fmt.Println("\nCache Control")
	customCache := &erlc.CacheOptions{
		Enabled: true,
		TTL:     2 * time.Minute,
	}
	vehicles, err := client.GetVehiclesWithCache(ctx, customCache)
	if err != nil {
		log.Fatalf("Failed to get vehicles: %v", err)
	}
	fmt.Printf("Vehicles (cached 2 minutes): %d\n", vehicles.Count)

	noCache := &erlc.CacheOptions{
		Enabled: false,
	}
	queue, err := client.GetQueueWithCache(ctx, noCache)
	if err != nil {
		log.Fatalf("Failed to get queue: %v", err)
	}
	fmt.Printf("Queue (live data): %d\n", queue.Count)

	client.ClearCache()
	fmt.Println("Cache cleared")
}
