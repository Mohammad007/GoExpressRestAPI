package main

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/middleware"
    "github.com/Mohammad007/GoExpressRestAPI/internal/routes"
    "log"
    "strings"
)

func main() {
    // Database configuration
    dbConfig := database.Config{
        // Default to sqlite-pure for better compatibility without CGO
        Type:     "sqlite-pure", // Options: "sqlite", "sqlite-pure", "mysql", "postgres", "mongodb"
        FilePath: "user_api.db",
        // MySQL/PostgreSQL configuration (uncomment and configure if needed)
        // Host:     "localhost",
        // Port:     "3306", // MySQL: 3306, PostgreSQL: 5432, MongoDB: 27017
        // User:     "root",
        // Password: "password",
        // DBName:   "user_api",
    }
    
    app, err := framework.NewApp(dbConfig)
    if err != nil {
        // Provide more helpful error message for common issues
        if strings.Contains(err.Error(), "CGO_ENABLED=0") || 
           strings.Contains(err.Error(), "requires cgo") || 
           strings.Contains(err.Error(), "SQLite requires CGO") {
            log.Fatalf("Database error: SQLite requires CGO_ENABLED=1 to work.\n" +
                "Please either:\n" +
                "1. Rebuild with CGO_ENABLED=1\n" +
                "2. Use 'sqlite-pure' type in main.go (pure Go SQLite implementation)\n" +
                "3. Configure an alternative database in main.go (MySQL/PostgreSQL/MongoDB)")
        }
        log.Fatal("Failed to initialize app:", err)
    }
    defer app.DB().Close()

    app.Use(middleware.ErrorHandler)
    app.Use(middleware.Logger)

    routes.RegisterUserRoutes(app)

    if err := app.Listen(":8080"); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}