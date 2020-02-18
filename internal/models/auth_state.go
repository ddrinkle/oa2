package models

import (
	"encoding/base64"
	"encoding/json"
)

type AuthState struct {
	AuthProviderID string `json:"auth_provider_id"`
	AppID          string `json:"app_id"`
	ReturnURL      string `json:"return"`
	Random         int    `json:"random"`
}

func (a AuthState) Encode() (string, error) {
	state, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.WithPadding(base64.StdPadding).EncodeToString(state), nil
}

func (a *AuthState) Decode(encState string) error {

	decodedState, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(encState)
	if err != nil {
		return err
	}

	err = json.Unmarshal(decodedState, a)
	return err
}
