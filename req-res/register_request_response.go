package reqres

type RegisterReq struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}
