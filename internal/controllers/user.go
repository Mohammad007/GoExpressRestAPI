// internal/controllers/user.go
package controllers

import (
    "github.com/go-playground/validator/v10"
    "github.com/gorilla/mux"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
    "net/http"
    "strconv"
)

func CreateUser(app *framework.App) func(r *http.Request, res *framework.Response) {
    return func(r *http.Request, res *framework.Response) {
        var user models.User
        if err := app.ParseBody(r, &user); err != nil {
            res.Error(http.StatusBadRequest, "Invalid request payload")
            return
        }
        if err := user.Validate(); err != nil {
            res.Error(http.StatusBadRequest, "Validation failed: "+err.Error())
            return
        }
        if err := app.DB().CreateUser(app.Context(), &user); err != nil {
            res.Error(http.StatusInternalServerError, "Failed to create user")
            return
        }
        res.Status(http.StatusCreated).Success("User created successfully", user)
    }
}

func GetAllUsers(app *framework.App) func(r *http.Request, res *framework.Response) {
    return func(r *http.Request, res *framework.Response) {
        users, err := app.DB().GetAllUsers(app.Context())
        if err != nil {
            res.Error(http.StatusInternalServerError, "Failed to fetch users")
            return
        }
        res.Success("Users fetched successfully", users)
    }
}

func GetUserByID(app *framework.App) func(r *http.Request, res *framework.Response) {
    return func(r *http.Request, res *framework.Response) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            res.Error(http.StatusBadRequest, "Invalid user ID")
            return
        }
        user, err := app.DB().GetUserByID(app.Context(), uint(id))
        if err != nil {
            res.Error(http.StatusNotFound, "User not found")
            return
        }
        res.Success("User fetched successfully", user)
    }
}

func UpdateUser(app *framework.App) func(r *http.Request, res *framework.Response) {
    return func(r *http.Request, res *framework.Response) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            res.Error(http.StatusBadRequest, "Invalid user ID")
            return
        }
        var user models.User
        if err := app.ParseBody(r, &user); err != nil {
            res.Error(http.StatusBadRequest, "Invalid request payload")
            return
        }
        if err := user.Validate(); err != nil {
            res.Error(http.StatusBadRequest, "Validation failed: "+err.Error())
            return
        }
        user.ID = uint(id)
        if err := app.DB().UpdateUser(app.Context(), &user); err != nil {
            res.Error(http.StatusInternalServerError, "Failed to update user")
            return
        }
        res.Success("User updated successfully", user)
    }
}

func DeleteUser(app *framework.App) func(r *http.Request, res *framework.Response) {
    return func(r *http.Request, res *framework.Response) {
        vars := mux.Vars(r)
        id, err := strconv.Atoi(vars["id"])
        if err != nil {
            res.Error(http.StatusBadRequest, "Invalid user ID")
            return
        }
        if err := app.DB().DeleteUser(app.Context(), uint(id)); err != nil {
            res.Error(http.StatusInternalServerError, "Failed to delete user")
            return
        }
        res.Success("User deleted successfully", nil)
    }
}