package gorm

import (
	"fmt"

	core "github.com/DrWeltschmerz/users-core"
)

// ParseStringUint parses a string ID to uint, returns 0 and error if invalid.
func ParseStringUint(s string) (uint, error) {
	var id uint
	_, err := fmt.Sscanf(s, "%d", &id)
	if err != nil {
		return 0, fmt.Errorf("invalid uint string: %w", err)
	}
	return id, nil
}

func coreToGormUser(user *core.User) (*GormUser, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}
	var id, roleID uint
	var err error
	if user.ID != "" {
		id, err = ParseStringUint(user.ID)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %w", err)
		}
	}
	if user.RoleID != "" {
		roleID, err = ParseStringUint(user.RoleID)
		if err != nil {
			return nil, fmt.Errorf("invalid role ID: %w", err)
		}
	}
	return &GormUser{
		ID:             id,
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		RoleID:         roleID,
		LastSeen:       user.LastSeen,
	}, nil
}

func gormToCoreUser(gormUser *GormUser) *core.User {
	if gormUser == nil {
		return nil
	}
	return &core.User{
		ID:             fmt.Sprintf("%d", gormUser.ID),
		Email:          gormUser.Email,
		HashedPassword: gormUser.HashedPassword,
		Username:       gormUser.Username,
		LastSeen:       gormUser.LastSeen,
		RoleID:         fmt.Sprintf("%d", gormUser.RoleID),
	}
}

func gormToCoreRole(gormRole *GormRole) *core.Role {
	if gormRole == nil {
		return nil
	}

	return &core.Role{
		ID:   fmt.Sprintf("%d", gormRole.ID),
		Name: gormRole.Name,
	}
}

func coreToGormRole(role *core.Role) *GormRole {
	if role == nil {
		return nil
	}

	var id uint
	if role.ID != "" {
		parsedID, err := ParseStringUint(role.ID)
		if err != nil {
			id = 0
		} else {
			id = parsedID
		}
	}

	return &GormRole{
		ID:   id,
		Name: role.Name,
	}
}
