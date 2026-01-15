package course

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ListCoursesInputDto struct {
	TableId string `uri:"tableId"`
}

func (r ListCoursesInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TableId, validation.Required, is.UUID),
	)
}

type createCourseInputDto struct {
	TableId string `uri:"tableId"`
}

func (r createCourseInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.TableId, validation.Required, is.UUID),
	)
}

type getCourseInputDto struct {
	ID string `uri:"courseId"`
}

func (r getCourseInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type updateCourseInputDto struct {
	ID    string `uri:"courseId"`
	Close *bool  `json:"close"`
}

func (r updateCourseInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Close, validation.In(true, false)),
	)
}

type deleteCourseInputDto struct {
	ID string `uri:"courseId"`
}

func (r deleteCourseInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
