package views

import (
	"4ctf/models"
	"time"

	"github.com/volatiletech/null/v8"
)

type userView struct {
	ID                     *uint64      `json:"id,omitempty" visible:"admin,user,other"`
	Username               *string      `json:"username,omitempty" visible:"admin,user,other"`
	PasswordHash           *string      `json:"password_hash,omitempty" visible:"nobody"`
	Email                  *string      `json:"email,omitempty" visible:"admin,user"`
	EmailVerified          *bool        `json:"email_verified,omitempty" visible:"admin,user"`
	EmailVerificationToken *null.String `json:"email_verification_token,omitempty" visible:"admin"`
	IsAdmin                *bool        `json:"is_admin,omitempty" visible:"admin"`
	IsHidden               *bool        `json:"is_hidden,omitempty" visible:"admin"`
	CreatedAt              *time.Time   `json:"-" visible:"admin"`
	UpdatedAt              *time.Time   `json:"-" visible:"admin"`
	DeletedAt              *null.Time   `json:"-" visible:"admin"`
}

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
