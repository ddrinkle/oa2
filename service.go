package oa2

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

type AUTHORIZATION_METHOD string

const (
	AUTHORIZATION_METHOD_HEADER_OAUTH    AUTHORIZATION_METHOD = "header_oauth"
	AUTHORIZATION_METHOD_HEADER_BEARER                        = "header_bearer"
	AUTHORIZATION_METHOD_QUERY_STRING                         = "qs_access_token"
	AUTHORIZATION_METHOD_QUERY_STRING_V2                      = "qs_oauth2_access_token"
	AUTHORIZATION_METHOD_QUERY_STRING_V3                      = "qs_apikey"
	AUTHORIZATION_METHOD_QUERY_STRING_V4                      = "qs_auth"
)

type OA2ServiceI interface {
	GetOAuth2Config() oauth2.Config

	GetTokenCollection() *oauth2.Token
	GetClientKey() string
	GetClientSecret() string

	SetTokenCollection(collection *oauth2.Token)
	SetRedirectURL(redirectURL string)
	SetContext(ctx context.Context)
	SetClientKey(string)
	SetClientSecret(string)

	GetAuthorizationMethod() AUTHORIZATION_METHOD

	//Access Token Methods
	RequestAccessTokenCollection(code string) (*oauth2.Token, error)
	RefreshAccessTokenInCollection(token *oauth2.Token) (*oauth2.Token, error)

	//UserInfo Endpoint Methods
	GetUserInfoEndpoint() string
	ParseUserInfoResponse(response []byte) (AuthProviderUser, error)
	RequestUserInfo(authMethod AUTHORIZATION_METHOD) ([]byte, error)

	//TrustedLogin
	IsTrustedLoginService() bool
}

type OA2Service struct {
	UserInfoEndpoint string
	oauth2.Config
	Token *oauth2.Token
	Ctx   *context.Context
}

type AuthProviderUser struct {
	Id            string `json:"id,omitempty"`
	UserName      string `json:"user_name,omitempty"`
	FirstName     string `json:"first_name,omitempty"`
	LastName      string `json:"last_name,omitempty"`
	Email         string `json:"email,omitempty"`
	VerifiedEmail bool   `json:"verified_email,omitempty"`
}

func (s *OA2Service) GetTokenCollection() *oauth2.Token {
	return s.Token
}
func (s *OA2Service) SetTokenCollection(collection *oauth2.Token) {
	s.Token = collection
}
func (s *OA2Service) SetRedirectURL(redirectURL string) {
	s.RedirectURL = redirectURL
}
func (s *OA2Service) GetClientKey() string {
	return s.ClientID
}

func (s *OA2Service) GetClientSecret() string {
	return s.ClientSecret
}

func (s *OA2Service) GetOAuth2Config() oauth2.Config {
	return s.Config
}

func (s *OA2Service) GetUserInfoEndpoint() string {
	return s.UserInfoEndpoint
}

func (s *OA2Service) SetContext(ctx context.Context) {
	s.Ctx = &ctx
}
func (s *OA2Service) SetClientKey(key string) {
	s.ClientID = key
}
func (s *OA2Service) SetClientSecret(secret string) {
	s.ClientSecret = secret
}
func (s *OA2Service) RequestAccessTokenCollection(code string) (*oauth2.Token, error) {

	if s.Ctx == nil {
		return nil, errors.New("Context can't be nil")
	}
	s.Config.RedirectURL = "http://localhost:8001/oauth/token"

	tokenCollection, err := s.Config.Exchange(*s.Ctx, code)

	fmt.Println(*tokenCollection)

	s.Token = tokenCollection
	return tokenCollection, err
}
func (s *OA2Service) RefreshAccessTokenInCollection(tokenCollection *oauth2.Token) (*oauth2.Token, error) {
	if s.Ctx == nil {
		return nil, errors.New("Context can't be nil")
	}

	tokenSource := s.Config.TokenSource(*s.Ctx, tokenCollection)

	tokenCollection, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return tokenCollection, nil
}

func (s *OA2Service) ParseUserInfoResponse(response []byte) (AuthProviderUser, error) {
	return AuthProviderUser{}, errors.New("ParseUserInfoResponse should be overridden by each service")
}
func (s *OA2Service) RequestUserInfo(authMethod AUTHORIZATION_METHOD) ([]byte, error) {

	userInfoURL, err := url.Parse(s.GetUserInfoEndpoint())
	if err != nil {
		return nil, err
	}

	v, err := url.ParseQuery(userInfoURL.RawQuery)
	if err != nil {
		return nil, err
	}

	if s.Token.AccessToken == "" {
		return nil, errors.New("Missing Access Token")
	}

	var headers = map[string]string{}

	switch authMethod {
	case AUTHORIZATION_METHOD_HEADER_OAUTH:
		headers["Authorization"] = "OAuth " + s.Token.AccessToken
	case AUTHORIZATION_METHOD_HEADER_BEARER:
		headers["Authorization"] = "Bearer " + s.Token.AccessToken
	case AUTHORIZATION_METHOD_QUERY_STRING:
		v.Add("access_token", s.Token.AccessToken)
	case AUTHORIZATION_METHOD_QUERY_STRING_V2:
		v.Add("oauth2_access_token", s.Token.AccessToken)
	case AUTHORIZATION_METHOD_QUERY_STRING_V3:
		v.Add("apikey", s.Token.AccessToken)
	case AUTHORIZATION_METHOD_QUERY_STRING_V4:
		v.Add("auth", s.Token.AccessToken)
	}

	userInfoURL.RawQuery = v.Encode()

	client := &http.Client{}

	req, err := http.NewRequest("GET", userInfoURL.String(), nil)

	for head, value := range headers {
		req.Header.Set(head, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (s *OA2Service) GetAuthorizationMethod() AUTHORIZATION_METHOD {
	return AUTHORIZATION_METHOD_HEADER_OAUTH
}
func (s *OA2Service) IsTrustedLoginService() bool {
	return false
}
