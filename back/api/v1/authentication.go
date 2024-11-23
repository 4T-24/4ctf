package v1

import (
	"4ctf/models"
	"4ctf/utils"
	"4ctf/views"
	"database/sql"
	"errors"
	"time"

	"github.com/savsgio/atreugo/v11"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func setupAuthRoutes(api *Api, router *atreugo.Router) {
	group := router.NewGroupPath("/auth")

	group.POST("/login", WithDefaults(api, WithBody(login)))
	group.POST("/register", WithDefaults(api, WithBody(register)))
	group.GET("/logout", WithDefaults(api, logout))
	group.GET("/me", WithDefaults(api, WithValidSession(profile)))
}

type LoginRequest struct {
	Username string `json:"username" validate:"min=2,max=40,regexp=^[a-zA-Z]*$"`
	Password string `json:"password" validate:"max=40,password"`
}

func login(api *Api) func(ctx *atreugo.RequestCtx, lr *LoginRequest) error {
	return func(ctx *atreugo.RequestCtx, lr *LoginRequest) error {
		// Check if the user exists in the database
		user, err := models.Users(models.UserWhere.Username.EQ(lr.Username)).OneG(ctx)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				api.
					WithField("username", lr.Username).
					WithError(err).
					Error("cannot find user")
			} else {
				api.
					WithField("username", lr.Username).
					WithError(err).
					Warn("cannot find user")
			}
			return ctx.JSONResponse(NewErrorResponse(404, []ResponseError{{Message: "User not found"}}))
		}

		// Check if the password is correct
		if !utils.VerifyPassword(lr.Password, user.PasswordHash) {
			api.
				WithField("username", lr.Username).
				Warn("invalid password")
			return ctx.JSONResponse(NewErrorResponse(404, []ResponseError{{Message: "User not found"}}))
		}

		// Create a new userSession
		userSession := models.UserSession{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}

		if err := userSession.InsertG(ctx, boil.Whitelist(models.UserSessionColumns.UserID, models.UserSessionColumns.ExpiresAt)); err != nil {
			api.
				WithField("username", lr.Username).
				WithError(err).
				Error("cannot create user session")
			return ctx.JSONResponse(NewErrorResponse(500, []ResponseError{{Message: "Internal server error"}}))
		}

		session := ctx.UserValue("session").(*Session)
		session.Valid = true
		session.UserSessionID = userSession.ID

		return ctx.JSONResponse(NewResponse(200, true))
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"nonzero,email"`
	Username string `json:"username" validate:"min=2,max=40,regexp=^[a-zA-Z]*$"`
	Password string `json:"password" validate:"max=40,password"`
}

func register(api *Api) func(ctx *atreugo.RequestCtx, rr *RegisterRequest) error {
	return func(ctx *atreugo.RequestCtx, rr *RegisterRequest) error {
		// Check if the username or email already exists
		exists, err := models.Users(models.UserWhere.Username.EQ(rr.Username)).ExistsG(ctx)
		if err != nil {
			api.
				WithField("username", rr.Username).
				WithError(err).
				Error("cannot check user existence")
			return ctx.JSONResponse(NewErrorResponse(500, []ResponseError{{Message: "Internal server error"}}))
		}
		if exists {
			return ctx.JSONResponse(NewErrorResponse(400, []ResponseError{{Message: "Username already taken"}}))
		}

		exists, err = models.Users(models.UserWhere.Email.EQ(rr.Email)).ExistsG(ctx)
		if err != nil {
			api.
				WithField("email", rr.Email).
				WithError(err).
				Error("cannot check email existence")
			return ctx.JSONResponse(NewErrorResponse(500, []ResponseError{{Message: "Internal server error"}}))
		}
		if exists {
			return ctx.JSONResponse(NewErrorResponse(400, []ResponseError{{Message: "Email already in use"}}))
		}

		// Hash the password
		hashedPassword := utils.HashPassword(rr.Password)

		// Create a new user in the database
		newUser := models.User{
			Username:     rr.Username,
			Email:        rr.Email,
			PasswordHash: hashedPassword,
		}
		if err := newUser.InsertG(ctx, boil.Whitelist(models.UserColumns.Username, models.UserColumns.Email, models.UserColumns.PasswordHash)); err != nil {
			api.
				WithField("username", rr.Username).
				WithError(err).
				Error("cannot create user")
			return ctx.JSONResponse(NewErrorResponse(500, []ResponseError{{Message: "Internal server error"}}))
		}

		api.
			WithField("username", rr.Username).
			Info("user registered successfully")

		return ctx.JSONResponse(NewResponse(201, true), 201)
	}
}

func logout(api *Api) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {
		api.session.DeleteSession(ctx.RequestCtx)

		return ctx.JSONResponse(NewResponse(200, true))
	}
}

func profile(api *Api) func(ctx *atreugo.RequestCtx) error {
	return func(ctx *atreugo.RequestCtx) error {
		session := ctx.UserValue("session").(*Session)

		userSession, err := models.UserSessions(models.UserSessionWhere.ID.EQ(session.UserSessionID)).OneG(ctx)
		if err != nil {
			api.
				WithField("userSessionID", session.UserSessionID).
				WithError(err).
				Error("cannot find user session")
			return ctx.JSONResponse(NewErrorResponse(404, []ResponseError{{Message: "User not found"}}))
		}

		user, err := userSession.User().OneG(ctx)
		if err != nil {
			api.
				WithField("userSessionID", session.UserSessionID).
				WithError(err).
				Error("cannot find user")
			return ctx.JSONResponse(NewErrorResponse(404, []ResponseError{{Message: "User not found"}}))
		}

		return ctx.JSONResponse(NewResponse(200, views.Return(user, user, views.UserView(user))))
	}
}
