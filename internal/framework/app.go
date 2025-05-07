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

type App struct {
    router      *mux.Router
    middlewares []func(http.HandlerFunc) http.HandlerFunc
    db          database.Database
    ctx         context.Context
}

type Response struct {
    w      http.ResponseWriter
    status int
}

func NewResponse(w http.ResponseWriter) *Response {
    return &Response{w: w, status: http.StatusOK}
}

type Router struct {
    app    *App
    prefix string
    mux    *mux.Router
}

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

func (app *App) Use(middleware func(http.HandlerFunc) http.HandlerFunc) {
    app.middlewares = append(app.middlewares, middleware)
}

func (app *App) Route(prefix string) *Router {
    subRouter := app.router.PathPrefix(prefix).Subrouter()
    return &Router{
        app:    app,
        prefix: prefix,
        mux:    subRouter,
    }
}

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

func (app *App) ParseBody(r *http.Request, v interface{}) error {
    return json.NewDecoder(r.Body).Decode(v)
}

func (res *Response) Respond(status int, data interface{}) {
    res.status = status
    res.w.Header().Set("Content-Type", "application/json")
    res.w.WriteHeader(status)
    if err := json.NewEncoder(res.w).Encode(data); err != nil {
        log.Printf("Error encoding JSON: %v", err)
    }
}

func (res *Response) Status(status int) *Response {
    res.status = status
    return res
}

func (res *Response) JSON(data interface{}) {
    res.Respond(res.status, data)
}

func (res *Response) Success(message string, data interface{}) {
    res.JSON(utils.SuccessResponse{Message: message, Data: data})
}

func (res *Response) Error(status int, message string) {
    res.Status(status).JSON(utils.ErrorResponse{Error: message})
}

func (app *App) DB() database.Database {
    return app.db
}

func (app *App) Context() context.Context {
    return app.ctx
}

func (app *App) Listen(port string) error {
    log.Printf("Server running on port %s", port)
    return http.ListenAndServe(port, app.router)
}