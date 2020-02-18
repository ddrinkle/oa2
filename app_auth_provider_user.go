package oa2

import uuid "github.com/satori/go.uuid"

// AppAuthProviderUser is a Model for a AppAuthProviderUser
type AppAuthProviderUser struct {
	UUID         uuid.UUID     `json:"id,omitempty" jsonapi:"primary,app_auth_provider_user"`
	UserID       string        `json:"user_id,omitempty" jsonapi:"attr,user_id,omitempty"`
	UserName     string        `json:"user_name,omitempty" jsonapi:"attr,user_name,omitempty"`
	AccessToken  string        `json:"access_token,omitempty" jsonapi:"attr,access_token,omitempty"`
	RefreshToken *string       `json:"refresh_token,omitempty" jsonapi:"attr,refresh_token,omitempty"`
	AuthProvider *AuthProvider `json:"auth_provider,omitempty" jsonapi:"relation,auth_provider,omitempty"`
	App          *App          `json:"app" jsonapi:"relation,app,omitempty"`
}

//AppAuthProviderUsers is a convience type for an slice of AuthProviders
type AppAuthProviderUsers []AppAuthProviderUser
