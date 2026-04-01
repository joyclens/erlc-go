package erlc

import "time"

type Server struct {
	ID             string    `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Owner          string    `json:"owner,omitempty"`
	PlayerCount    int       `json:"playerCount,omitempty"`
	MaxPlayerCount int       `json:"maxPlayerCount,omitempty"`
	Status         string    `json:"status,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty"`
}

type CommandRequest struct {
	Command string `json:"command"`
	Args    string `json:"args,omitempty"`
}

type CommandResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Output  interface{} `json:"output,omitempty"`
}

type Player struct {
	ID              string    `json:"id,omitempty"`
	RobloxID        int64     `json:"robloxId,omitempty"`
	Username        string    `json:"username,omitempty"`
	DisplayName     string    `json:"displayName,omitempty"`
	Team            string    `json:"team,omitempty"`
	Rank            string    `json:"rank,omitempty"`
	JoinTime        time.Time `json:"joinTime,omitempty"`
	PlayTime        int       `json:"playTime,omitempty"`
	IsStaff         bool      `json:"isStaff,omitempty"`
	IsModerator     bool      `json:"isModerator,omitempty"`
	IsAdministrator bool      `json:"isAdministrator,omitempty"`
	IsOwner         bool      `json:"isOwner,omitempty"`
	Permissions     []string  `json:"permissions,omitempty"`
}

type PlayerQueue struct {
	ID       string    `json:"id,omitempty"`
	RobloxID int64     `json:"robloxId,omitempty"`
	Username string    `json:"username,omitempty"`
	QueuePos int       `json:"queuePos,omitempty"`
	QueuedAt time.Time `json:"queuedAt,omitempty"`
}

type JoinLog struct {
	ID            string    `json:"id,omitempty"`
	RobloxID      int64     `json:"robloxId,omitempty"`
	Username      string    `json:"username,omitempty"`
	Timestamp     time.Time `json:"timestamp,omitempty"`
	Status        string    `json:"status,omitempty"`
	JoinCount     int       `json:"joinCount,omitempty"`
	TotalPlayTime int       `json:"totalPlayTime,omitempty"`
}

type CommandLog struct {
	ID        string    `json:"id,omitempty"`
	RobloxID  int64     `json:"robloxId,omitempty"`
	Username  string    `json:"username,omitempty"`
	Command   string    `json:"command,omitempty"`
	Arguments string    `json:"arguments,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
	Success   bool      `json:"success,omitempty"`
	Output    string    `json:"output,omitempty"`
}

type KillLog struct {
	ID         string    `json:"id,omitempty"`
	KillerID   int64     `json:"killerId,omitempty"`
	KillerName string    `json:"killerName,omitempty"`
	VictimID   int64     `json:"victimId,omitempty"`
	VictimName string    `json:"victimName,omitempty"`
	Weapon     string    `json:"weapon,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Location   string    `json:"location,omitempty"`
	Distance   float64   `json:"distance,omitempty"`
	Headshot   bool      `json:"headshot,omitempty"`
}

type ModCall struct {
	ID         string     `json:"id,omitempty"`
	RobloxID   int64      `json:"robloxId,omitempty"`
	Username   string     `json:"username,omitempty"`
	Message    string     `json:"message,omitempty"`
	Timestamp  time.Time  `json:"timestamp,omitempty"`
	Status     string     `json:"status,omitempty"`
	ResolvedBy string     `json:"resolvedBy,omitempty"`
	ResolvedAt *time.Time `json:"resolvedAt,omitempty"`
}

type Ban struct {
	ID        string     `json:"id,omitempty"`
	RobloxID  int64      `json:"robloxId,omitempty"`
	Username  string     `json:"username,omitempty"`
	Reason    string     `json:"reason,omitempty"`
	BannedBy  string     `json:"bannedBy,omitempty"`
	BannedAt  time.Time  `json:"bannedAt,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	Active    bool       `json:"active,omitempty"`
	Notes     string     `json:"notes,omitempty"`
}

type Vehicle struct {
	ID          string    `json:"id,omitempty"`
	Model       string    `json:"model,omitempty"`
	Name        string    `json:"name,omitempty"`
	OwnerID     int64     `json:"ownerId,omitempty"`
	OwnerName   string    `json:"ownerName,omitempty"`
	Team        string    `json:"team,omitempty"`
	License     string    `json:"license,omitempty"`
	Location    string    `json:"location,omitempty"`
	Health      float64   `json:"health,omitempty"`
	Fuel        float64   `json:"fuel,omitempty"`
	SpawnedAt   time.Time `json:"spawnedAt,omitempty"`
	LastUpdated time.Time `json:"lastUpdated,omitempty"`
	IsStolen    bool      `json:"isStolen,omitempty"`
}

type Staff struct {
	ID          string    `json:"id,omitempty"`
	RobloxID    int64     `json:"robloxId,omitempty"`
	Username    string    `json:"username,omitempty"`
	DisplayName string    `json:"displayName,omitempty"`
	Role        string    `json:"role,omitempty"`
	Department  string    `json:"department,omitempty"`
	JoinedAt    time.Time `json:"joinedAt,omitempty"`
	IsActive    bool      `json:"isActive,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
}

type PlayerList struct {
	Total   int      `json:"total,omitempty"`
	Players []Player `json:"players,omitempty"`
	Count   int      `json:"count,omitempty"`
}

type QueueList struct {
	Total int            `json:"total,omitempty"`
	Queue []PlayerQueue  `json:"queue,omitempty"`
	Count int            `json:"count,omitempty"`
}

type JoinLogList struct {
	Total int       `json:"total,omitempty"`
	Logs  []JoinLog `json:"logs,omitempty"`
	Count int       `json:"count,omitempty"`
}

type CommandLogList struct {
	Total int          `json:"total,omitempty"`
	Logs  []CommandLog `json:"logs,omitempty"`
	Count int          `json:"count,omitempty"`
}

type KillLogList struct {
	Total int       `json:"total,omitempty"`
	Logs  []KillLog `json:"logs,omitempty"`
	Count int       `json:"count,omitempty"`
}

type ModCallList struct {
	Total int       `json:"total,omitempty"`
	Calls []ModCall `json:"calls,omitempty"`
	Count int       `json:"count,omitempty"`
}

type BanList struct {
	Total int   `json:"total,omitempty"`
	Bans  []Ban `json:"bans,omitempty"`
	Count int   `json:"count,omitempty"`
}

type VehicleList struct {
	Total    int       `json:"total,omitempty"`
	Vehicles []Vehicle `json:"vehicles,omitempty"`
	Count    int       `json:"count,omitempty"`
}

type StaffList struct {
	Total int     `json:"total,omitempty"`
	Staff []Staff `json:"staff,omitempty"`
	Count int     `json:"count,omitempty"`
}

type ListOptions struct {
	Limit  int    `json:"limit,omitempty"`
	Offset int    `json:"offset,omitempty"`
	Sort   string `json:"sort,omitempty"`
	Filter string `json:"filter,omitempty"`
}

type CacheOptions struct {
	Enabled bool
	TTL     time.Duration
}
