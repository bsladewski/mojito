package delivery

// loginRequest is used to read a request to the login endpoint.
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse is used to format responses from the login endpoint.
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// refreshRequest is used to read a request to the refresh endpoint.
type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// refreshResponse is used to format responses from the refresh endpoint.
type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
