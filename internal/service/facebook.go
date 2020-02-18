package service

import (
	"github.com/ddrinkle/oa2"
	"encoding/json"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type FacebookOA2Service struct {
	oa2.OA2Service
}
type facebookUser struct {
	Id        string `json:"id"`
	UserName  string `json:"name"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func NewFacebookService(env string) (oa2.OA2ServiceI, error) {

	return &FacebookOA2Service{
		oa2.OA2Service{
			UserInfoEndpoint: "https://graph.facebook.com/v2.5/me?fields=name,first_name,last_name,email,verified",
			Config: oauth2.Config{
				Endpoint: facebook.Endpoint,
				Scopes:   []string{"email", "public_profile"},
			},
		},
	}, nil
}

func (srv *FacebookOA2Service) ParseUserInfoResponse(response []byte) (oa2.AuthProviderUser, error) {
	user := facebookUser{}

	err := json.Unmarshal(response, &user)
	if err != nil {
		return oa2.AuthProviderUser{}, err
	}

	return oa2.AuthProviderUser{
		Id:            user.Id,
		UserName:      user.UserName,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		VerifiedEmail: user.Email != "",
		Email:         user.Email,
	}, nil

}
