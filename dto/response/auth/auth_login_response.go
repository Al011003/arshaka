package auth

type LoginResponse struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	MustChangePassword bool   `json:"must_change_password"`
	MustFillEmail      bool   `json:"must_fill_email"`
	Role               string `jason:"role"`
}