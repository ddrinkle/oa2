package oa2

type TokenResponse struct {
	AccessToken                 string `json:"access_token"`
	RefreshToken                string `json:"refresh_token"`
	UserInfo                    string `json:"user"`
	IdentityProviderAccessToken string `json:"ip_token"`
}
