package menuItem

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ListMenuItemsInputDto struct {
	MenuCategoryId string `uri:"menuCategoryId"`
}

func (r ListMenuItemsInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuCategoryId, validation.Required, is.UUID),
	)
}

type createMenuItemInputDto struct {
	MenuCategoryId string `uri:"menuCategoryId"`
	Title          string `json:"title"`
	Price          int64  `json:"price"`
}

func (r createMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MenuCategoryId, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Price, validation.Required, validation.Min(1), validation.Max(10000)),
	)
}

type getMenuItemInputDto struct {
	ID string `uri:"menuItemId"`
}

func (r getMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type updateMenuItemInputDto struct {
	ID       string  `uri:"menuItemId"`
	Title    *string `json:"title"`
	Position *int64  `json:"position"`
	Active   *bool   `json:"active"`
	Price    *int64  `json:"price"`
}

func (r updateMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Title, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Position, validation.NilOrNotEmpty, validation.Min(1)),
		validation.Field(&r.Active, validation.In(true, false)),
		validation.Field(&r.Price, validation.Min(1), validation.Max(10000)),
	)
}

type deleteMenuItemInputDto struct {
	ID string `uri:"menuItemId"`
}

func (r deleteMenuItemInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
