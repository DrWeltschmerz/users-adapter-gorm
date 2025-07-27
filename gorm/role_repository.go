package gorm

import (
	"context"
	"errors"

	core "github.com/DrWeltschmerz/users-core"
	"gorm.io/gorm"
)

type GormRoleRepository struct {
	db *gorm.DB
}

func NewGormRoleRepository(db *gorm.DB) *GormRoleRepository {
	return &GormRoleRepository{db: db}
}

func (r *GormRoleRepository) Create(ctx context.Context, role core.Role) (*core.Role, error) {
	gormRole := coreToGormRole(&role)
	if err := r.db.WithContext(ctx).Create(gormRole).Error; err != nil {
		return nil, err
	}
	return gormToCoreRole(gormRole), nil
}

func (r *GormRoleRepository) Update(ctx context.Context, role core.Role) (*core.Role, error) {
	gormRole := coreToGormRole(&role)
	if err := r.db.WithContext(ctx).Model(&GormRole{}).Where("id = ?", gormRole.ID).Updates(gormRole).Error; err != nil {
		return nil, err
	}
	var updated GormRole
	if err := r.db.WithContext(ctx).First(&updated, gormRole.ID).Error; err != nil {
		return nil, err
	}
	return gormToCoreRole(&updated), nil
}

func (r *GormRoleRepository) Delete(ctx context.Context, id string) error {
	uid, err := ParseStringUint(id)
	if err != nil {
		return errors.New("invalid role id")
	}
	res := r.db.WithContext(ctx).Delete(&GormRole{}, uid)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return core.ErrRoleNotFound
	}
	return nil
}

func (r *GormRoleRepository) GetByID(ctx context.Context, id string) (*core.Role, error) {
	var gormRole GormRole
	uid, err := ParseStringUint(id)
	if err != nil {
		return nil, errors.New("invalid role id")
	}
	if err := r.db.WithContext(ctx).First(&gormRole, uid).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRoleNotFound
		}
		return nil, err
	}
	return gormToCoreRole(&gormRole), nil
}

func (r *GormRoleRepository) GetByName(ctx context.Context, name string) (*core.Role, error) {
	var gormRole GormRole
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&gormRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRoleNotFound
		}
		return nil, err
	}
	return gormToCoreRole(&gormRole), nil
}

func (r *GormRoleRepository) List(ctx context.Context) ([]core.Role, error) {
	var gormRoles []GormRole
	if err := r.db.WithContext(ctx).Find(&gormRoles).Error; err != nil {
		return nil, err
	}
	roles := make([]core.Role, 0, len(gormRoles))
	for _, gr := range gormRoles {
		if r := gormToCoreRole(&gr); r != nil {
			roles = append(roles, *r)
		}
	}
	return roles, nil
}
