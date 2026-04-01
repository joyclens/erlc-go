package erlc

import (
	"context"
	"net/http"
	"time"
)

func (c *Client) GetPlayers(ctx context.Context) (*PlayerList, error) {
	var result PlayerList
	err := c.doWithCache(ctx, EndpointPlayers, &result, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPlayersWithCache(ctx context.Context, cacheOpt *CacheOptions) (*PlayerList, error) {
	var result PlayerList
	err := c.doWithCache(ctx, EndpointPlayers, &result, cacheOpt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetQueue(ctx context.Context) (*QueueList, error) {
	var result QueueList
	err := c.doWithCache(ctx, EndpointQueue, &result, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetQueueWithCache(ctx context.Context, cacheOpt *CacheOptions) (*QueueList, error) {
	var result QueueList
	err := c.doWithCache(ctx, EndpointQueue, &result, cacheOpt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetJoinLogs(ctx context.Context) (*JoinLogList, error) {
	var result JoinLogList
	err := c.do(ctx, http.MethodGet, EndpointJoinLogs, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) FilterPlayersByTeam(ctx context.Context, team string) ([]Player, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Player
	for _, p := range players.Players {
		if p.Team == team {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (c *Client) FindPlayerByUsername(ctx context.Context, username string) (*Player, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range players.Players {
		if p.Username == username {
			return &p, nil
		}
	}
	return nil, nil
}

func (c *Client) FindPlayerByRobloxID(ctx context.Context, robloxID int64) (*Player, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range players.Players {
		if p.RobloxID == robloxID {
			return &p, nil
		}
	}
	return nil, nil
}

func (c *Client) GetStaffMembers(ctx context.Context) ([]Player, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	var staff []Player
	for _, p := range players.Players {
		if p.IsStaff {
			staff = append(staff, p)
		}
	}
	return staff, nil
}

func (c *Client) GetPlayersWithPermission(ctx context.Context, permission string) ([]Player, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Player
	for _, p := range players.Players {
		for _, perm := range p.Permissions {
			if perm == permission {
				filtered = append(filtered, p)
				break
			}
		}
	}
	return filtered, nil
}

func (c *Client) GetPlayerCount(ctx context.Context) (int, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return 0, err
	}
	return players.Count, nil
}

func (c *Client) GetQueueCount(ctx context.Context) (int, error) {
	queue, err := c.GetQueue(ctx)
	if err != nil {
		return 0, err
	}
	return queue.Count, nil
}

func (c *Client) GetPlayerStats(ctx context.Context) (map[string]interface{}, error) {
	players, err := c.GetPlayers(ctx)
	if err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total":       players.Count,
		"staff":       0,
		"moderators":  0,
		"admins":      0,
		"owner":       0,
	}

	for _, p := range players.Players {
		if p.IsOwner {
			stats["owner"] = true
		} else if p.IsAdministrator {
			stats["admins"] = stats["admins"].(int) + 1
		} else if p.IsModerator {
			stats["moderators"] = stats["moderators"].(int) + 1
		} else if p.IsStaff {
			stats["staff"] = stats["staff"].(int) + 1
		}
	}

	return stats, nil
}

func (c *Client) WaitForPlayerJoin(ctx context.Context, robloxID int64, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			player, err := c.FindPlayerByRobloxID(ctx, robloxID)
			if err != nil {
				return err
			}
			if player != nil {
				return nil
			}
		}
	}
}
