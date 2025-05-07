// internal/routes/user.go
package routes

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/controllers"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
)

// RegisterUserRoutes registers user-related routes
func RegisterUserRoutes(app *framework.App) {
    router := app.Route("/users")
    router.
        GET("/", controllers.GetAllUsers(app)).
        GET("/{id}", controllers.GetUserByID(app)).
        POST("/", controllers.CreateUser(app)).
        PUT("/{id}", controllers.UpdateUser(app)).
        DELETE("/{id}", controllers.DeleteUser(app))
}