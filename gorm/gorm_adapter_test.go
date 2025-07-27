package gorm

import (
	"context"
	"os"
	"testing"

	core "github.com/DrWeltschmerz/users-core"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	testDB   *gorm.DB
	userRepo core.UserRepository
	roleRepo core.RoleRepository
)

func TestMain(m *testing.M) {
	var err error
	testDB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = testDB.AutoMigrate(&GormUser{}, &GormRole{})
	if err != nil {
		panic(err)
	}

	roleRepo = NewGormRoleRepository(testDB)
	userRepo = NewGormUserRepository(testDB)

	code := m.Run()
	os.Exit(code)
}

func clearDB(t *testing.T) {
	require.NoError(t, testDB.Exec("DELETE FROM gorm_users").Error)
	require.NoError(t, testDB.Exec("DELETE FROM gorm_roles").Error)
}

// --- Helpers ---

func createTestRole(t *testing.T) *core.Role {
	role, err := roleRepo.Create(context.Background(), core.Role{Name: "test-role"})
	require.NoError(t, err)
	return role
}

func createTestUser(t *testing.T, roleID string) *core.User {
	user := core.User{
		Email:          "test@example.com",
		Username:       "testuser",
		HashedPassword: "hashedpw",
		RoleID:         roleID,
	}
	created, err := userRepo.Create(context.Background(), user)
	require.NoError(t, err)
	return created
}

// --- RoleRepository Tests ---

func TestRoleRepository_Create(t *testing.T) {
	clearDB(t)
	role := core.Role{Name: "admin"}
	created, err := roleRepo.Create(context.Background(), role)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, role.Name, created.Name)
}

func TestRoleRepository_GetByID(t *testing.T) {
	clearDB(t)
	created := createTestRole(t)

	got, err := roleRepo.GetByID(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, created.Name, got.Name)
}

func TestRoleRepository_GetByID_InvalidID(t *testing.T) {
	clearDB(t)
	_, err := roleRepo.GetByID(context.Background(), "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid role id")
}

func TestRoleRepository_GetByID_NotFound(t *testing.T) {
	clearDB(t)
	_, err := roleRepo.GetByID(context.Background(), "9999")
	require.ErrorIs(t, err, core.ErrRoleNotFound)
}

func TestRoleRepository_GetByName(t *testing.T) {
	clearDB(t)
	created := createTestRole(t)

	got, err := roleRepo.GetByName(context.Background(), created.Name)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)
}

func TestRoleRepository_GetByName_NotFound(t *testing.T) {
	clearDB(t)
	_, err := roleRepo.GetByName(context.Background(), "no-such-role")
	require.ErrorIs(t, err, core.ErrRoleNotFound)
}

func TestRoleRepository_Update(t *testing.T) {
	clearDB(t)
	created := createTestRole(t)

	created.Name = "updated-role"
	updated, err := roleRepo.Update(context.Background(), *created)
	require.NoError(t, err)
	require.Equal(t, "updated-role", updated.Name)
}

func TestRoleRepository_Delete(t *testing.T) {
	clearDB(t)
	created := createTestRole(t)

	err := roleRepo.Delete(context.Background(), created.ID)
	require.NoError(t, err)

	_, err = roleRepo.GetByID(context.Background(), created.ID)
	require.ErrorIs(t, err, core.ErrRoleNotFound)
}

func TestRoleRepository_Delete_InvalidID(t *testing.T) {
	clearDB(t)
	err := roleRepo.Delete(context.Background(), "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid role id")
}

func TestRoleRepository_Delete_NotFound(t *testing.T) {
	clearDB(t)
	err := roleRepo.Delete(context.Background(), "9999")
	require.ErrorIs(t, err, core.ErrRoleNotFound)
}

func TestRoleRepository_List(t *testing.T) {
	clearDB(t)
	r1 := createTestRole(t)
	r2, err := roleRepo.Create(context.Background(), core.Role{Name: "second-role"})
	require.NoError(t, err)

	roles, err := roleRepo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, roles, 2)

	// Check both roles are present
	var found1, found2 bool
	for _, r := range roles {
		if r.ID == r1.ID {
			found1 = true
		}
		if r.ID == r2.ID {
			found2 = true
		}
	}
	require.True(t, found1)
	require.True(t, found2)
}

// --- UserRepository additional negative tests ---

func TestUserRepository_GetByID_InvalidID(t *testing.T) {
	clearDB(t)
	_, err := userRepo.GetByID(context.Background(), "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid user id")
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	clearDB(t)
	_, err := userRepo.GetByID(context.Background(), "9999")
	require.ErrorIs(t, err, core.ErrUserNotFound)
}

func TestUserRepository_Delete_InvalidID(t *testing.T) {
	clearDB(t)
	err := userRepo.Delete(context.Background(), "invalid")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid user id")
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	clearDB(t)
	err := userRepo.Delete(context.Background(), "9999")
	require.ErrorIs(t, err, core.ErrUserNotFound)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	clearDB(t)
	_, err := userRepo.GetByEmail(context.Background(), "noone@example.com")
	require.ErrorIs(t, err, core.ErrUserNotFound)
}

func TestUserRepository_GetByUsername_NotFound(t *testing.T) {
	clearDB(t)
	_, err := userRepo.GetByUsername(context.Background(), "nousername")
	require.ErrorIs(t, err, core.ErrUserNotFound)
}
