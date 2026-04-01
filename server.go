package erlc

import (
	"context"
	"net/http"
)

func (c *Client) GetServer(ctx context.Context) (*Server, error) {
	var result Server
	err := c.do(ctx, http.MethodGet, EndpointServer, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ExecuteCommand(ctx context.Context, command string, args string) (*CommandResponse, error) {
	if command == "" {
		return nil, NewValidationError("command", "command is required", nil)
	}

	req := CommandRequest{
		Command: command,
		Args:    args,
	}

	var result CommandResponse
	err := c.do(ctx, http.MethodPost, EndpointServerCommand, req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) ResetAPIKey(ctx context.Context) error {
	var result map[string]interface{}
	return c.do(ctx, http.MethodPost, EndpointAPIKeyReset, nil, &result)
}
