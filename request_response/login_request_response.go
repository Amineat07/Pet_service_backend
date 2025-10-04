package requestresponse

type LoginReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Firstname  string `json:"first_name"`
	Lastname   string `json:"last_name"`
	Email      string `json:"email"`
	IsCustomer bool   `json:"is_customer"`
	IsProvider bool   `json:"is_provider"`
	IsAdmin    bool   `json:"is_admin"`
	Token      string `json:"token"`
}
