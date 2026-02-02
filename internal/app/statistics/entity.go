package statistics

import "time"

type paymentTakingsEntity struct {
	PaymentType string `json:"paymentType"`
	Takings     int64  `json:"takings"`
}

type menuItemStatEntity struct {
	Title    string `json:"title"`
	Quantity int64  `json:"quantity"`
	Takings  int64  `json:"takings"`
}

type statisticsEntity struct {
	AvgTableDuration time.Duration          `json:"avgTableDuration"`
	PaymentsTakins   []paymentTakingsEntity `json:"paymentsTakings"`
	MenuItemStats    []menuItemStatEntity   `json:"menuItemStats"`
}
