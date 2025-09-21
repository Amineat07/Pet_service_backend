package requestresponse

import "time"

type UpdateUserReq struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Password  *string `json:"password"`
}
type UsersResponse struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	IsCustomer bool      `json:"is_Customer"`
	IsProvider bool      `json:"is_Provider"`
	Created_At time.Time `json:"created_at"`
}

type ProviderResponse struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Email      string    `json:"email"`
	IsProvider bool      `json:"is_Provider"`
	Created_At time.Time `json:"created_at"`
}
