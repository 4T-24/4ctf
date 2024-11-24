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
