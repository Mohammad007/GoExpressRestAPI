# Database Configuration Guide

## SQLite and CGO Requirements

This project uses SQLite as its default database. The standard SQLite driver requires CGO to be enabled during compilation. If you encounter the following error:

```
failed to initialize database, got error Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
```

This means your Go binary was compiled without CGO support, which is required for the standard SQLite driver to work properly.

### Solutions

#### Option 1: Enable CGO (Recommended for Standard SQLite)

Recompile your application with CGO enabled:

```bash
# On Windows
set CGO_ENABLED=1
go build -o api.exe ./cmd/api

# On Linux/macOS
export CGO_ENABLED=1
go build -o api ./cmd/api
```

This requires a C compiler to be installed on your system:
- Windows: Install MinGW or TDM-GCC
- macOS: Install Xcode Command Line Tools
- Linux: Install GCC via your package manager

#### Option 2: Use Pure Go SQLite Implementation

The application now supports a pure Go SQLite implementation that doesn't require CGO. To use it, update your database configuration:

1. Edit `cmd/api/main.go` and change the database type to `sqlite-pure`:

```go
dbConfig := database.Config{
    Type:     "sqlite-pure", // Use pure Go SQLite implementation
    FilePath: "user_api.db",
}
```

This uses the `github.com/glebarez/sqlite` package, which is a pure Go implementation of SQLite that doesn't require CGO.

#### Option 3: Use an Alternative Database

If you prefer to use a different database type, you can configure the application accordingly:

1. Edit `cmd/api/main.go` and update the database configuration:

```go
dbConfig := database.Config{
    Type:     "mysql", // Change to "mysql", "postgres", or "mongodb"
    Host:     "localhost",
    Port:     "3306",
    User:     "root",
    Password: "your_password",
    DBName:   "user_api",
}
```

## Database Types

The application supports the following database types:

- `sqlite`: Standard SQLite database (requires CGO)
- `sqlite-pure`: Pure Go SQLite implementation (no CGO required)
- `mysql`: MySQL database
- `postgres`: PostgreSQL database
- `mongodb`: MongoDB database

Each database type requires different configuration parameters in the `database.Config` struct:

### SQLite and SQLite-Pure
```go
dbConfig := database.Config{
    Type:     "sqlite" or "sqlite-pure",
    FilePath: "path/to/database.db",
}
```

### MySQL and PostgreSQL
```go
dbConfig := database.Config{
    Type:     "mysql" or "postgres",
    Host:     "localhost",
    Port:     "3306" or "5432",
    User:     "username",
    Password: "password",
    DBName:   "database_name",
}
```

### MongoDB
```go
dbConfig := database.Config{
    Type:     "mongodb",
    Host:     "localhost",
    Port:     "27017",
    User:     "username",
    Password: "password",
    DBName:   "database_name",
}
```