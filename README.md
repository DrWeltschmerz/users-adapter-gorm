# users-adapter-gorm

GORM adapter for [users-core](https://github.com/DrWeltschmerz/users-core).

## Features

- Implements `UserRepository` and `RoleRepository` interfaces using GORM
- Conversion utilities between core and GORM models

## Usage

```go
import (
    "github.com/DrWeltschmerz/users-core"
    gormadapter "github.com/DrWeltschmerz/users-adapter-gorm/gorm"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
db.AutoMigrate(&gormadapter.GormUser{}, &gormadapter.GormRole{})

userRepo := gormadapter.NewGormUserRepository(db)
roleRepo := gormadapter.NewGormRoleRepository(db)
```

Use these repositories with the [`users-core` service](https://github.com/DrWeltschmerz/users-core).

## Testing

Run integration tests with:

```sh
go test ./gorm
```

## License

This project is licensed under the GNU GPL v3. See [LICENSE](LICENSE) for details.


## How to Use All Modules Together

1. **Set up your database** (e.g. SQLite, PostgreSQL).
2. **Auto-migrate** the GORM models (`GormUser`, `GormRole`).
3. **Create GORM repositories** from [users-adapter-gorm](https://github.com/DrWeltschmerz/users-adapter-gorm).
4. **Create a hasher** from [jwt-auth](https://github.com/DrWeltschmerz/jwt-auth).
5. **Create the service** from [users-core](https://github.com/DrWeltschmerz/users-core) using the above.
6. **Use the service** for user management in your app.


## TODO

- [ ] Create a Gin adapter to expose the service as a REST API (see [gin-gonic/gin](https://github.com/gin-gonic/gin)) (In Progress)
- [ ] Consider adding support for additional features:
    - Email verification and password reset flows
    - User profile endpoints
    - Role-based access control helpers
    - Activity/audit logging
    - OAuth2/social login integration
    - Rate limiting and brute-force protection
- [ ] Improve documentation and add more usage examples
- [ ] Expand integration test coverage