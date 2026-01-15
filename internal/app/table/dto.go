package table

import (
	"github.com/casari-eat-n-go/backend/internal/pkg/ceng_pubsub"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

type ListTablesInputDto struct {
	IncludeClosed *bool `form:"includeClosed"`
	UserId        *uuid.UUID
}

func (r ListTablesInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.IncludeClosed, validation.In(true, false)),
	)
}

type getTableInputDto struct {
	ID     string `uri:"tableId"`
	UserId *uuid.UUID
}

func (r getTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}

type createTableInputDto struct {
	Name   string `json:"name"`
	UserId *uuid.UUID
}

func (r createTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
	)
}

type updateTableInputDto struct {
	ID            string                     `uri:"tableId"`
	Name          *string                    `json:"name"`
	IsOpen        *bool                      `json:"isOpen"`
	PaymentMethod *ceng_pubsub.PaymentMethod `json:"paymentMethod"`
	UserId        *uuid.UUID
}

func (r updateTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
		validation.Field(&r.Name, validation.NilOrNotEmpty, validation.Length(1, 255)),
		validation.Field(&r.IsOpen, validation.In(true, false)),
		validation.Field(&r.PaymentMethod, validation.NilOrNotEmpty, validation.In("card", "cash")),
	)
}

type deleteTableInputDto struct {
	ID     string `uri:"tableId"`
	UserId *uuid.UUID
}

func (r deleteTableInputDto) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUID),
	)
}
