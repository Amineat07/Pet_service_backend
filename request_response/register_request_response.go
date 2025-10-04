package requestresponse

import "time"

type RegisterReq struct {
	Firstname         string `json:"first_Name" validate:"required"`
	Lastname          string `json:"last_Name" validate:"required"`
	Email             string `json:"email" validate:"required"`
	Password          string `json:"password" validate:"required"`
	IsCustomer        *bool   `json:"is_Customer" validate:"required"`
	IsServiceProvider *bool   `json:"is_Service_Provider" validate:"required"`
	CreatedAt         time.Time
}
