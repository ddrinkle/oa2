package service

import (
	"github.com/ddrinkle/oa2"
	"encoding/json"
	"strconv"

	"golang.org/x/oauth2"
)

type EveOA2Service struct {
	oa2.OA2Service
}

type EveUser struct {
	Id                 int64  `json:"CharacterID"`
	CharacterName      string `json:"CharacterName"`
	CharacterOwnerHash string `json:"CharacterOwnerHash"`
}

func NewEveService(env string) (oa2.OA2ServiceI, error) {

	return &EveOA2Service{
		oa2.OA2Service{
			UserInfoEndpoint: "https://login.eveonline.com/oauth/verify",
			Config: oauth2.Config{
				Endpoint: oauth2.Endpoint{
					AuthURL:  "https://login.eveonline.com/v2/oauth/authorize",
					TokenURL: "https://login.eveonline.com/v2/oauth/token",
				},
				Scopes: []string{"esi-location.read_location.v1", "esi-planets.manage_planets.v1", "esi-assets.read_assets.v1", "esi-assets.read_corporation_assets.v1", "esi-universe.read_structures.v1", "esi-wallet.read_character_wallet.v1", "esi-markets.read_character_orders.v1", "esi-assets.read_assets.v1"},
				//Scopes: []string{"esi-planets.manage_planets.v1"},
			},
		},
	}, nil
}

func (srv *EveOA2Service) ParseUserInfoResponse(response []byte) (oa2.AuthProviderUser, error) {

	user := EveUser{}

	err := json.Unmarshal(response, &user)
	if err != nil {
		return oa2.AuthProviderUser{}, err
	}

	return oa2.AuthProviderUser{
		Id:       strconv.FormatInt(user.Id, 10),
		UserName: user.CharacterName,
	}, nil

}
func (srv *EveOA2Service) GetAuthorizationMethod() oa2.AUTHORIZATION_METHOD {
	return oa2.AUTHORIZATION_METHOD_HEADER_BEARER
}
