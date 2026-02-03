package table

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type listTablesInputDto struct {
	IncludeClosed *bool  `form:"includeClosed"`
	Target        string `form:"target"`
}

func (r listTablesInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.IncludeClosed, validation.In(true, false)),
		validation.Field(&r.Target, validation.Required, validation.In("inside", "outside")),
	)
}

type getTableInputDto struct {
	ID string `uri:"tableId"`
}

func (r getTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type createTableInputDto struct {
	Name   string `json:"name"`
	Inside *bool  `json:"inside"`
}

func (r createTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Inside, validation.In(true, false)),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
	)
}

type updateTableInputDto struct {
	ID            string  `uri:"tableId"`
	Name          *string `json:"name"`
	Inside        *bool   `json:"inside"`
	Close         *bool   `json:"close"`
	PaymentMethod *string `json:"paymentMethod"`
}

func (r updateTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Name, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.Inside, validation.In(true, false)),
		validation.Field(&r.Close, validation.In(true, false)),
		validation.Field(&r.PaymentMethod, validation.NilOrNotEmpty, validation.In("card", "cash")),
	)
}

type deleteTableInputDto struct {
	ID string `uri:"tableId"`
}

func (r deleteTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
