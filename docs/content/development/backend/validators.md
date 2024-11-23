# Validator Documentation

The `validators` package is used to ensure that incoming request data is valid according to predefined rules. The package integrates with the API to perform validation checks before the request is processed further.

## Overview

Each incoming request body is validated using a set of validation rules defined in the struct tags. The validation is performed automatically by the `WithBody` middleware before passing control to the corresponding handler.

### Example: Register Request

Here’s an example of a `RegisterRequest` struct with validation rules:

```go
type RegisterRequest struct {
    Email    string `json:"email" validate:"nonzero,email"`
    Username string `json:"username" validate:"min=2,max=40,regexp=^[a-zA-Z]*$"`
    Password string `json:"password" validate:"max=40,password"`
}
```

In this example:

- `Email` must be non-empty and a valid email format.
- `Username` must be between 2 and 40 characters long and only contain alphabetic characters.
- `Password` must not exceed 40 characters and must meet certain password criteria (such as length, complexity, etc.).

The validation rules are specified using tags in the struct definition, following the syntax supported by the validator library.

## How It Works

The validation process is integrated into the request handling using a middleware function. The middleware automatically checks the request body against the validation rules defined in the model struct.

### `WithBody` Middleware

The `WithBody` middleware is responsible for:

1. **Parsing the body**: It unmarshals the incoming JSON request body into the specified struct (`K`).
2. **Validating the body**: It validates the parsed struct using the `validator.Validate` function. If the body does not meet the validation rules, the middleware returns a `400 Bad Request` response with the corresponding error messages.
3. **Calling the handler**: If the body passes validation, the middleware calls the original handler with the validated data.

Here is the middleware function:

```go
func WithBody[K any](fn func(api *Api) func(ctx *atreugo.RequestCtx, body K) error) func(api *Api) func(ctx *atreugo.RequestCtx) error {
    return func(api *Api) func(ctx *atreugo.RequestCtx) error {
        return func(ctx *atreugo.RequestCtx) error {
            var body K
            if err := json.Unmarshal(ctx.PostBody(), &body); err != nil {
                logrus.
                    WithField("request_ip", ctx.RemoteIP().String()).
                    WithError(err).
                    Warn("bad request")
                return ctx.JSONResponse(NewErrorResponse(400, []ResponseError{{Message: "Bad request"}}))
            }

            if err := validator.Validate(body); err != nil {
                logrus.
                    WithField("request_ip", ctx.RemoteIP().String()).
                    WithError(err).
                    Warn("invalid request")
                return ctx.JSONResponse(NewErrorResponse(400, validatorErrorToResponseError(err)))
            }

            return fn(api)(ctx, body)
        }
    }
}
```

### Example of Usage in a Route

You can use the `WithBody` middleware to validate the request body in your routes. For example, in the `register` route:

```go
group.POST("/register", WithDefaults(api, WithBody(register)))

func register(api *Api) func(ctx *atreugo.RequestCtx, rr *RegisterRequest) error {
    return func(ctx *atreugo.RequestCtx, rr *RegisterRequest) error {
        // handler logic
    }
}
```

In this example, the `register` handler receives a `RegisterRequest` object that has already been validated.

## Validator Rules

For default validators, check out the [package's documentation](https://github.com/go-validator/validator)

- `email`: Ensures the field is a valid email address.
- `password`: Custom password validation rule that include checks for complexity and length.

## Error Handling

If validation fails, the middleware will log the error and return a `400 Bad Request` response with detailed error messages. The `validatorErrorToResponseError` function formats the validation errors into a response-friendly format.

Example error response:

```json
{
  "status": 400,
  "errors": [
    {
      "field": "Email",
      "message": "Email must be a valid email address"
    }
  ]
}
```

## Best Practices

- **Use specific validation rules**: Apply the correct validation rule to each field (e.g., `email`, `min`, `max`, `regexp`).
- **Limit validation complexity**: Keep the validation rules simple and focused to avoid unnecessary overhead in request processing.
- **Log validation errors**: Ensure that all validation errors are logged to help with debugging and monitoring.
