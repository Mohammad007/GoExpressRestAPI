// internal/framework/app.go
package framework

import (
    "context"
    "encoding/json"
    "github.com/gorilla/mux"
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/utils"
    "log"
    "net/http"
)

// App is the core of the framework
type App struct {
    router      *mux.Router
    middlewares []func(http.HandlerFunc) http.HandlerFunc
    db          database.Database
    ctx         context.Context
}

// Response wraps http.ResponseWriter for chainable methods
type Response struct {
    w      http.ResponseWriter
    status int
}

// Router defines a route group, similar to Express.js router
type Router struct {
    app    *App
    prefix string
    mux    *mux.Router
}

// NewApp creates a new framework instance
func NewApp(dbConfig database.Config) (*App, error) {
    db, err := database.NewDatabase(dbConfig)
    if err != nil {
        return nil, err
    }
    if err := db.Connect(); err != nil {
        return nil, err
    }

    return &App{
        router:      mux.NewRouter(),
        middlewares: []func(http.HandlerFunc) http.HandlerFunc{},
        db:          db,
        ctx:         context.Background(),
    }, nil
}

// Use adds a global middleware
func (app *App) Use(middleware func(http.HandlerFunc) http.HandlerFunc) {
    app.middlewares = append(app.middlewares, middleware)
}

// Route creates a new router for a path prefix
func (app *App) Route(prefix string) *Router {
    subRouter := app.router.PathPrefix(prefix).Subrouter()
    return &Router{
        app:    app,
        prefix: prefix,
        mux:    subRouter,
    }
}

// GET, POST, PUT, DELETE for Router
func (r *Router) GET(path string, handler func(r *http.Request, res *Response)) *Router {
    r.registerRoute(path, http.MethodGet, handler)
    return r
}

func (r *Router) POST(path string, handler func(r *http.Request, res *Response)) *Router {
    r.registerRoute(path, http.MethodPost, handler)
    return r
}

func (r *Router) PUT(path string, handler func(r *http.Request, res *Response)) *Router {
    r.registerRoute(path, http.MethodPut, handler)
    return r
}

func (r *Router) DELETE(path string, handler func(r *http.Request, res *Response)) *Router {
    r.registerRoute(path, http.MethodDelete, handler)
    return r
}

func (r *Router) registerRoute(path, method string, handler func(r *http.Request, res *Response)) {
    wrappedHandler := func(w http.ResponseWriter, req *http.Request) {
        res := &Response{w: w, status: http.StatusOK}
        handler(req, res)
    }
    for _, middleware := range r.app.middlewares {
        wrappedHandler = middleware(wrappedHandler)
    }
    r.mux.HandleFunc(path, wrappedHandler).Methods(method)
}

// ParseBody parses JSON body
func (app *App) ParseBody(r *http.Request, v interface{}) error {
    return json.NewDecoder(r.Body).Decode(v)
}

// Respond sends a JSON response
func (res *Response) Respond(status int, data interface{}) {
    res.status = status
    res.w.Header().Set("Content-Type", "application/json")
    res.w.WriteHeader(status)
    if err := json.NewEncoder(res.w).Encode(data); err != nil {
        log.Printf("Error encoding JSON: %v", err)
    }
}

// Status sets the response status code
func (res *Response) Status(status int) *Response {
    res.status = status
    return res
}

// JSON sends a JSON response
func (res *Response) JSON(data interface{}) {
    res.Respond(res.status, data)
}

// Success sends a success response
func (res *Response) Success(message string, data interface{}) {
    res.JSON(utils.SuccessResponse{Message: message, Data: data})
}

// Error sends an error response
func (res *Response) Error(status int, message string) {
    res.Status(status).JSON(utils.ErrorResponse{Error: message})
}

// DB returns the database instance
func (app *App) DB() database.Database {
    return app.db
}

// Context returns the context
func (app *App) Context() context.Context {
    return app.ctx
}

// Listen starts the server
func (app *App) Listen(port string) error {
    log.Printf("Server running on port %s", port)
    return http.ListenAndServe(port, app.router)
}