package v1

import (
	"4ctf/models"
	"4ctf/utils"
	"4ctf/views"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/savsgio/atreugo/v11"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func setupAuthRoutes(api *Api, router *atreugo.Router) {
	group := router.NewGroupPath("/auth")

	group.POST("/login", Default(api, Body(login)))
	group.POST("/register", Default(api, Body(register)))
	group.GET("/logout", Default(api, logout))
	group.GET("/me", Default(api, RequireValidSession(profile)))
}

type LoginRequest struct {
	Username string `json:"username" validate:"min=2,max=40,regexp=^[a-zA-Z]*$"`
	Password string `json:"password" validate:"max=256,password"`
}

func login(api *Api) func(ctx *atreugo.RequestCtx, lr *LoginRequest) *Response[any] {
	return func(ctx *atreugo.RequestCtx, lr *LoginRequest) *Response[any] {
		// Check if the user exists in the database
		user, err := models.Users(models.UserWhere.Username.EQ(lr.Username)).OneG(ctx)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				api.
					WithField("username", lr.Username).
					WithError(err).
					Warn("cannot find user because it does not exist")
				return Error(http.StatusNotFound, "Invalid credentials")
			} else {
				api.
					WithField("username", lr.Username).
					WithError(err).
					Error("cannot find user")
				return Error(http.StatusInternalServerError, "Internal server error")
			}
		}

		// Check if the password is correct
		if !utils.VerifyPassword(lr.Password, user.PasswordHash) {
			api.
				WithField("username", lr.Username).
				Warn("invalid password")
			return Error(http.StatusUnauthorized, "Invalid credentials")
		}

		// Create a new userSession, this is used to link the cookie to the session, and not use the DB provider
		// The DB provider could break auto-generated code
		userSession := models.UserSession{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}

		if err := userSession.InsertG(ctx, boil.Whitelist(models.UserSessionColumns.UserID, models.UserSessionColumns.ExpiresAt)); err != nil {
			api.
				WithField("username", lr.Username).
				WithError(err).
				Error("cannot create user session")
			return Error(http.StatusInternalServerError, "Internal server error")
		}

		session := ctx.UserValue("session").(*Session)
		session.Valid = true
		session.UserSessionID = userSession.ID

		return Success(http.StatusOK, true)
	}
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"nonzero,email"`
	Username string `json:"username" validate:"min=2,max=40,regexp=^[a-zA-Z]*$"`
	Password string `json:"password" validate:"max=256,password"`
}

func register(api *Api) func(ctx *atreugo.RequestCtx, rr *RegisterRequest) *Response[any] {
	return func(ctx *atreugo.RequestCtx, rr *RegisterRequest) *Response[any] {
		// Check if the username or email already exists
		exists, err := models.Users(models.UserWhere.Username.EQ(rr.Username)).ExistsG(ctx)
		if err != nil {
			api.
				WithField("username", rr.Username).
				WithError(err).
				Error("cannot check user existence")
			return Error(http.StatusInternalServerError, "Internal server error")
		}
		if exists {
			return Error(400, "Username already taken")
		}

		exists, err = models.Users(models.UserWhere.Email.EQ(rr.Email)).ExistsG(ctx)
		if err != nil {
			api.
				WithField("email", rr.Email).
				WithError(err).
				Error("cannot check email existence")
			return Error(http.StatusInternalServerError, "Internal server error")
		}
		if exists {
			return Error(http.StatusBadRequest, "Email already in use")
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
			return Error(http.StatusInternalServerError, "Internal server error")
		}

		api.
			WithField("username", rr.Username).
			Info("user registered successfully")

		return Success(http.StatusCreated, true)
	}
}

func logout(api *Api) func(ctx *atreugo.RequestCtx) *Response[any] {
	return func(ctx *atreugo.RequestCtx) *Response[any] {
		api.session.DeleteSession(ctx.RequestCtx)

		return Success(http.StatusOK, true)
	}
}

func profile(api *Api) func(ctx *atreugo.RequestCtx) *Response[any] {
	return func(ctx *atreugo.RequestCtx) *Response[any] {
		session := ctx.UserValue("session").(*Session)

		userSession, err := models.UserSessions(models.UserSessionWhere.ID.EQ(session.UserSessionID)).OneG(ctx)
		if err != nil {
			api.
				WithField("userSessionID", session.UserSessionID).
				WithError(err).
				Error("cannot find user session")
			return Error(http.StatusNotFound, "User not found")
		}

		user, err := userSession.User().OneG(ctx)
		if err != nil {
			api.
				WithField("userSessionID", session.UserSessionID).
				WithError(err).
				Error("cannot find user")
			return Error(http.StatusNotFound, "User not found")
		}

		return Success(http.StatusOK, views.Return(user, user, user.View()))
	}
}
