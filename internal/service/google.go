package service

import (
	"github.com/ddrinkle/oa2"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOA2Service struct {
	oa2.OA2Service
}
type GoogleUser struct {
	Id            string `json:"id"`
	UserName      string `json:"name"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
}

func NewGoogleService(env string) (oa2.OA2ServiceI, error) {

	return &GoogleOA2Service{
		oa2.OA2Service{
			UserInfoEndpoint: "https://www.googleapis.com/oauth2/v1/userinfo",
			Config: oauth2.Config{
				Endpoint: google.Endpoint,
				Scopes:   []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
			},
		},
	}, nil
}

func (srv *GoogleOA2Service) ParseUserInfoResponse(response []byte) (oa2.AuthProviderUser, error) {

	user := GoogleUser{}

	err := json.Unmarshal(response, &user)
	if err != nil {
		return oa2.AuthProviderUser{}, err
	}

	return oa2.AuthProviderUser{
		Id:            user.Id,
		UserName:      user.UserName,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		VerifiedEmail: user.VerifiedEmail,
		Email:         user.Email,
	}, nil

}
