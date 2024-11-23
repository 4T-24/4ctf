# Middleware Documentation

The following middlewares are used to handle common tasks such as session management, request validation, and ensuring the integrity of incoming data. These middlewares are used in the API routes to ensure that the incoming requests meet certain conditions before passing them to the corresponding handlers.

## Overview

- **WithDefaults**: Ensures that session information is loaded, and the logger is initialized for each request.
- **WithBody**: Validates the body of the request according to the defined validation rules.
- **WithValidSession**: Ensures that the session is valid before processing the request.

### Middleware Functions

Each middleware function takes an API instance (`*Api`) and returns a handler function that performs some processing before calling the original handler.

---

## `WithDefaults`

The `WithDefaults` middleware sets up the session and logger for each request. It ensures that the session information is properly loaded and updated during the request handling.

### Function Signature

```go
func WithDefaults(api *Api, fn func(api *Api) func(ctx *atreugo.RequestCtx) error) func(ctx *atreugo.RequestCtx) error
```

### How It Works

1. The middleware retrieves the session information using `api.session.GetSession`.
2. It sets the session into the request context so that it can be accessed by downstream handlers.
3. The `api` object is reinitialized with the session and logger for the specific request.
4. The original handler is executed, and after the handler runs, it checks if the session has been updated. If so, the session is saved back.

### Example Usage

```go
group.POST("/some-endpoint", WithDefaults(api, someHandler))
```

This will ensure that the session is available in the request context and that logging is correctly handled.

---

## `WithBody`

The `WithBody` middleware is used to parse and validate the request body. It unmarshals the incoming JSON data into a Go struct and validates it according to the rules defined in the struct tags. If validation fails, a `400 Bad Request` response is returned.

### Function Signature

```go
func WithBody[K any](fn func(api *Api) func(ctx *atreugo.RequestCtx, body K) error) func(api *Api) func(ctx *atreugo.RequestCtx) error
```

### How It Works

1. The middleware unmarshals the incoming JSON request body into a struct of type `K`.
2. It validates the struct using the `validator.Validate` function.
3. If the body is invalid, a `400 Bad Request` response is returned with details of the validation errors.
4. If the body is valid, the original handler is called with the validated data.

### Example Usage

```go
func register(api *Api) func(ctx *atreugo.RequestCtx, rr *RegisterRequest) error {
    return func(ctx *atreugo.RequestCtx, rr *RegisterRequest) error {
        // Handler logic
    }
}

group.POST("/register", WithDefaults(api, WithBody(register)))
```

In this example, the body of the request is validated using the rules defined in the `RegisterRequest` struct before the `register` handler is invoked.

---

## `WithValidSession`

The `WithValidSession` middleware ensures that the current user session is valid before proceeding with the request. If the session is not valid, it returns a `401 Unauthorized` response.

### Function Signature

```go
func WithValidSession(fn func(api *Api) func(ctx *atreugo.RequestCtx) error) func(api *Api) func(ctx *atreugo.RequestCtx) error
```

### How It Works

1. The middleware checks if the session associated with the request is valid.
2. If the session is not valid, it returns a `401 Unauthorized` response.
3. If the session is valid, the original handler is called.

### Example Usage

```go
group.GET("/protected-endpoint", WithDefaults(api, WithValidSession(protectedHandler)))
```

In this example, the `protectedHandler` will only be executed if the session is valid.

---

## Error Handling

For each middleware, if an error occurs (e.g., invalid request body, invalid session), an appropriate error response is returned:

- `400 Bad Request` for invalid request bodies or validation errors.
- `401 Unauthorized` for invalid or expired sessions.

Each middleware logs the error details, including the IP address of the requester, to assist in debugging and monitoring.

---

## Best Practices

- **Session Management**: Use `WithDefaults` and `WithValidSession` to ensure that session data is consistently available and valid.
- **Request Validation**: Use `WithBody` for automatic validation of incoming request bodies to prevent invalid data from reaching the handler.
- **Error Handling**: Ensure that each middleware returns appropriate error responses and logs relevant information for debugging.

---
