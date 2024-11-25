package v1

import (
	"4ctf/models"
	"4ctf/utils"
	"4ctf/views"
	"context"
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
	Username string `json:"username" validate:"min=2,max=40"`
	Password string `json:"password" validate:"max=256,password"`
}

func login(api *Api) func(ctx *atreugo.RequestCtx, lr *LoginRequest) *Response[any] {
	return func(ctx *atreugo.RequestCtx, lr *LoginRequest) *Response[any] {
		user, err := models.Users(models.UserWhere.Username.EQ(lr.Username)).OneG(context.Background())
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return Error(http.StatusNotFound, api.Translate("InvalidCredentials"))
			}
			return Error(http.StatusInternalServerError, api.Translate("InternalServerError"))
		}

		if !utils.VerifyPassword(lr.Password, user.PasswordHash) {
			return Error(http.StatusUnauthorized, api.Translate("InvalidCredentials"))
		}

		userSession := models.UserSession{
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		}
		if err := userSession.InsertG(context.Background(), boil.Whitelist(models.UserSessionColumns.UserID, models.UserSessionColumns.ExpiresAt)); err != nil {
			return Error(http.StatusInternalServerError, api.Translate("InternalServerError"))
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
		// Check if the username already exists
		exists, err := models.Users(models.UserWhere.Username.EQ(rr.Username)).ExistsG(context.Background())
		if err != nil {
			api.
				WithField("username", rr.Username).
				WithError(err).
				Error("cannot check user existence")
			return Error(http.StatusInternalServerError, api.Translate("InternalServerError"))
		}
		if exists {
			return Error(http.StatusBadRequest, api.Translate("UsernameAlreadyTaken"))
		}

		// Check if the email already exists
		exists, err = models.Users(models.UserWhere.Email.EQ(rr.Email)).ExistsG(context.Background())
		if err != nil {
			api.
				WithField("email", rr.Email).
				WithError(err).
				Error("cannot check email existence")
			return Error(http.StatusInternalServerError, api.Translate("InternalServerError"))
		}
		if exists {
			return Error(http.StatusBadRequest, api.Translate("EmailAlreadyInUse"))
		}

		// Hash the password
		hashedPassword := utils.HashPassword(rr.Password)

		// Create a new user in the database
		newUser := models.User{
			Username:     rr.Username,
			Email:        rr.Email,
			PasswordHash: hashedPassword,
		}
		if err := newUser.InsertG(context.Background(), boil.Whitelist(models.UserColumns.Username, models.UserColumns.Email, models.UserColumns.PasswordHash)); err != nil {
			api.
				WithField("username", rr.Username).
				WithError(err).
				Error("cannot create user")
			return Error(http.StatusInternalServerError, api.Translate("InternalServerError"))
		}

		api.
			WithField("username", rr.Username).
			Info("user registered successfully")

		return Success(http.StatusCreated, api.Translate("UserRegisteredSuccessfully"))
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

		userSession, err := models.UserSessions(models.UserSessionWhere.ID.EQ(session.UserSessionID)).OneG(context.Background())
		if err != nil {
			api.
				WithField("userSessionID", session.UserSessionID).
				WithError(err).
				Error("cannot find user session")
			return Error(http.StatusNotFound, api.Translate("UserNotFound"))
		}

		user, err := userSession.User().OneG(context.Background())
		if err != nil {
			api.
				WithField("userSessionID", session.UserSessionID).
				WithError(err).
				Error("cannot find user")
			return Error(http.StatusNotFound, api.Translate("UserNotFound"))
		}

		return Success(http.StatusOK, views.Return(user, user, user.View()))
	}
}
