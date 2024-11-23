package v1

import (
	"errors"
	"net/mail"
	"reflect"

	"gopkg.in/validator.v2"
)

func validatePassword(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("password must be string")
	}

	// Password must be at least 12 characters long
	if len(st.String()) < 12 {
		return errors.New("password must be at least 12 characters long")
	}

	// Password must contain :
	// - at least one lowercase letter
	// - at least one uppercase letter
	// - at least one digit
	// - at least one special character
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSpecial := false
	for _, c := range st.String() {
		switch {
		case 'a' <= c && c <= 'z':
			hasLower = true
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case '0' <= c && c <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func validateEmail(v interface{}, param string) error {
	st := reflect.ValueOf(v)
	if st.Kind() != reflect.String {
		return errors.New("email must be string")
	}

	_, err := mail.ParseAddress(st.String())
	return err
}

func setupValidators() {
	validator.SetValidationFunc("password", validatePassword)
	validator.SetValidationFunc("email", validateEmail)
}
