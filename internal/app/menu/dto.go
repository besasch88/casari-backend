package menu

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type getMenuInputDto struct {
	Target string `form:"target"`
}

func (r getMenuInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Target, validation.Required, validation.In("inside", "outside")),
	)
}
