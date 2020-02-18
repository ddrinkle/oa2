package oa2

import (
	"github.com/ddrinkle/platform/crypt"

	uuid "github.com/satori/go.uuid"
)

//AuthProvider is a storage structure for a AuthProvider
type AuthProvider struct {
	UUID                uuid.UUID     `json:"-" jsonapi:"primary,auth_provider"`
	Name                string        `json:"name,omitempty" jsonapi:"attr,name,omitempty"`
	OA2ServiceName      string        `json:"oa2_service_name,omitempty" jsonapi:"attr,oa2_service_name,omitempty"`
	ClientKey           *crypt.String `json:"client_key,omitempty" jsonapi:"attr,client_key,omitempty"`
	ClientSecret        *crypt.String `json:"secret_key,omitempty" jsonapi:"attr,secret_key,omitempty"`
	AuthorizationMethod string        `json:"authorization_method,omitempty" jsonapi:"attr,authorization_method,omitempty"`
}

type AuthProviders []AuthProvider

func (a AuthProvider) GetName() string {
	return "auth_provider"
}

func (a AuthProvider) GetID() string {
	return a.UUID.String()
}

func (a *AuthProvider) SetID(input string) error {
	var err error

	if input != "" {
		a.UUID, err = uuid.FromString(input)
	}

	return err
}
