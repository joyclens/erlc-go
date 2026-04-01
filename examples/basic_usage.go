package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joyclens/erlc-go"
)

func main() {
	serverKey := "your-server-key-here"

	client, err := erlc.NewClient(serverKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	server, err := client.GetServer(ctx)
	if err != nil {
		log.Fatalf("Failed to get server info: %v", err)
	}
	fmt.Printf("Server: %s\n", server.Name)
	fmt.Printf("Owner: %s\n", server.Owner)
	fmt.Printf("Players: %d/%d\n", server.PlayerCount, server.MaxPlayerCount)

	players, err := client.GetPlayers(ctx)
	if err != nil {
		log.Fatalf("Failed to get players: %v", err)
	}

	fmt.Printf("\nOnline Players: %d\n", players.Count)
	for _, player := range players.Players {
		status := ""
		if player.IsOwner {
			status = "[Owner]"
		} else if player.IsAdministrator {
			status = "[Admin]"
		} else if player.IsModerator {
			status = "[Mod]"
		} else if player.IsStaff {
			status = "[Staff]"
		}

		fmt.Printf("  %s (%d) %s Team: %s\n", player.Username, player.RobloxID, status, player.Team)
	}

	queue, err := client.GetQueue(ctx)
	if err != nil {
		log.Fatalf("Failed to get queue: %v", err)
	}
	fmt.Printf("\nQueued Players: %d\n", queue.Count)

	vehicles, err := client.GetVehicles(ctx)
	if err != nil {
		log.Fatalf("Failed to get vehicles: %v", err)
	}

	fmt.Printf("\nVehicles: %d\n", vehicles.Count)
	for _, vehicle := range vehicles.Vehicles {
		fmt.Printf("  %s (%s) - License: %s - Owner: %s\n",
			vehicle.Name, vehicle.Model, vehicle.License, vehicle.OwnerName)
	}

	kills, err := client.GetKillLogs(ctx)
	if err != nil {
		log.Fatalf("Failed to get kill logs: %v", err)
	}

	fmt.Printf("\nRecent Kills (up to 5):\n")
	for i, kill := range kills.Logs {
		if i >= 5 {
			break
		}
		headshot := ""
		if kill.Headshot {
			headshot = " [HEADSHOT]"
		}
		fmt.Printf("  %s killed %s with %s at %.0fm%s\n",
			kill.KillerName, kill.VictimName, kill.Weapon, kill.Distance, headshot)
	}
}
