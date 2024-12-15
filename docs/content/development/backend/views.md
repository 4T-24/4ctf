# Model views

When generated, models will produce a struct that contains all the fields of the table, and will also produce the same struct with all optional fields.

## Overview

A **view** is a struct that defines which fields should be serialized and sent to the user, based on their visibility level. Each field in a view has a `visible` tag that specifies the roles or contexts in which the field should be included.

Those tags are controlled by the comment in the SQL schema, for example:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY COMMENT 'visible:"admin,user,other"',
    ...
);
```

This will generate the following view :

```go
type userView struct {
    ID                     *uint64      `json:"id,omitempty" visible:"admin,user,other"`
    ...
}
```

To convert a model to a view, you can use the `model.View` functionn :

```go
view := user.View()
```

### Visibility Levels

- `admin`: Visible to administrators.
- `user`: Visible to the user themselves.
- `other`: Visible to other roles (e.g., public viewers).
- `nobody`: Not visible to anyone.
- Fields can also use `json:"-"` to ensure they are not serialized at all.

## Creating a View

To create a new view, define a struct with the desired fields and `visible` tags. Implement a function to map the model data to the view struct. For example:

```go
func UserView(user *models.User) *userView {
    return &userView{
        ID:                     &user.ID,
        Username:               &user.Username,
        PasswordHash:           &user.PasswordHash,
        Email:                  &user.Email,
        EmailVerified:          &user.EmailVerified,
        EmailVerificationToken: &user.EmailVerificationToken,
        IsAdmin:                &user.IsAdmin,
        IsHidden:               &user.IsHidden,
        CreatedAt:              &user.CreatedAt,
        UpdatedAt:              &user.UpdatedAt,
        DeletedAt:              &user.DeletedAt,
    }
}
```

## Using Views in Handlers

To use a view in an API handler, call the appropriate view function and return it in the response. For example:

```go
func profile(api *Api) func(ctx *atreugo.RequestCtx) error {
    return func(ctx *atreugo.RequestCtx) error {
        session := ctx.UserValue("session").(*Session)

        userSession, err := models.UserSessions(models.UserSessionWhere.ID.EQ(session.UserSessionID)).OneG(ctx)
        if err != nil {
            ...
        }

        user, err := userSession.User().OneG(ctx)
        if err != nil {
            ...
        }

        return ctx.JSONResponse(NewResponse(200, views.Return(user, user, views.UserView(user))))
    }
}
```

The params are explained by your IDE, the first one is the viewer, the second one the owner of the resource, and the third one is the model you want to send.

### Explanation of `views.Return`

The `views.Return` function is responsible for applying the visibility rules to the view. It takes the current user (to determine their role) and the data to be serialized.

## Updating Views

To update an existing view:

1. Modify the struct to add, remove, or update fields.
2. Update the `visible` tags to reflect the new visibility rules.
3. Update the corresponding view function to map any new or removed fields.

## Adding a New View

To add a new view:

1. Define a new struct in the `views` package.
2. Add appropriate `json` and `visible` tags to each field.
3. Implement a function to map the model data to the view struct.
4. Use the new view in your handlers as needed.

### Example: Creating a `productView`

```go
type productView struct {
    ID          *uint64 `json:"id,omitempty" visible:"admin,user,other"`
    Name        *string `json:"name,omitempty" visible:"admin,user,other"`
    Description *string `json:"description,omitempty" visible:"admin,user"`
    Price       *float64 `json:"price,omitempty" visible:"admin,user"`
}

func ProductView(product *models.Product) *productView {
    return &productView{
        ID:          &product.ID,
        Name:        &product.Name,
        Description: &product.Description,
        Price:       &product.Price,
    }
}
```

## Best Practices

- **Use `visible` tags carefully**: Avoid exposing sensitive information (e.g., passwords, tokens) to unintended roles.
- **Keep view structs minimal**: Only include fields that are necessary for the given context.
- **Test visibility rules**: Ensure fields are visible only to the intended roles.

## Future Improvements

- Implement a utility function to automatically filter fields based on the `visible` tag, reducing boilerplate in handlers.
  - This one might be tricky to implement and the benefits might not be worth the complexity.
