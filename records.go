package erlc

import (
	"context"
	"net/http"
	"time"
)

func (c *Client) GetBans(ctx context.Context) (*BanList, error) {
	var result BanList
	err := c.do(ctx, http.MethodGet, EndpointBans, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetActiveBans(ctx context.Context) ([]Ban, error) {
	bans, err := c.GetBans(ctx)
	if err != nil {
		return nil, err
	}

	var active []Ban
	for _, ban := range bans.Bans {
		if ban.Active {
			active = append(active, ban)
		}
	}
	return active, nil
}

func (c *Client) GetPermanentBans(ctx context.Context) ([]Ban, error) {
	bans, err := c.GetActiveBans(ctx)
	if err != nil {
		return nil, err
	}

	var permanent []Ban
	for _, ban := range bans {
		if ban.ExpiresAt == nil {
			permanent = append(permanent, ban)
		}
	}
	return permanent, nil
}

func (c *Client) GetTemporaryBans(ctx context.Context) ([]Ban, error) {
	bans, err := c.GetActiveBans(ctx)
	if err != nil {
		return nil, err
	}

	var temporary []Ban
	for _, ban := range bans {
		if ban.ExpiresAt != nil && ban.ExpiresAt.After(time.Now()) {
			temporary = append(temporary, ban)
		}
	}
	return temporary, nil
}

func (c *Client) GetExpiredBans(ctx context.Context) ([]Ban, error) {
	bans, err := c.GetBans(ctx)
	if err != nil {
		return nil, err
	}

	var expired []Ban
	for _, ban := range bans.Bans {
		if !ban.Active || (ban.ExpiresAt != nil && ban.ExpiresAt.Before(time.Now())) {
			expired = append(expired, ban)
		}
	}
	return expired, nil
}

func (c *Client) FindBanByUsername(ctx context.Context, username string) (*Ban, error) {
	bans, err := c.GetBans(ctx)
	if err != nil {
		return nil, err
	}

	for _, ban := range bans.Bans {
		if ban.Username == username {
			return &ban, nil
		}
	}
	return nil, nil
}

func (c *Client) FindBanByRobloxID(ctx context.Context, robloxID int64) (*Ban, error) {
	bans, err := c.GetBans(ctx)
	if err != nil {
		return nil, err
	}

	for _, ban := range bans.Bans {
		if ban.RobloxID == robloxID {
			return &ban, nil
		}
	}
	return nil, nil
}

func (c *Client) GetVehicles(ctx context.Context) (*VehicleList, error) {
	var result VehicleList
	err := c.doWithCache(ctx, EndpointVehicles, &result, nil)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetVehiclesWithCache(ctx context.Context, cacheOpt *CacheOptions) (*VehicleList, error) {
	var result VehicleList
	err := c.doWithCache(ctx, EndpointVehicles, &result, cacheOpt)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) FilterVehiclesByTeam(ctx context.Context, team string) ([]Vehicle, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Vehicle
	for _, v := range vehicles.Vehicles {
		if v.Team == team {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}

func (c *Client) FilterVehiclesByModel(ctx context.Context, model string) ([]Vehicle, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Vehicle
	for _, v := range vehicles.Vehicles {
		if v.Model == model {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}

func (c *Client) FilterVehiclesByOwner(ctx context.Context, ownerName string) ([]Vehicle, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Vehicle
	for _, v := range vehicles.Vehicles {
		if v.OwnerName == ownerName {
			filtered = append(filtered, v)
		}
	}
	return filtered, nil
}

func (c *Client) FindVehicleByLicense(ctx context.Context, license string) (*Vehicle, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range vehicles.Vehicles {
		if v.License == license {
			return &v, nil
		}
	}
	return nil, nil
}

func (c *Client) GetStolenVehicles(ctx context.Context) ([]Vehicle, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return nil, err
	}

	var stolen []Vehicle
	for _, v := range vehicles.Vehicles {
		if v.IsStolen {
			stolen = append(stolen, v)
		}
	}
	return stolen, nil
}

func (c *Client) GetVehiclesByHealthStatus(ctx context.Context, healthThreshold float64) ([]Vehicle, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return nil, err
	}

	var damaged []Vehicle
	for _, v := range vehicles.Vehicles {
		if v.Health < healthThreshold {
			damaged = append(damaged, v)
		}
	}
	return damaged, nil
}

func (c *Client) GetVehicleCount(ctx context.Context) (int, error) {
	vehicles, err := c.GetVehicles(ctx)
	if err != nil {
		return 0, err
	}
	return vehicles.Count, nil
}

func (c *Client) GetStaff(ctx context.Context) (*StaffList, error) {
	var result StaffList
	err := c.do(ctx, http.MethodGet, EndpointStaff, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) FilterStaffByRole(ctx context.Context, role string) ([]Staff, error) {
	staff, err := c.GetStaff(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Staff
	for _, s := range staff.Staff {
		if s.Role == role {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}

func (c *Client) FilterStaffByDepartment(ctx context.Context, department string) ([]Staff, error) {
	staff, err := c.GetStaff(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []Staff
	for _, s := range staff.Staff {
		if s.Department == department {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}

func (c *Client) FindStaffByUsername(ctx context.Context, username string) (*Staff, error) {
	staff, err := c.GetStaff(ctx)
	if err != nil {
		return nil, err
	}

	for _, s := range staff.Staff {
		if s.Username == username {
			return &s, nil
		}
	}
	return nil, nil
}

func (c *Client) GetActiveStaff(ctx context.Context) ([]Staff, error) {
	staff, err := c.GetStaff(ctx)
	if err != nil {
		return nil, err
	}

	var active []Staff
	for _, s := range staff.Staff {
		if s.IsActive {
			active = append(active, s)
		}
	}
	return active, nil
}
