package gorm

import (
	"context"
	"errors"

	core "github.com/DrWeltschmerz/users-core"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user core.User) (*core.User, error) {
	gormUser, err := coreToGormUser(&user)
	if err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Create(gormUser).Error; err != nil {
		return nil, err
	}
	return gormToCoreUser(gormUser), nil
}

func (r *GormUserRepository) Update(ctx context.Context, user core.User) (*core.User, error) {
	// Fetch current user from DB
	var current GormUser
	gormID, err := ParseStringUint(user.ID)
	if err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).First(&current, gormID).Error; err != nil {
		return nil, err
	}

	// Only update fields that are non-zero in input
	updates := map[string]interface{}{}
	if user.Username != "" {
		updates["username"] = user.Username
	}
	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.HashedPassword != "" {
		updates["hashed_password"] = user.HashedPassword
	}
	if user.RoleID != "" {
		roleID, err := ParseStringUint(user.RoleID)
		if err == nil {
			updates["role_id"] = roleID
		}
	}
	if !user.LastSeen.IsZero() {
		updates["last_seen"] = user.LastSeen
	}

	if len(updates) == 0 {
		// Nothing to update
		return gormToCoreUser(&current), nil
	}

	if err := r.db.WithContext(ctx).Model(&GormUser{}).Where("id = ?", gormID).Updates(updates).Error; err != nil {
		return nil, err
	}
	var updated GormUser
	if err := r.db.WithContext(ctx).First(&updated, gormID).Error; err != nil {
		return nil, err
	}
	return gormToCoreUser(&updated), nil
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*core.User, error) {
	var gormUser GormUser
	uid, err := ParseStringUint(id)
	if err != nil {
		return nil, errors.New("invalid user id")
	}
	if err := r.db.WithContext(ctx).First(&gormUser, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}
		// ...existing code...
	return gormToCoreUser(&gormUser), nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*core.User, error) {
	var gormUser GormUser
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}
	return gormToCoreUser(&gormUser), nil
}

func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*core.User, error) {
	var gormUser GormUser
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrUserNotFound
		}
		return nil, err
	}
	return gormToCoreUser(&gormUser), nil
}

func (r *GormUserRepository) List(ctx context.Context) ([]core.User, error) {
	var gormUsers []GormUser
	if err := r.db.WithContext(ctx).Find(&gormUsers).Error; err != nil {
		return nil, err
	}
	users := make([]core.User, 0, len(gormUsers))
	for _, gu := range gormUsers {
		if u := gormToCoreUser(&gu); u != nil {
			users = append(users, *u)
		}
	}
	return users, nil
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	uid, err := ParseStringUint(id)
	if err != nil {
		return errors.New("invalid user id")
	}
	res := r.db.WithContext(ctx).Delete(&GormUser{}, uid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return core.ErrUserNotFound
	}
	return nil
}
