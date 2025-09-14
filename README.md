# users-adapter-gorm

Gorm adapter for modular user system

## Overview

This package provides Gorm-based implementations of user and role repositories for the [`users-core`](https://github.com/DrWeltschmerz/users-core) modular user system.

## Features

- Gorm models for users and roles
- Repository implementations for CRUD operations on users and roles
- Conversion utilities between core and Gorm models
- Comprehensive unit tests using SQLite in-memory DB

## Project Structure

```
.
├── gorm/
│   ├── gorm_adapter_test.go   # Unit tests for repositories
│   ├── role.go                # GormRole model
│   ├── role_repository.go     # GormRoleRepository implementation
│   ├── user.go                # GormUser model
│   ├── user_repository.go     # GormUserRepository implementation
│   └── utils.go               # Conversion and parsing utilities
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Usage

1. **Install dependencies**  
   Make sure you have Go 1.18+ and run:
   ```sh
   go mod tidy
   ```

2. **Initialize Gorm and repositories**
   ```go
   import (
       "gorm.io/driver/sqlite"
       "gorm.io/gorm"
       adapter "github.com/DrWeltschmerz/users-adapter-gorm/gorm"
   )

   db, err := gorm.Open(sqlite.Open("your.db"), &gorm.Config{})
   if err != nil {
       // handle error
   }

   // Auto-migrate models
   db.AutoMigrate(&adapter.GormUser{}, &adapter.GormRole{})

   userRepo := adapter.NewGormUserRepository(db)
   roleRepo := adapter.NewGormRoleRepository(db)
   ```

3. **Implement the `users-core` interfaces**  
   The repositories implement the `core.UserRepository` and `core.RoleRepository` interfaces from [`users-core`](https://github.com/DrWeltschmerz/users-core).

## Testing

Run all tests with:
```sh
go test ./gorm/...
```

## License

This project is licensed under the GNU GPL v3. See [LICENSE](LICENSE) for details.
