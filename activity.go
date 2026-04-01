package erlc

import (
	"context"
	"net/http"
)

func (c *Client) GetCommandLogs(ctx context.Context) (*CommandLogList, error) {
	var result CommandLogList
	err := c.do(ctx, http.MethodGet, EndpointCommandLogs, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetKillLogs(ctx context.Context) (*KillLogList, error) {
	var result KillLogList
	err := c.do(ctx, http.MethodGet, EndpointKillLogs, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetModCalls(ctx context.Context) (*ModCallList, error) {
	var result ModCallList
	err := c.do(ctx, http.MethodGet, EndpointModCalls, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) FilterCommandLogsByPlayer(ctx context.Context, username string) ([]CommandLog, error) {
	logs, err := c.GetCommandLogs(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []CommandLog
	for _, log := range logs.Logs {
		if log.Username == username {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

func (c *Client) FilterCommandLogsByCommand(ctx context.Context, command string) ([]CommandLog, error) {
	logs, err := c.GetCommandLogs(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []CommandLog
	for _, log := range logs.Logs {
		if log.Command == command {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

func (c *Client) FilterKillLogsByKiller(ctx context.Context, killerName string) ([]KillLog, error) {
	logs, err := c.GetKillLogs(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []KillLog
	for _, log := range logs.Logs {
		if log.KillerName == killerName {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

func (c *Client) FilterKillLogsByVictim(ctx context.Context, victimName string) ([]KillLog, error) {
	logs, err := c.GetKillLogs(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []KillLog
	for _, log := range logs.Logs {
		if log.VictimName == victimName {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

func (c *Client) FilterKillLogsByWeapon(ctx context.Context, weapon string) ([]KillLog, error) {
	logs, err := c.GetKillLogs(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []KillLog
	for _, log := range logs.Logs {
		if log.Weapon == weapon {
			filtered = append(filtered, log)
		}
	}
	return filtered, nil
}

func (c *Client) GetHeadshotCount(ctx context.Context) (int, error) {
	logs, err := c.GetKillLogs(ctx)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, log := range logs.Logs {
		if log.Headshot {
			count++
		}
	}
	return count, nil
}

func (c *Client) GetResolvedModCalls(ctx context.Context) ([]ModCall, error) {
	calls, err := c.GetModCalls(ctx)
	if err != nil {
		return nil, err
	}

	var resolved []ModCall
	for _, call := range calls.Calls {
		if call.Status == "resolved" {
			resolved = append(resolved, call)
		}
	}
	return resolved, nil
}

func (c *Client) GetOpenModCalls(ctx context.Context) ([]ModCall, error) {
	calls, err := c.GetModCalls(ctx)
	if err != nil {
		return nil, err
	}

	var open []ModCall
	for _, call := range calls.Calls {
		if call.Status != "resolved" {
			open = append(open, call)
		}
	}
	return open, nil
}

func (c *Client) GetModCallsByPlayer(ctx context.Context, username string) ([]ModCall, error) {
	calls, err := c.GetModCalls(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []ModCall
	for _, call := range calls.Calls {
		if call.Username == username {
			filtered = append(filtered, call)
		}
	}
	return filtered, nil
}
