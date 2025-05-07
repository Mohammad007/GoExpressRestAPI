# GoExpressRestAPI

A lightweight, Express.js-inspired REST API framework built in Go, designed to feel familiar to Node.js developers. It combines the simplicity of Express.js with the performance and type safety of Go, supporting multiple databases and modern API development practices.

> **Important Note**: This project uses SQLite as its default database, which requires CGO to be enabled during compilation. If you encounter errors about CGO being disabled, please refer to the [Database Configuration Guide](docs/database_guide.md) for solutions.

## Table of Contents
1. [Introduction](#introduction)
2. [Features](#features)
3. [Installation](#installation)
4. [Project Setup](#project-setup)
5. [Core Concepts](#core-concepts)
   - [App](#app)
   - [Router](#router)
   - [Models](#models)
   - [Controllers](#controllers)
   - [Routes](#routes)
   - [Middleware](#middleware)
   - [Database](#database)
6. [Building a CRUD API](#building-a-crud-api)
   - [Step 1: Define a Model](#step-1-define-a-model)
   - [Step 2: Create Controllers](#step-2-create-controllers)
   - [Step 3: Define Routes](#step-3-define-routes)
   - [Step 4: Setup Main File](#step-4-setup-main-file)
   - [Step 5: Test the API](#step-5-test-the-api)
7. [Advanced Features](#advanced-features)
   - [Custom Middleware](#custom-middleware)
   - [Environment Variables](#environment-variables)
   - [Route-Specific Middleware](#route-specific-middleware)
   - [Custom Validation](#custom-validation)
8. [Database Configuration](#database-configuration)
   - [MySQL](#mysql)
   - [PostgreSQL](#postgresql)
   - [SQLite](#sqlite)
   - [MongoDB](#mongodb)
9. [Testing](#testing)
10. [Examples](#examples)
    - [Simple Hello World API](#simple-hello-world-api)
    - [User Management CRUD API](#user-management-crud-api)
11. [Contributing](#contributing)
12. [License](#license)

## Introduction

**GoExpressRestAPI** is a REST API framework for Go that mirrors the simplicity and developer experience of Express.js. It’s built for developers who love the Node.js ecosystem but want the performance, type safety, and concurrency of Go. With features like chainable routing, Mongoose-like models, and support for multiple databases (MySQL, PostgreSQL, SQLite, MongoDB), it’s perfect for building modern APIs.

This documentation guides you through installing, configuring, and using the framework to build production-ready APIs. Whether you’re a Node.js developer transitioning to Go or a Go developer looking for an Express.js-like experience, this framework is designed for you.

## Features

- **Express.js-Like Routing**: Modular routing with a dedicated `routes/` folder, supporting chainable methods (`router.GET().POST()`).
- **Mongoose-Like Models**: Schema-based model definitions with built-in validation and hooks, inspired by Mongoose.
- **Chainable Responses**: Intuitive response handling with methods like `res.Status(200).JSON(...)`.
- **Middleware Support**: Global and route-specific middleware for logging, authentication, error handling, and more.
- **Database Support**: Seamless integration with MySQL, PostgreSQL, SQLite, and MongoDB via a unified `Database` interface.
- **Type Safety**: Leverages Go’s static typing for robust and maintainable code.
- **Developer-Friendly**: Familiar folder structure (`controllers/`, `models/`, `routes/`) and Express.js-inspired syntax.

## Installation

To use **GoExpressRestAPI**, you need Go (version 1.18 or higher) installed. Run the following command to add the framework to your project:

```bash
go get github.com/Mohammad007/GoExpressRestAPI
```

You’ll also need dependencies for database support and validation:

```bash
go get github.com/go-playground/validator/v10
go get github.com/gorilla/mux
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/driver/postgres
go get gorm.io/driver/sqlite
go get go.mongodb.org/mongo-driver/mongo
```

For environment variable support (optional):

```bash
go get github.com/joho/godotenv
```

## Project Setup

To start a new project with **GoExpressRestAPI**:

1. **Create a New Project**:
   ```bash
   mkdir my-api
   cd my-api
   go mod init my-api
   go get github.com/Mohammad007/GoExpressRestAPI
   ```

2. **Recommended Folder Structure**:
   ```
   my-api/
   ├── main.go              # Entry point
   ├── controllers/         # Business logic for API endpoints
   │   └── user.go
   ├── models/              # Data models with schema and validation
   │   └── user.go
   ├── routes/              # Route definitions
   │   └── user.go
   ├── go.mod               # Go module file
   ├── go.sum               # Dependency checksums
   └── .env                 # Environment variables (optional)
   ```

3. **Basic Example**:
   Create a simple API in `main.go`:

   ```go
   package main

   import (
       "github.com/Mohammad007/GoExpressRestAPI/internal/database"
       "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
       "github.com/Mohammad007/GoExpressRestAPI/internal/middleware"
   )

   func main() {
       dbConfig := database.Config{Type: "sqlite", FilePath: "test.db"}
       app, err := framework.NewApp(dbConfig)
       if err != nil {
           panic(err)
       }
       defer app.DB().Close()

       app.Use(middleware.Logger)

       app.Route("/hello").GET(func(r *http.Request, res *framework.Response) {
           res.Success("Hello, World!", nil)
       })

       app.Listen(":8080")
   }
   ```

4. **Run the Project**:
   ```bash
   go run main.go
   ```

   Test the endpoint:
   ```bash
   curl http://localhost:8080/hello
   # Output: {"message":"Hello, World!"}
   ```

## Core Concepts

### App
The `App` struct is the core of the framework, managing routing, middleware, and database connections.

- **Key Methods**:
  - `NewApp(dbConfig)`: Initializes the app with a database configuration.
  - `Use(middleware)`: Adds global middleware.
  - `Route(prefix)`: Creates a route group for a path prefix.
  - `Listen(port)`: Starts the server.
  - `ParseBody(r, v)`: Parses JSON request body.
  - `DB()`: Returns the database instance.
  - `Context()`: Returns the request context.

### Router
The `Router` struct allows modular route definitions, similar to `express.Router()`.

- **Key Methods**:
  - `GET(path, handler)`: Defines a GET route.
  - `POST(path, handler)`: Defines a POST route.
  - `PUT(path, handler)`: Defines a PUT route.
  - `DELETE(path, handler)`: Defines a DELETE route.
  - Chainable: `router.GET().POST().PUT()`

### Models
Models are defined using a Mongoose-like schema structure with validation and hooks.

- **Structure**:
  - A `Schema` struct defines fields, types, and validation rules.
  - A factory function (`NewModel`) creates instances.
  - Validation is handled via `validator` tags and a `Validate()` method.
  - GORM hooks (`BeforeCreate`) provide pre-save logic.

### Controllers
Controllers contain the business logic for API endpoints, similar to Express.js controllers.

- **Structure**:
  - Functions with signature `func(r *http.Request, res *framework.Response)`.
  - Use `app.ParseBody` for request parsing and `res.Success`/`res.Error` for responses.
  - Validation is typically handled at the model level.

### Routes
Routes are defined in the `routes/` folder, keeping routing logic separate from controllers.

- **Structure**:
  - A `RegisterXRoutes(app)` function defines routes for a resource.
  - Uses `app.Route(prefix)` to create a route group.
  - Chainable methods for defining HTTP methods.

### Middleware
Middleware functions process requests before they reach controllers.

- **Types**:
  - Global middleware: Applied via `app.Use`.
  - Route-specific middleware: Applied via `router.Use`.
- **Built-in Middleware**:
  - `middleware.Logger`: Logs requests.
  - `middleware.ErrorHandler`: Handles errors.

### Database
The framework supports multiple databases through a unified `Database` interface.

- **Supported Databases**: MySQL, PostgreSQL, SQLite, MongoDB.
- **Configuration**: Defined via `database.Config`.
- **Methods**:
  - `Connect()`: Establishes a connection.
  - `Close()`: Closes the connection.
  - CRUD operations: `CreateUser`, `GetUserByID`, etc.

## Building a CRUD API

Let’s build a **User Management API** with CRUD operations (`POST /users`, `GET /users`, `GET /users/:id`, `PUT /users/:id`, `DELETE /users/:id`).

### Step 1: Define a Model
Create a Mongoose-like model in `models/user.go`:

```go
package models

import (
    "github.com/go-playground/validator/v10"
    "gorm.io/gorm"
    "time"
)

type UserSchema struct {
    ID        uint       `json:"id" gorm:"primaryKey" bson:"id"`
    Name      string     `json:"name" validate:"required,min=2" gorm:"type:varchar(100)" bson:"name"`
    Email     string     `json:"email" validate:"required,email" gorm:"unique;type:varchar(100)" bson:"email"`
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime" bson:"created_at"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime" bson:"updated_at"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index" bson:"deleted_at"`
}

type User struct {
    UserSchema
    validator *validator.Validate
}

func NewUser(name, email string) *User {
    user := &User{
        UserSchema: UserSchema{
            Name:      name,
            Email:     email,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        },
        validator: validator.New(),
    }
    return user
}

func (u *User) Validate() error {
    return u.validator.Struct(u)
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    return u.Validate()
}
```

### Step 2: Create Controllers
Define controllers in `controllers/user.go`:

```go
package controllers

import (
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
```

### Step 3: Define Routes
Define routes in `routes/user.go`:

```go
package routes

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/controllers"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
)

func RegisterUserRoutes(app *framework.App) {
    router := app.Route("/users")
    router.
        GET("/", controllers.GetAllUsers(app)).
        GET("/{id}", controllers.GetUserByID(app)).
        POST("/", controllers.CreateUser(app)).
        PUT("/{id}", controllers.UpdateUser(app)).
        DELETE("/{id}", controllers.DeleteUser(app))
}
```

### Step 4: Setup Main File
Configure the app in `main.go`:

```go
package main

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/middleware"
    "github.com/Mohammad007/GoExpressRestAPI/internal/routes"
    "log"
)

func main() {
    dbConfig := database.Config{
        Type:     "sqlite",
        FilePath: "user_api.db",
    }
    app, err := framework.NewApp(dbConfig)
    if err != nil {
        log.Fatal("Failed to initialize app:", err)
    }
    defer app.DB().Close()

    app.Use(middleware.ErrorHandler)
    app.Use(middleware.Logger)

    routes.RegisterUserRoutes(app)

    app.Listen(":8080")
}
```

### Step 5: Test the API
Run the server:

```bash
go run main.go
```

Test the endpoints:

- **Create User**:
  ```bash
  curl -X POST -H "Content-Type: application/json" -d '{"name":"Jane Doe","email":"jane@example.com"}' http://localhost:8080/users
  # Output: {"message":"User created successfully","data":{"id":1,"name":"Jane Doe","email":"jane@example.com","created_at":"2025-05-08T09:30:00Z","updated_at":"2025-05-08T09:30:00Z"}}
  ```

- **Get All Users**:
  ```bash
  curl http://localhost:8080/users
  ```

- **Get User by ID**:
  ```bash
  curl http://localhost:8080/users/1
  ```

- **Update User**:
  ```bash
  curl -X PUT -H "Content-Type: application/json" -d '{"name":"Jane Smith","email":"jane.smith@example.com"}' http://localhost:8080/users/1
  ```

- **Delete User**:
  ```bash
  curl -X DELETE http://localhost:8080/users/1
  ```

## Advanced Features

### Custom Middleware
Create custom middleware in `middleware/auth.go`:

```go
package middleware

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "net/http"
)

func AuthMiddleware() func(http.HandlerFunc) http.HandlerFunc {
    return func(next http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            token := r.Header.Get("Authorization")
            if token != "Bearer my-secret-token" {
                res := &framework.Response{w: w}
                res.Error(http.StatusUnauthorized, "Invalid token")
                return
            }
            next(w, r)
        }
    }
}
```

Apply it globally:

```go
app.Use(middleware.AuthMiddleware())
```

Or to specific routes:

```go
router.Use(middleware.AuthMiddleware())
```

### Environment Variables
Use `godotenv` to load configuration from a `.env` file:

```bash
go get github.com/joho/godotenv
```

Create a `.env` file:

```
DB_TYPE=sqlite
DB_FILE=user_api.db
```

Update `main.go`:

```go
package main

import (
    "github.com/joho/godotenv"
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/middleware"
    "github.com/Mohammad007/GoExpressRestAPI/internal/routes"
    "log"
    "os"
)

func main() {
    godotenv.Load()
    dbConfig := database.Config{
        Type:     os.Getenv("DB_TYPE"),
        Host:     os.Getenv("DB_HOST"),
        Port:     os.Getenv("DB_PORT"),
        User:     os.Getenv("DB_USER"),
        Password: os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        FilePath: os.Getenv("DB_FILE"),
    }
    app, err := framework.NewApp(dbConfig)
    if err != nil {
        log.Fatal("Failed to initialize app:", err)
    }
    defer app.DB().Close()

    app.Use(middleware.ErrorHandler)
    app.Use(middleware.Logger)

    routes.RegisterUserRoutes(app)

    app.Listen(":8080")
}
```

### Route-Specific Middleware
Apply middleware to specific routes:

```go
func RegisterUserRoutes(app *framework.App) {
    router := app.Route("/users")
    router.Use(middleware.AuthMiddleware())
    router.
        GET("/", controllers.GetAllUsers(app)).
        GET("/{id}", controllers.GetUserByID(app)).
        POST("/", controllers.CreateUser(app))
}
```

### Custom Validation
Add custom validation rules to models:

```go
func (u *User) Validate() error {
    v := validator.New()
    v.RegisterValidation("custom_name", func(fl validator.FieldLevel) bool {
        return len(fl.Field().String()) > 3
    })
    return v.Struct(u)
}
```

Use it in the schema:

```go
Name string `json:"name" validate:"required,custom_name"`
```

## Database Configuration

The framework supports MySQL, PostgreSQL, SQLite, and MongoDB. Configure the database via `database.Config`.

### MySQL
Install MySQL and create a database:

```sql
CREATE DATABASE user_api;
```

Configure:

```go
dbConfig := database.Config{
    Type:     "mysql",
    Host:     "localhost",
    Port:     "3306",
    User:     "root",
    Password: "yourpassword",
    DBName:   "user_api",
}
```

### PostgreSQL
Install PostgreSQL and create a database:

```sql
CREATE DATABASE user_api;
```

Configure:

```go
dbConfig := database.Config{
    Type:     "postgres",
    Host:     "localhost",
    Port:     "5432",
    User:     "postgres",
    Password: "yourpassword",
    DBName:   "user_api",
}
```

### SQLite
No server setup required. Specify a file path:

```go
dbConfig := database.Config{
    Type:     "sqlite",
    FilePath: "user_api.db",
}
```

### MongoDB
Install MongoDB and ensure it’s running. Configure:

```go
dbConfig := database.Config{
    Type:     "mongodb",
    Host:     "localhost",
    Port:     "27017",
    User:     "",
    Password: "",
    DBName:   "user_api",
}
```

## Testing

Write unit tests for controllers and routes.

### Controller Test
Test the `CreateUser` controller:

```go
package controllers

import (
    "bytes"
    "encoding/json"
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/models"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestCreateUser(t *testing.T) {
    app, _ := framework.NewApp(database.Config{Type: "sqlite", FilePath: ":memory:"})
    user := models.User{UserSchema: models.UserSchema{Name: "Jane", Email: "jane@example.com"}}
    body, _ := json.Marshal(user)
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    CreateUser(app)(req, &framework.Response{w: w})

    if w.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", w.Code)
    }
}
```

### Route Test
Test the user routes:

```go
package routes

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestUserRoutes(t *testing.T) {
    app, _ := framework.NewApp(database.Config{Type: "sqlite", FilePath: ":memory:"})
    RegisterUserRoutes(app)

    req := httptest.NewRequest(http.MethodGet, "/users", nil)
    w := httptest.NewRecorder()
    app.router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", w.Code)
    }
}
```

Run tests:

```bash
go test ./...
```

## Examples

### Simple Hello World API
A minimal API:

```go
package main

import (
    "github.com/Mohammad007/GoExpressRestAPI/internal/database"
    "github.com/Mohammad007/GoExpressRestAPI/internal/framework"
    "github.com/Mohammad007/GoExpressRestAPI/internal/middleware"
)

func main() {
    dbConfig := database.Config{Type: "sqlite", FilePath: "test.db"}
    app, err := framework.NewApp(dbConfig)
    if err != nil {
        panic(err)
    }
    defer app.DB().Close()

    app.Use(middleware.Logger)

    app.Route("/hello").GET(func(r *http.Request, res *framework.Response) {
        res.Success("Hello, World!", nil)
    })

    app.Listen(":8080")
}
```

Test:

```bash
curl http://localhost:8080/hello
# Output: {"message":"Hello, World!"}
```

### User Management CRUD API
The full CRUD API example is covered in the [Building a CRUD API](#building-a-crud-api) section.

## Contributing

Contributions are welcome! To contribute:

1. Fork the repository: `https://github.com/Mohammad007/GoExpressRestAPI`
2. Create a feature branch: `git checkout -b feature-name`
3. Commit your changes: `git commit -m "Add feature-name"`
4. Push to the branch: `git push origin feature-name`
5. Open a pull request.

Please include tests and update the documentation for new features.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.