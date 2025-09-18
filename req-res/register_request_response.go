package reqres

type RegisterReq struct {
	Firstname         string `json:"first_name"`
	Lastname          string `json:"lastname"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	IsCustomer        bool   `json:"is_customer"`
	IsServiceProvider bool   `json:"is_service_provider"`
}
