package views

import (
	"4ctf/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ViewTest struct {
	Field1 *string `json:"field1" visible:"user,admin"`
	Field2 *string `json:"field2" visible:"other"`
	Field3 *string `json:"field3" visible:"nobody"`
	Field4 *string `json:"field4" visible:"admin"`
	Field5 *string `json:"field5" visible:"user"`
}

var (
	field1 = "field1"
	field2 = "field2"
	field3 = "field3"
	field4 = "field4"
	field5 = "field5"
)

func filledViewTest() ViewTest {
	return ViewTest{
		Field1: &field1,
		Field2: &field2,
		Field3: &field3,
		Field4: &field4,
		Field5: &field5,
	}
}

func TestFilter_AsUser(t *testing.T) {
	viewer := &models.User{ID: 1}
	owner := &models.User{ID: 1}
	data := filledViewTest()

	filteredData := Return(viewer, owner, &data)

	assert.Equal(t, data.Field1, filteredData.Field1)
	assert.Equal(t, data.Field2, filteredData.Field2)
	assert.Nil(t, filteredData.Field3)
	assert.Nil(t, filteredData.Field4)
	assert.Equal(t, data.Field5, filteredData.Field5)
}

func TestFilter_AsAdmin(t *testing.T) {
	viewer := &models.User{ID: 1, IsAdmin: true}
	owner := &models.User{ID: 2}
	data := filledViewTest()

	filteredData := Return(viewer, owner, &data)

	assert.Equal(t, data.Field1, filteredData.Field1)
	assert.Equal(t, data.Field2, filteredData.Field2)
	assert.Nil(t, filteredData.Field3)
	assert.Equal(t, data.Field4, filteredData.Field4)
	assert.Nil(t, filteredData.Field5)
}

func TestFilter_AsOther(t *testing.T) {
	viewer := &models.User{ID: 2}
	owner := &models.User{ID: 1}
	data := filledViewTest()

	filteredData := Return(viewer, owner, &data)

	assert.Nil(t, filteredData.Field1)
	assert.Equal(t, data.Field2, filteredData.Field2)
	assert.Nil(t, filteredData.Field3)
	assert.Nil(t, filteredData.Field4)
	assert.Nil(t, filteredData.Field5)
}

func TestFilter_AsOtherNil(t *testing.T) {
	owner := &models.User{ID: 1}
	data := filledViewTest()

	filteredData := Return(nil, owner, &data)

	assert.Nil(t, filteredData.Field1)
	assert.Equal(t, data.Field2, filteredData.Field2)
	assert.Nil(t, filteredData.Field3)
	assert.Nil(t, filteredData.Field4)
	assert.Nil(t, filteredData.Field5)
}

func TestFilter_AsUserNoOwner(t *testing.T) {
	viewer := &models.User{ID: 1}
	data := filledViewTest()

	filteredData := Return(viewer, nil, &data)

	assert.Nil(t, filteredData.Field1)
	assert.Equal(t, data.Field2, filteredData.Field2)
	assert.Nil(t, filteredData.Field3)
	assert.Nil(t, filteredData.Field4)
	assert.Nil(t, filteredData.Field5)
}

func TestFilter_AsAdminNoOwner(t *testing.T) {
	viewer := &models.User{ID: 1, IsAdmin: true}
	data := filledViewTest()

	filteredData := Return(viewer, nil, &data)

	assert.Equal(t, data.Field1, filteredData.Field1)
	assert.Equal(t, data.Field2, filteredData.Field2)
	assert.Nil(t, filteredData.Field3)
	assert.Equal(t, data.Field4, filteredData.Field4)
	assert.Nil(t, filteredData.Field5)
}
