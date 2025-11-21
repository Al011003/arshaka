package auth

type RegisterResponse struct {
	NRA  string `json:"nra"`
	Role string `json:"role"`
}