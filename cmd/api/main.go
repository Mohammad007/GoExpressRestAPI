// cmd/api/main.go
package main

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/middleware"
    "github.com/Mohammad007/GoExpressRestAPI/internal/routes"
    "log"
)

func main() {
    // Database configuration
    dbConfig := database.Config{
        Type:     "sqlite",
        FilePath: "user_api.db",
    }

    // Create framework instance
    app, err := framework.NewApp(dbConfig)
    if err != nil {
        log.Fatal("Failed to initialize app:", err)
    }
    defer app.DB().Close()

    // Add middlewares
    app.Use(middleware.ErrorHandler)
    app.Use(middleware.Logger)

    // Register routes
    routes.RegisterUserRoutes(app)

    // Start server
    if err := app.Listen(":8080"); err != nil {
        log.Fatal("Server failed to start:", err)
    }
}