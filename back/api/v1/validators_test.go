package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator_String(t *testing.T) {
	err := validatePassword("aaaaaaaaaaaaA1$", "")
	assert.Nil(t, err)

	err = validatePassword(1234568791254, "")
	assert.NotNil(t, err)
}

func TestPasswordValidator_Length(t *testing.T) {
	err := validatePassword("aaaA1$", "")
	assert.NotNil(t, err)
}

func TestPasswordValidator_Lowercase(t *testing.T) {
	err := validatePassword("AAAAAAAAAAAAA1$", "")
	assert.NotNil(t, err)
}

func TestPasswordValidator_Uppercase(t *testing.T) {
	err := validatePassword("aaaaaaaaaaaaa1$", "")
	assert.NotNil(t, err)
}

func TestPasswordValidator_Digit(t *testing.T) {
	err := validatePassword("aaaaaaaaaaaaA$", "")
	assert.NotNil(t, err)
}

func TestPasswordValidator_Special(t *testing.T) {
	err := validatePassword("aaaaaaaaaaaaA1", "")
	assert.NotNil(t, err)
}

func TestEmailValidator_Correct(t *testing.T) {
	err := validateEmail("test@test.com", "")
	assert.Nil(t, err)
}

func TestEmailValidator_Incorrect(t *testing.T) {
	err := validateEmail("test", "")
	assert.NotNil(t, err)

	err = validateEmail("test@.com", "")
	assert.NotNil(t, err)
}
