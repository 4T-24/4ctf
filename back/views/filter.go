package views

import (
	"4ctf/models"
	"reflect"
	"strings"
)

// Return filters the fields of a struct based on the `visible` tag
// It sets the fields to nil if the viewer is not allowed to see them
// The viewer is allowed to see the field if the `visible` tag contains:
// - `user` and the viewer is the owner of the data
// - `admin` and the viewer is an admin
// - `other` and the viewer is not the owner of the data
// - `nobody` and the viewer is not allowed to see the field
// The owner can be nil if the data is not owned by a user
func Return[K any](viewer *models.User, owner *models.User, data *K) *K {
	dataValue := reflect.ValueOf(data).Elem()

	// For each field in the struct
	for i := 0; i < dataValue.NumField(); i++ {
		field := dataValue.Field(i)
		fieldType := dataValue.Type().Field(i)

		// Get the `visible` tag
		visibleTag := fieldType.Tag.Get("visible")
		if visibleTag == "" {
			continue
		}

		// If the field is not visible, set it to nil
		if !isFieldVisible(visibleTag, viewer, owner) {
			field.Set(reflect.Zero(field.Type()))
		}
	}

	return data
}

func isFieldVisible(visibleTag string, viewer *models.User, owner *models.User) bool {
	roles := strings.Split(visibleTag, ",")
	for _, r := range roles {
		if owner != nil && viewer != nil && r == "user" {
			return owner.ID == viewer.ID
		}
		if viewer != nil && r == "admin" && viewer.IsAdmin {
			return true
		}
		if r == "other" {
			return true
		}
		if r == "nobody" {
			return false
		}
	}
	return false
}
