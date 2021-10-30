package models

// swagger:model
type LoginRequest struct {
	IDToken string `json:"id_token"`
}

// swagger:model
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// swagger:model
type LoginResponse struct {
	TokenPair TokenPair `json:"tokens_pair"`
	User      User      `json:"user"`
}
