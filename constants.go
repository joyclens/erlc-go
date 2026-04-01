package erlc

import "time"

const (
	BaseURL        = "https://api.policeroleplay.community"
	V1BaseURL      = BaseURL + "/v1"
	V2BaseURL      = BaseURL + "/v2"
	ServerKeyHeader = "Server-Key"
	ContentTypeHeader = "Content-Type"
	ContentTypeJSON = "application/json"
	UserAgentHeader = "User-Agent"
	UserAgent       = "erlc-go/1.0"
	MaxRetries      = 3
	DefaultTimeout  = 30 * time.Second
	DefaultRateLimitPerSecond = 100
	DefaultCacheTTL = 5 * time.Minute
	DefaultMaxCacheSize = 10 * 1024 * 1024
)

const (
	EndpointServer        = "/server"
	EndpointServerCommand = "/server/command"
	EndpointAPIKeyReset   = "/api-key/reset"
	EndpointPlayers       = "/server/players"
	EndpointQueue         = "/server/queue"
	EndpointJoinLogs      = "/server/joinlogs"
	EndpointCommandLogs   = "/server/commandlogs"
	EndpointKillLogs      = "/server/killlogs"
	EndpointModCalls      = "/server/modcalls"
	EndpointBans          = "/server/bans"
	EndpointVehicles      = "/server/vehicles"
	EndpointStaff         = "/server/staff"
)

const (
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusUnprocessableEntity = 422
	StatusTooManyRequests     = 429
	StatusInternalServerError = 500
	StatusServiceUnavailable  = 503
)
