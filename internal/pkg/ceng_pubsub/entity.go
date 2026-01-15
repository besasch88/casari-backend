package ceng_pubsub

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentCash PaymentMethod = "cash"
	PaymentCard PaymentMethod = "card"
)

type MenuCategoryEventEntity struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Position  int64     `json:"position"`
	Active    *bool     `json:"active"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type MenuItemEventEntity struct {
	ID             uuid.UUID `json:"id"`
	MenuCategoryID uuid.UUID `json:"menuCategoryId"`
	Title          string    `json:"title"`
	Position       int64     `json:"position"`
	Active         *bool     `json:"active"`
	Price          int64     `json:"price"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type MenuOptionEventEntity struct {
	ID         uuid.UUID `json:"id"`
	MenuItemID uuid.UUID `json:"menuItemId"`
	Title      string    `json:"title"`
	Position   int64     `json:"position"`
	Active     *bool     `json:"active"`
	Price      int64     `json:"price"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TableEventEntity struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.UUID      `json:"userId"`
	Name          string         `json:"name"`
	Close         *bool          `json:"close"`
	PaymentMethod *PaymentMethod `json:"paymentMethod"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
}
