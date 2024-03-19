package models

type Response struct {
	Status int `json:"status"`
	Body   any `json:"body"`
}

type SigninRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type SignupRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthCheckResponse struct {
	Login string `json:"login"`
}
