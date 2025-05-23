package utils

type SuccessResponse struct {
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}