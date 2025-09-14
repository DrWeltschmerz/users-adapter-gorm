# users-adapter-gorm

GORM adapter for [users-core](https://github.com/DrWeltschmerz/users-core).

## Features

- Implements `UserRepository` and `RoleRepository` interfaces using GORM
- Conversion utilities between core and GORM models

## Usage

This module provides GORM-based implementations for the repository interfaces defined in [users-core](https://github.com/DrWeltschmerz/users-core).

```go
import (
    gormadapter "github.com/DrWeltschmerz/users-adapter-gorm/gorm"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
db.AutoMigrate(&gormadapter.GormUser{}, &gormadapter.GormRole{})

userRepo := gormadapter.NewGormUserRepository(db)
roleRepo := gormadapter.NewGormRoleRepository(db)
```

See [users-core](https://github.com/DrWeltschmerz/users-core) for how to wire these repositories into a service.

## Testing

Run integration tests with:

```sh
go test ./gorm
```

## License

This project is licensed under the GNU GPL v3. See [LICENSE](LICENSE) for details.
    - Rate limiting and brute-force protection
- [ ] Improve documentation and add more usage examples
- [ ] Expand integration test coverage