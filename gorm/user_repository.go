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
	gormUser, err := coreToGormUser(&user)
	if err != nil {
		return nil, err
	}
	updates := map[string]interface{}{
		"username":        gormUser.Username,
		"email":           gormUser.Email,
		"hashed_password": gormUser.HashedPassword,
		"role_id":         gormUser.RoleID,
		"last_seen":       gormUser.LastSeen,
	}
	if err := r.db.WithContext(ctx).Model(&GormUser{}).Where("id = ?", gormUser.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	var updated GormUser
	if err := r.db.WithContext(ctx).First(&updated, gormUser.ID).Error; err != nil {
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
