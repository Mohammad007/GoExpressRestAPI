package middleware

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/utils"
    "net/http"
)

func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        res := &framework.Response{w: w, status: http.StatusOK}
        defer func() {
            if err := recover(); err != nil {
                res.Error(http.StatusInternalServerError, "Internal server error")
            }
        }()
        next(w, r)
    }
}